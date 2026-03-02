#!/bin/bash
set -euo pipefail

mode=${1:-run}
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

MOUNT="$BASE_DIR/mnt"

CERTBOT_CONF="$MOUNT/certbot/conf"
CERTBOT_WWW="$MOUNT/certbot/www"
mkdir -p "$CERTBOT_CONF"
mkdir -p "$CERTBOT_WWW"

NGINX_SRC="$BASE_DIR/nginx"
NGINX_LOGS="$MOUNT/nginx/logs"
mkdir -p "$NGINX_LOGS"
if [[ ! -e "$NGINX_LOGS/access.log" ]]; then 
    touch  "$NGINX_LOGS/access.log"
    touch  "$NGINX_LOGS/error.log"
fi

OHSAL="$MOUNT/ohsal"
mkdir -p "$OHSAL/logs"
mkdir -p "$OHSAL/data"
if [[ ! -e "$OHSAL/logs/ohsal.log" ]]; then 
    touch  "$OHSAL/logs/ohsal.log"
fi

##################################################
#############   DOCKER    ########################
##################################################

NETWORK_NAME="web-servers"
create_network() {
    if ! docker network ls | grep -q "$NETWORK_NAME"; then   
        docker network create "$NETWORK_NAME"
    fi
}


##################################################
#############   PROXY    #########################
##################################################


RENEW_CRON_TAG="# ssl_renew"

start_proxy_prod() {


    ############## NGINX ################
    if ! docker ps | grep -q "nginx"; then
        docker run -itd \
            --name nginx \
            --network "$NETWORK_NAME" \
            -p 80:80 -p 443:443 \
            -v "$NGINX_SRC/prod:/etc/nginx/conf.d:ro" \
            -v "$CERTBOT_CONF:/etc/letsencrypt" \
            -v "$CERTBOT_WWW:/var/www/certbot" \
            -v "$NGINX_LOGS:/var/log/nginx" \
            --restart unless-stopped \
            nginx:latest
    fi


   
    ############## CERTBOT CRON JOB ################
   
    CRON_SCHEDULE="0 */12 * * *"
    CRON_CMD="$BASE_DIR/renew_ssl.sh >> $BASE_DIR/proxy/cron.log 2>&1"
    CRON_LINE="$CRON_SCHEDULE $CRON_CMD $RENEW_CRON_TAG"
    
    if ! crontab -l 2>/dev/null | grep -Fq "$RENEW_CRON_TAG"; then
        (crontab -l 2>/dev/null; echo "$CRON_LINE") | crontab -
    fi
}

start_proxy_first_time() {

    if [[ ! -d "$CERTBOT_CONF/accounts" ]]; then
        ############## NGINX ################
        docker run -itd --rm \
            --name nginx \
            --network "$NETWORK_NAME" \
            -p 80:80 \
            -v "$NGINX_SRC/temp:/etc/nginx/conf.d:ro" \
            -v "$CERTBOT_CONF:/etc/letsencrypt" \
            -v "$CERTBOT_WWW:/var/www/certbot" \
            -v "$NGINX_LOGS:/var/log/nginx" \
            nginx:latest

        ############## CERTBOT ################
        docker run --rm \
            -v "$CERTBOT_CONF:/etc/letsencrypt" \
            -v "$CERTBOT_WWW:/var/www/certbot" \
            certbot/certbot certonly --webroot -w /var/www/certbot \
            --non-interactive --agree-tos \
            --email joshua.olson1@yahoo.com \
            -d www.testohsal.com \
            -d testohsal.com

    fi
}

stop_proxy() {

    docker stop nginx
    docker rm nginx

    if crontab -l 2>/dev/null | grep -Fq "$RENEW_CRON_TAG"; then
        crontab -l 2>/dev/null | grep -Fv "$RENEW_CRON_TAG" | crontab -
    fi
}

##################################################
#############   OHSAL    #########################
##################################################
start_ohsal() {
    if ! docker images | grep -Fq ohsal.com; then
        docker build -t ohsal.com .
    fi
    if docker ps -a | grep -Fq ohsal; then
        docker start ohsal
    else
        docker run -d \
            --name ohsal --network "$NETWORK_NAME" \
            -v "$OHSAL:/var/ohsal" \
            ohsal.com:latest
    fi
}

stop_ohsal() {
    docker stop ohsal
    docker rm ohsal
}

##################################################
#############   FAIL2BAN    ######################
##################################################
set_fail_to_ban() {
    local action_file="/etc/fail2ban/action.d/iptables-docker.conf"
    cat "$NGINX_SRC/fail2ban/iptables-docker.conf" > "$action_file" 

    local filter_file="/etc/fail2ban/filter.d/nginx-scan.conf"
    cat "$NGINX_SRC/fail2ban/nginx-scan.conf" > "$filter_file" 

    local jail_file="/etc/fail2ban/jail.d/nginx-scan.local"
    sed "s|NGINX_LOGS|$NGINX_LOGS|g" /$NGINX_SRC/fail2ban/nginx-scan.local > "$jail_file"
    chmod 644 "$jail_file"

    systemctl reload fail2ban
}


##################################################
#############   MAIN    ##########################
##################################################
if [[ "$mode" == "run" ]]; then
    set_fail_to_ban
    create_network
    start_proxy_first_time
    start_proxy_prod
    start_ohsal

elif [[ "$mode" == "stop" ]]; then
    stop_ohsal
    stop_proxy
else
    echo "Unknown Mode"
fi
