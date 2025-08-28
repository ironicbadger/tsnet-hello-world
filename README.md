# Tailscale tsnet Demo Application

A minimal Go application demonstrating Tailscale's tsnet functionality. This app creates its own node on your tailnet and identifies connecting users via the id headers that they present as part of a web request.

## Usage

1. Generate an Auth Key at: https://login.tailscale.com/admin/settings/keys

2. Rename `.env.example` to `.env` and insert the Auth Key

3. Build the app with `docker compose up --build`

4. The app will be accessible at `https://ts-hello-world.your-tailnet.ts.net`.

## Customization

| Variable | Description | Default |
|----------|-------------|---------|
| `TS_AUTHKEY` | Tailscale authentication key | (required) |
| `TS_HOSTNAME` | Hostname for the tsnet node | ts-hello-world |
| `TS_STATE_DIR` | Directory for persistent state | /var/lib/ts-hello-world |
| `TS_ENABLE_SERVE` | Enable Tailscale Serve | true |
| `TS_ENABLE_FUNNEL` | Enable public access via Funnel | false |