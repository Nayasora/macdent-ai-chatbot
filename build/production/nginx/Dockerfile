FROM nginx:alpine

RUN apk add --no-cache openssl
RUN mkdir -p /etc/nginx/ssl2

RUN openssl dhparam -out /etc/nginx/dhparam.pem 2048

COPY build/development/nginx/docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

CMD ["/docker-entrypoint.sh"]