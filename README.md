# Car Rental Directory

A directory site listing multiple car rental companies — Vue 3 frontend, Go (Gin + GORM) backend, SQLite database.

## Structure

```
car-rental-directory/
  backend/     Go API + SQLite
  frontend/    Vue 3 + Vite
```

## Backend setup

```bash
cd backend
go mod tidy      # downloads gin, gorm, sqlite driver, cors
go run .
```

This starts the API on `http://localhost:8080` and creates `car_rental.db` in the
`backend/` folder on first run, seeded with placeholder partners/config/icons.

Endpoints:

| Method | Path                    | Auth required | Purpose                     |
|--------|-------------------------|:--------------:|------------------------------|
| POST   | /api/auth/login         | no  | Log in, returns a JWT         |
| PUT    | /api/auth/password      | yes | Change the admin password     |
| GET    | /api/config             | no  | Site name, logo, banner, socials |
| PUT    | /api/config             | yes | Update site config           |
| GET    | /api/partners           | no  | List rental partners         |
| POST   | /api/partners           | yes | Add a partner                |
| PUT    | /api/partners/:id       | yes | Update a partner              |
| DELETE | /api/partners/:id       | yes | Remove a partner              |
| GET    | /api/payment-icons      | no  | List bank/social icons        |
| POST   | /api/payment-icons      | yes | Add an icon                   |
| PUT    | /api/payment-icons/:id  | yes | Update an icon                 |
| DELETE | /api/payment-icons/:id  | yes | Remove an icon                 |
| POST   | /api/upload             | yes | Upload an image file, returns `{ url }` |

Auth-required routes expect `Authorization: Bearer <token>` from the login response.

### Image uploads

`POST /api/upload` takes a multipart form with a `file` field (png/jpg/jpeg/gif/webp),
saves it into `backend/uploads/`, and returns `{ "url": "/uploads/169..." }`. The
admin UI's file inputs (Partners → Logo / Media image) call this automatically and
fill in the URL field for you — no need to host images elsewhere.

### Default admin account

On first run, the server creates one admin user and prints the password to the
console:

```
======================================================
 Created default admin account:
   username: admin
   password: 123123123
 Log in and change this password right away.
======================================================
```

Override it before first run with an env var so it's never the same default twice:

```bash
ADMIN_DEFAULT_PASSWORD="something-only-you-know" go run .
```

Also set `JWT_SECRET` to a real random string in production (the code falls back to
a dev-only value otherwise):

```bash
JWT_SECRET="$(openssl rand -hex 32)" go run .
```

**Log in and change the password immediately** via the admin UI's Password tab
(or `PUT /api/auth/password`) — this is a single-admin setup, so there's no
separate user-management screen.

Images: drop files into `backend/uploads/` and reference them as `/uploads/filename.png`
in `logo_url` / `icon_url` / `banner_image_url` fields — the server serves that folder
as static files.

## Frontend setup

```bash
cd frontend
npm install
npm run dev
```

Vite proxies `/api` and `/uploads` to `http://localhost:8080` in dev
(see `vite.config.js`), so just run the backend alongside it.

Visit `/admin/login` to sign in, or `/` for the public storefront.

For production, `npm run build` outputs static files in `frontend/dist/`
that you can serve with Nginx, pointing `/api` and `/uploads` at your
Go backend the same way. Since this uses Vue Router in history mode
(clean URLs like `/admin`, no `#`), add an SPA fallback to your Nginx
config so deep links work on refresh:

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

## Replacing placeholder content

Everything (site name, logo, banner, partner list, bank/social icons) lives in the
SQLite database, not hardcoded in the frontend. Update it via the API
(e.g. `curl -X PUT http://localhost:8080/api/config -d '{...}'`) or the admin UI at
`/admin`.

## Recent additions

- **Bank/Social icons** are now fully editable (not just add/delete) from the
  Icons tab in the admin panel.
- **Telegram/Facebook buttons** can each have a custom icon image
  (`telegram_icon_url` / `facebook_icon_url`) and background color
  (`telegram_color` / `facebook_color`) — set from the Site Config tab. Leave the
  icon URL blank to fall back to the default text glyph.
- **Hero banner** is now a fixed 300px tall, image cropped with `object-fit: cover`
  so banners of any aspect ratio look consistent.
- **Footer note** now scrolls as a marquee directly under the site name/logo in
  the footer, and respects `prefers-reduced-motion` (falls back to static
  centered text for people with that setting on).

If you're updating an existing `car_rental.db`, GORM's `AutoMigrate` adds the new
`SiteConfig` columns automatically on next `go run .` — no manual migration
needed.

## Deploying to a DigitalOcean Droplet (public hosting)

One Droplet hosts every client from this one codebase. Each client is a
separate **site** with its own domain(s), database, uploaded images and admin
login, so editing one client's content never affects another. For example:

| Site    | Domains                  | Data folder                |
|---------|--------------------------|----------------------------|
| `net97` | net97.co, www.net97.co   | `/opt/net97.co/data/net97` |

How it fits together: the repo builds one Docker image, `net97co-app` (the Go
API serves `/api`, `/uploads` and the built Vue frontend). Each site runs its
own container from that image on a private port (`127.0.0.1:8101`, `8102`, …),
using `docker-compose.yml` with the site's settings from `sites/<name>.env`.
nginx is the public web server and sends each domain to its site's port
(`deploy/nginx.conf` is the template). Certbot provides free HTTPS
certificates and renews them automatically.

**Sharing a Droplet with bp24 (or another project):** this project can run on
the same Droplet as bp24. Its image (`net97co-app`), containers and nginx
files (`net97co-<name>`) and ports (from 8101; bp24 uses 8081 up) are all
separate, so the two projects never overwrite each other's sites, even when a
site has the same name in both. Each project deploys only its own sites.

