server {
    listen 80;
    server_name jolsho.com www.jolsho.com;

    # ACME challenge path for Certbot
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # Redirect all other HTTP requests to HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

http {
    # Shared memory zone for tracking client IPs
    limit_req_zone $binary_remote_addr zone=one:10m rate=5r/s;

    # Single certificate for all domains
    ssl_certificate     /etc/letsencrypt/live/common.crt;
    ssl_certificate_key /etc/letsencrypt/live/common.key;

    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    server {
        listen 443 ssl http2;
        server_name jolsho.com www.jolsho.com;

        location / {
            # Apply the rate limit
            limit_req zone=one burst=15 nodelay;
            limit_req_status 429;  # Too Many Requests if limit exceeded

            proxy_pass http://jolsho:8080;
            proxy_http_version 1.1;

            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";

            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            proxy_read_timeout 3600s;
            proxy_send_timeout 3600s;
        }

        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }
    }
}
