# https://github.com/mattermost/mattermost-docker/blob/master/web/Dockerfile
FROM nginxinc/nginx-unprivileged:mainline-alpine

USER root

RUN mkdir -p /usr/local/bin && \
    cd /usr/local/bin && \
    wget https://github.com/gohugoio/hugo/releases/download/v0.91.2/hugo_0.91.2_Linux-64bit.tar.gz -O hugo.tar.gz && \
    tar xvf hugo.tar.gz && \
    rm hugo.tar.gz

RUN rm /etc/nginx/conf.d/default.conf \
    && apk add --no-cache curl
COPY ./html /synmodoro/
COPY entrypoint.sh /
COPY synmodoro /etc/nginx/sites-available/

RUN chown -R nginx:nginx /etc/nginx/sites-available && \
         chown -R nginx:nginx /var/cache/nginx && \
         chown -R nginx:nginx /var/log/nginx && \
         chown -R nginx:nginx /etc/nginx/conf.d && \
         chown nginx:nginx entrypoint.sh
RUN touch /var/run/nginx.pid && \
         chown -R nginx:nginx /var/run/nginx.pid

COPY ./security.conf /etc/nginx/conf.d/

RUN chown -R nginx:nginx /etc/nginx/conf.d/security.conf

RUN chmod u+x /entrypoint.sh

RUN sed -i "/^http {/a \    proxy_buffering off;\n" /etc/nginx/nginx.conf
RUN sed -i '/temp_path/d' /etc/nginx/nginx.conf \
    && sed -i 's!/tmp/nginx.pid!/var/run/nginx.pid!g' /etc/nginx/nginx.conf

RUN cd /synmodoro/synmodoro && \
    /usr/local/bin/hugo

USER nginx

ENTRYPOINT ["/entrypoint.sh"]
VOLUME ["/var/run", "/etc/nginx/conf.d/", "/var/cache/nginx/"]
