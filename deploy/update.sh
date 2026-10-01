#!/usr/bin/env bash
# Pull the latest code, rebuild the image once and restart every client site
# on it. Data in ./data is left untouched. Run on the Droplet:
#   bash deploy/update.sh
set -euo pipefail
cd "$(dirname "$0")/.."

# Pull, then run the rest from the freshly pulled copy of this script, so the
# build and restart steps always match the new docker-compose.yml.
if [ "${1:-}" != "--pulled" ]; then
  git pull --ff-only
  exec bash deploy/update.sh --pulled
fi

docker build -t portfolio-app .
for env_file in sites/*.env; do
  [ -f "$env_file" ] || continue
  docker compose -p "portfolio-$(basename "$env_file" .env)" --env-file "$env_file" up -d
done
docker image prune -f
