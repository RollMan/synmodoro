#!/bin/sh
# https://github.com/mattermost/mattermost-docker/blob/master/web/entrypoint.sh

APP_HOST=${APP_HOST:-app}
APP_PORT_NUMBER=${APP_PORT_NUMBER:-80}

test -w /etc/nginx/conf.d/synmodoro.conf && rm /etc/nginx/conf.d/synmodoro.conf
ln -s -f /etc/nginx/sites-available/synmodoro /etc/nginx/conf.d/synmodoro.conf

sed -i "s/{%APP_HOST%}/${APP_HOST}/g" /etc/nginx/conf.d/synmodoro.conf
sed -i "s/{%APP_PORT%}/${APP_PORT_NUMBER}/g" /etc/nginx/conf.d/synmodoro.conf

exec nginx -g 'daemon off;'