1. **Create a Droplet** on [digitalocean.com](https://www.digitalocean.com):
   Ubuntu 24.04. 1 GB of RAM ($6/month) runs several sites, since each one
   uses little memory. Add your SSH key. Skip this if you're using the
   Droplet that already runs bp24.
2. **Point every client domain at it**: for each domain, create a DNS `A`
   record (and one for `www` if you want it) with the Droplet's IP address.
3. **SSH in and get the code**. The repo is private, so give the Droplet
   read-only access with a deploy key for this repo. GitHub allows each deploy
   key on only one repo, so this one gets its own key even if the Droplet
   already pulls bp24:
   ```bash
   ssh root@YOUR_DROPLET_IP
   ssh-keygen -t ed25519 -N "" -f ~/.ssh/net97co_deploy
   cat ~/.ssh/net97co_deploy.pub
   ```
   On GitHub, open the net97.co repo → Settings → Deploy keys → Add deploy
   key, paste the printed line, and leave "Allow write access" off. Then tell
   SSH to use that key for this repo, and clone:
   ```bash
   cat >> ~/.ssh/config <<'EOF'
   Host github-net97co
     HostName github.com
     IdentityFile ~/.ssh/net97co_deploy
     IdentitiesOnly yes
   EOF
   ssh-keyscan github.com >> ~/.ssh/known_hosts
   git clone git@github-net97co:2mwebstore/net97.co.git /opt/net97.co
   cd /opt/net97.co
   ```
4. **Set up the server** (one time per Droplet). This installs Docker, nginx
   and certbot, opens the firewall for SSH/HTTP/HTTPS and adds swap. If the
   Droplet already runs bp24 this was done already, but running it again is
   harmless:
   ```bash
   bash deploy/setup-droplet.sh
   ```
5. **Add each client site**: give it a short name, then its domains:
   ```bash
   bash deploy/add-site.sh net97 net97.co www.net97.co
   ```
   Each run creates the site's settings in `sites/<name>.env` (a random
   `JWT_SECRET`, the admin password `123123123`, and its private port), sets
   up the nginx site, and starts it. It prints the admin password and the `certbot`
   command for the next step. The first run builds the image, which takes a
   few minutes. Add more clients the same way. To give two domains different
   content, add them as two sites instead of one.
6. **Turn on HTTPS** for each site once its domain loads over HTTP, using the
   command `add-site.sh` printed:
   ```bash
   certbot --nginx -d net97.co -d www.net97.co
   ```
   Certbot asks for an email the first time, adds the certificate to that
   site's nginx config, redirects HTTP to HTTPS, and renews automatically.
7. For each site, open `https://<domain>/admin/login`, log in as `admin` with
   the printed password, and change it from the Password tab. Every new site
   starts with the placeholder content; replace it in that site's admin panel.
8. **Turn on auto-deploy** (one time only), from `/opt/net97.co`:
   ```bash
   bash deploy/setup-auto-deploy.sh
   ```
   It prints three values. On GitHub, open the net97.co repo → Settings →
   Secrets and variables → Actions → New repository secret, and add each one
   with the name it's printed under: `DROPLET_HOST`, `DROPLET_KNOWN_HOSTS` and
   `DROPLET_SSH_KEY`. The key it creates can only run this project's
   `deploy/update.sh` on the Droplet, so it can't be used to log in.

**Deploying changes:** push to `main` on GitHub. The "Deploy to Droplet"
workflow (`.github/workflows/deploy.yml`) connects to the Droplet and runs
`deploy/update.sh`, which pulls, rebuilds the image once and restarts every
site on the new version. Follow it in the repo's Actions tab: a red ❌ means
the deploy failed, and its log shows why. To redeploy without a push, use
Actions → Deploy to Droplet → Run workflow, or run
`cd /opt/net97.co && bash deploy/update.sh` on the Droplet. Don't edit code
directly on the Droplet, or the next `git pull` will fail.

**Logs:** `docker logs -f net97co-net97-app-1` for a site (replace the middle
`net97` with the site name), `/var/log/nginx/error.log` for nginx.

**Removing a site:**
```bash
docker compose -p net97co-net97 --env-file sites/net97.env down
rm /etc/nginx/sites-enabled/net97co-net97 && systemctl reload nginx
```
Its data stays in `data/net97` until you delete that folder.

**Upload size:** nginx accepts uploads up to 20 MB (`client_max_body_size` in
`deploy/nginx.conf`). The live copy for each site is
`/etc/nginx/sites-available/net97co-<name>`. Edit that copy on the Droplet,
then run `nginx -t && systemctl reload nginx`.

**Data and backups:** each site's SQLite database and uploaded images are in
`/opt/net97.co/data/<name>` on the Droplet (mounted at `/app/data` in its
container, where `DATA_DIR` points). Rebuilds and restarts never touch them.
Back up every site at once by copying the `data` folder, e.g. from your own
machine: `scp -r root@YOUR_DROPLET_IP:/opt/net97.co/data ./backup`, or turn on
DigitalOcean's Droplet backups. `sites/` holds each site's secrets and first
admin password, so keep it private.

## Partner background photo + upload

Partners now have a `media_image_url` field — a full background photo shown
behind the logo badge on the card (previously it was just a gradient). Set it via
the admin UI's file upload (Partners tab → "Media image") or paste a URL directly.
On hover, the photo zooms with a bouncy `cubic-bezier(0, .2, .5, 3)` easing over
.5s. A dark gradient overlay keeps the partner name readable regardless of the
photo's brightness.

git add -A
git commit -m "Initial commit: car rental directory ready for Railway"
git branch -M main
git remote add origin https://github.com/chansila5555-oss/vp168.git
git push -u origin main