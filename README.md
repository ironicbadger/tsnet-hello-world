# Tailscale tsnet Demo Application

A minimal Go application demonstrating Tailscale's tsnet functionality. This app creates its own node on your tailnet and identifies connecting users.

## Features

- **User Identification**: Displays connecting user's identity via Tailscale
- **Service Status**: JSON API endpoint for service status
- **Serve/Funnel Support**: Configurable access modes via environment variables
- **Minimal Design**: Clean, production-ready code following Go best practices

## Tech Stack

- **Language**: Go 1.23
- **Library**: tailscale.com/tsnet
- **Container**: Multi-stage Docker build with Alpine Linux

## Quick Start

### Prerequisites

1. Docker and Docker Compose installed
2. A Tailscale account
3. An auth key from https://login.tailscale.com/admin/settings/keys

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd ts-hello-world
```

2. Create a `.env` file with your Tailscale auth key:
```bash
echo "TS_AUTHKEY=tskey-auth-xxxxxxxxxxxxx" > .env
```

### Running the Application

```bash
docker compose up --build
```

The app will be accessible at `https://ts-hello-world.your-tailnet.ts.net`

To enable public access via Funnel:
```bash
TS_ENABLE_FUNNEL=true docker compose up --build
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TS_AUTHKEY` | Tailscale authentication key | (required) |
| `TS_HOSTNAME` | Hostname for the tsnet node | ts-hello-world |
| `TS_STATE_DIR` | Directory for persistent state | /var/lib/ts-hello-world |
| `TS_ENABLE_SERVE` | Enable Tailscale Serve | true |
| `TS_ENABLE_FUNNEL` | Enable public access via Funnel | false |

## API Endpoints

- `GET /` - Web interface showing user identity and connection info
- `GET /api/status` - JSON API returning service status

## Project Structure

```
ts-hello-world/
├── main.go           # Main application (165 lines)
├── templates/
│   └── index.html    # HTML template
├── static/
│   └── style.css     # Stylesheet
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
├── Dockerfile        # Multi-stage production build
├── docker-compose.yml # Container orchestration
└── README.md         # This file
```

## Building Locally

```bash
# Build Docker image
docker build -t ts-hello-world .

# Run container
docker run -d \
  --name ts-hello-world \
  -e TS_AUTHKEY=tskey-auth-xxx \
  -v ts-hello-world-state:/var/lib/ts-hello-world \
  ts-hello-world
```

## Security Notes

- Auth keys should be treated as secrets and never committed to version control
- Use `.env` files for local development only
- In production, use proper secret management
- Funnel mode exposes your service to the public internet - use with caution

## License

MIT