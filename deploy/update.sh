#!/usr/bin/env bash
# Pull the latest code, rebuild the image once and restart every client site
# on it. Data in ./data is left untouched. Run on the Droplet:
#   bash deploy/update.sh
set -euo pipefail
cd "$(dirname "$0")/.."

git pull --ff-only
docker build -t car-rental-app .
for env_file in sites/*.env; do
  [ -f "$env_file" ] || continue
  docker compose -p "$(basename "$env_file" .env)" --env-file "$env_file" up -d
done
docker image prune -f
