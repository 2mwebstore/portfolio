#!/usr/bin/env bash
# One-time setup for auto-deploy from GitHub (.github/workflows/deploy.yml).
# Creates an SSH key that can only run deploy/update.sh on this Droplet and
# prints the three values to save as GitHub repo secrets.
# Run as root from the repo folder:  bash deploy/setup-auto-deploy.sh
set -euo pipefail
cd "$(dirname "$0")/.."

tmp="$(mktemp -d)"
ssh-keygen -q -t ed25519 -N "" -C "github-actions-deploy" -f "$tmp/key"

# "restrict" blocks shells, port forwarding and the like: whatever command
# comes in with this key, the Droplet only runs the deploy script.
install -m 700 -d ~/.ssh
echo "command=\"cd $(pwd) && bash deploy/update.sh\",restrict $(cat "$tmp/key.pub")" >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

ip="$(curl -sf -m 3 http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address \
  || hostname -I | awk '{print $1}')"

echo "Add these 3 secrets on GitHub: repo -> Settings -> Secrets and variables"
echo "-> Actions -> New repository secret"
echo
echo "===== DROPLET_HOST ====="
echo "$ip"
echo
echo "===== DROPLET_KNOWN_HOSTS ====="
echo "$ip $(cut -d' ' -f1,2 /etc/ssh/ssh_host_ed25519_key.pub)"
echo
echo "===== DROPLET_SSH_KEY (copy every line, including BEGIN and END) ====="
cat "$tmp/key"
rm -rf "$tmp"
