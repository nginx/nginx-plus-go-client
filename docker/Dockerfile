# syntax=docker/dockerfile:1.17
FROM debian:12-slim

LABEL maintainer="NGINX Docker Maintainers <integrations@nginx.com>"

ARG NGINX_PLUS_VERSION=R32

# Install NGINX Plus
# Download certificate and key from the customer portal (https://my.f5.com)
# and copy to the build context
RUN --mount=type=secret,id=nginx-repo.crt,dst=/etc/ssl/nginx/nginx-repo.crt,mode=0644 \
    --mount=type=secret,id=nginx-repo.key,dst=/etc/ssl/nginx/nginx-repo.key,mode=0644 \
    <<"eot" bash -euo pipefail
    apt-get update
    apt-get install --no-install-recommends --no-install-suggests -y ca-certificates gnupg curl apt-transport-https
    curl -fsSL https://cs.nginx.com/static/keys/nginx_signing.key | gpg --dearmor > /etc/apt/trusted.gpg.d/nginx_signing.gpg
    curl -fsSL -o /etc/apt/apt.conf.d/90pkgs-nginx https://cs.nginx.com/static/files/90pkgs-nginx
    DEBIAN_VERSION=$(awk -F '=' '/^VERSION_CODENAME=/ {print $2}' /etc/os-release)
    printf "%s\n" "deb https://pkgs.nginx.com/plus/${NGINX_PLUS_VERSION}/debian ${DEBIAN_VERSION} nginx-plus" > /etc/apt/sources.list.d/nginx-plus.list
    apt-get update
    apt-get install -y nginx-plus
    apt-get remove --purge --auto-remove -y gnupg
    rm -rf /var/lib/apt/lists/*
    rm /etc/apt/apt.conf.d/90pkgs-nginx /etc/apt/sources.list.d/nginx-plus.list
eot

EXPOSE 8080 8081

STOPSIGNAL SIGQUIT

RUN rm -rf /etc/nginx/conf.d/*
COPY --link docker/test.conf /etc/nginx/conf.d/

CMD ["nginx", "-g", "daemon off;"]
