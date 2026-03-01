#!/bin/bash
set -euo pipefail

mode=${1:-run}
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

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

start_proxy() {

    CERTBOT_CONF="$BASE_DIR/certbot/conf"
    CERTBOT_WWW="$BASE_DIR/certbot/www"


    ############## NGINX ################
    if ! docker ps -a | grep -q "nginx"; then
        docker run -d \
            --name nginx \
            --network "$NETWORK_NAME" \
            -p 80:80 -p 443:443 \
            -v "$BASE_DIR/nginx/conf.d:/etc/nginx/conf.d" \
            -v "$CERTBOT_CONF:/etc/letsencrypt" \
            -v "$CERTBOT_WWW:/var/www/certbot" \
            --restart unless-stopped \
            nginx:latest
    fi


    ############## CERTBOT ################
    domains=(testohsal.com)

    docker run --rm \
        -v "$CERTBOT_CONF:/etc/letsencrypt" \
        -v "$CERTBOT_WWW:/var/www/certbot" \
        certbot/certbot certonly --webroot -w /var/www/certbot \
        "${domains[@]/#/-d }"

    ############## CERTBOT CRON JOB ################

    CRON_SCHEDULE="0 */12 * * *"
    CRON_CMD="$BASE_DIR/proxy/renew_ssl.sh >> $BASE_DIR/proxy/cron.log 2>&1"
    CRON_LINE="$CRON_SCHEDULE $CRON_CMD $RENEW_CRON_TAG"

    if ! crontab -l 2>/dev/null | grep -Fq "$RENEW_CRON_TAG"; then
        (crontab -l 2>/dev/null; echo "$CRON_LINE") | crontab -
    fi
}
stop_proxy() {

    docker stop nginx

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
        docker run --name ohsal --network "$NETWORK_NAME" ohsal.com:latest
    fi
}

stop_ohsal() {
    if docker ps | grep -Fq ohsal; then
        docker stop ohsal
    fi
}


##################################################
#############   MAIN    ##########################
##################################################
if [[ "$mode" == "run" ]]; then
    create_network();
    start_proxy();
    start_ohsal();

elif [[ "$mode" == "stop" ]]; then
    stop_proxy();
    stop_ohsal();
else
    echo "Unknown Mode"
fi
