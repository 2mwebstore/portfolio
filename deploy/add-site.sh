#!/usr/bin/env bash
# Add a client site: its own database, uploads, admin login and domain(s),
# all running from this one codebase. Run as root from the repo folder:
#   bash deploy/add-site.sh lik lik.me www.lik.me
set -euo pipefail
cd "$(dirname "$0")/.."

name="${1:-}"
if [[ ! "$name" =~ ^[a-z0-9][a-z0-9_-]*$ ]] || [ $# -lt 2 ]; then
  echo "usage: bash deploy/add-site.sh <name> <domain> [more domains...]" >&2
  echo "  <name> is lowercase letters, numbers, - or _ (e.g. lik)" >&2
  exit 1
fi
shift

env_file="sites/$name.env"
if [ -f "$env_file" ]; then
  echo "$env_file already exists, so site '$name' was already added" >&2
  exit 1
fi

# Each site gets its own local port, starting at 8081.
port=8081
while grep -qsx "APP_PORT=$port" sites/*.env; do
  port=$((port + 1))
done

password="$(openssl rand -hex 8)"
mkdir -p sites
cat > "$env_file" <<ENV
SITE=$name
APP_PORT=$port
JWT_SECRET=$(openssl rand -hex 32)
ADMIN_DEFAULT_PASSWORD=$password
ENV
chmod 600 "$env_file"

# Don't overwrite an existing nginx site: certbot adds its HTTPS settings to it.
site_conf="/etc/nginx/sites-available/$name"
if [ ! -f "$site_conf" ]; then
  sed -e "s/server_name DOMAIN;/server_name $*;/" \
      -e "s/127.0.0.1:APP_PORT;/127.0.0.1:$port;/" \
      deploy/nginx.conf > "$site_conf"
  ln -sf "$site_conf" "/etc/nginx/sites-enabled/$name"
else
  echo "$site_conf already exists, leaving it as is"
fi
nginx -t
systemctl reload nginx

docker image inspect car-rental-app >/dev/null 2>&1 || docker build -t car-rental-app .
docker compose -p "$name" --env-file "$env_file" up -d

echo
echo "Site '$name' is running at http://$1 (port $port)"
echo "Admin login: admin / $password  (saved in $env_file)"
echo "Once the domain's DNS points at this Droplet, turn on HTTPS with:"
echo "  certbot --nginx$(printf ' -d %s' "$@")"
