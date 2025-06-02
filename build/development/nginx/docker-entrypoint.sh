#!/bin/sh
set -e

envsubst "\$APP_INTERNAL_PORT" < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

exec nginx -g "daemon off;"