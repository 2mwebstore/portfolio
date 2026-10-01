#!/usr/bin/env bash
# Add a client site: its own database, uploads, admin login and domain(s),
# all running from this one codebase. Run as root from the repo folder:
#   bash deploy/add-site.sh net97 net97.co www.net97.co
set -euo pipefail
cd "$(dirname "$0")/.."

name="${1:-}"
if [[ ! "$name" =~ ^[a-z0-9][a-z0-9_-]*$ ]] || [ $# -lt 2 ]; then
  echo "usage: bash deploy/add-site.sh <name> <domain> [more domains...]" >&2
  echo "  <name> is lowercase letters, numbers, - or _ (e.g. net97)" >&2
  exit 1
fi
shift

env_file="sites/$name.env"
if [ -f "$env_file" ]; then
  echo "$env_file already exists, so site '$name' was already added" >&2
  exit 1
fi

# The container and nginx site are prefixed with this project's name, so a
# same-named site from another project on this Droplet (e.g. bp24) can't
# replace them.
project="net97co-$name"

# Each site gets its own local port. Starting at 8101 keeps clear of other
# projects on this Droplet (bp24 starts at 8081), and ports already in use
# are skipped.
port=8101
while grep -qsx "APP_PORT=$port" sites/*.env || ss -Hltn "sport = :$port" | grep -q .; do
  port=$((port + 1))
done

# First admin login for every new site. Change it from the admin panel's
# Password tab right after logging in.
password="123123123"
mkdir -p sites
cat > "$env_file" <<ENV
SITE=$name
APP_PORT=$port
JWT_SECRET=$(openssl rand -hex 32)
ADMIN_DEFAULT_PASSWORD=$password
ENV
chmod 600 "$env_file"

# Don't overwrite an existing nginx site: certbot adds its HTTPS settings to it.
site_conf="/etc/nginx/sites-available/$project"
if [ ! -f "$site_conf" ]; then
  sed -e "s/server_name DOMAIN;/server_name $*;/" \
      -e "s/127.0.0.1:APP_PORT;/127.0.0.1:$port;/" \
      deploy/nginx.conf > "$site_conf"
  ln -sf "$site_conf" "/etc/nginx/sites-enabled/$project"
else
  echo "$site_conf already exists, leaving it as is"
fi
nginx -t
systemctl reload nginx

docker image inspect net97co-app >/dev/null 2>&1 || docker build -t net97co-app .
docker compose -p "$project" --env-file "$env_file" up -d

echo
echo "Site '$name' is running at http://$1 (port $port)"
echo "Admin login: admin / $password  (saved in $env_file)"
echo "Once the domain's DNS points at this Droplet, turn on HTTPS with:"
echo "  certbot --nginx$(printf ' -d %s' "$@")"
