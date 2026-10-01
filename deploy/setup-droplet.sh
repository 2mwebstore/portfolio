#!/usr/bin/env bash
# One-time setup for a fresh Ubuntu Droplet: installs Docker, nginx and
# certbot, opens the firewall for SSH/HTTP/HTTPS, and adds 2 GB of swap so
# the image can build on a 1 GB Droplet. Then add each client with
# deploy/add-site.sh. Run as root:  bash deploy/setup-droplet.sh
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v docker >/dev/null; then
  curl -fsSL https://get.docker.com | sh
fi

apt-get update
apt-get install -y nginx certbot python3-certbot-nginx
rm -f /etc/nginx/sites-enabled/default
cp deploy/nginx-catch-all.conf /etc/nginx/sites-available/catch-all
ln -sf /etc/nginx/sites-available/catch-all /etc/nginx/sites-enabled/catch-all
nginx -t
systemctl reload nginx

ufw allow OpenSSH
ufw allow 'Nginx Full'
ufw --force enable

if [ ! -f /swapfile ]; then
  fallocate -l 2G /swapfile
  chmod 600 /swapfile
  mkswap /swapfile
  swapon /swapfile
  echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

echo "Droplet ready. Next, add a client site:"
echo "  bash deploy/add-site.sh <name> <domain> [more domains...]"
