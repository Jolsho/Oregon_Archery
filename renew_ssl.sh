#!/usr/bin/env bash

set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
CERTBOT_CONF="$BASE_DIR/proxy/certbot/conf"
CERTBOT_WWW="$BASE_DIR/proxy/certbot/www"
LOGFILE="$BASE_DIR/proxy/certbot/certbot.log"

docker run --rm \
    -v $CERTBOT_CONF:/etc/letsencrypt \
    -v $CERTBOT_WWW:/var/www/certbot \
    certbot/certbot renew --webroot -w /var/www/certbot && \
    docker exec nginx nginx -s reload >> $LOGFILE 2>&1
