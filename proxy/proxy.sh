#!/bin/bash

mode=${1:-run}

# Cron setup
CERTBOT_CONF="$(pwd)/certbot/conf"
CERTBOT_WWW="$(pwd)/certbot/www"
LOGFILE="$(pwd)/certbot.log"
JOB="docker run --rm -v $CERTBOT_CONF:/etc/letsencrypt -v $CERTBOT_WWW:/var/www/certbot certbot/certbot renew --webroot -w /var/www/certbot && docker exec nginx nginx -s reload >> $LOGFILE 2>&1"

if [[ "$mode" == "run" ]]; then

    if ! docker ps -q -f name=nginx >/dev/null; then
        docker run -d \
            --name nginx \
            -p 80:80 -p 443:443 \
            -v "$(pwd)/nginx/conf.d:/etc/nginx/conf.d" \
            -v "$CERTBOT_CONF:/etc/letsencrypt" \
            -v "$CERTBOT_WWW:/var/www/certbot" \
            --restart unless-stopped \
            nginx:latest
    else
        echo "Nginx container already running"
    fi


    docker run --rm -it \
        -v "$CERTBOT_CONF:/etc/letsencrypt" \
        -v "$CERTBOT_WWW:/var/www/certbot" \
        certbot/certbot certonly --webroot -w /var/www/certbot \
        "${domains[@]/#/-d }"


    if crontab -l 2>/dev/null | grep -Fq "$JOB"; then
        echo "Cron job already exists"
    else
        (crontab -l 2>/dev/null; echo "0 */12 * * * $JOB") | crontab -
        echo "Cron job installed"
    fi

elif [[ "$mode" == "stop" ]]; then

    docker stop nginx

    if crontab -l 2>/dev/null | grep -Fq "$JOB"; then
        # Remove the job by filtering it out
        crontab -l 2>/dev/null | grep -Fv "$JOB" | crontab -
        echo "Cron job removed"
    else
        echo "No matching cron job to remove"
    fi

else
    echo "Unknown Mode"
fi
