# Squix Demo Environment

A Docker container with GoTTY web terminal, Squix, Presenterm, and a sample squirrel database.

## Quick Start

```bash
# Build and run
docker compose up --build

# Access at http://localhost:8082
```

Or run directly:
```bash
docker run -d --rm --name demo-gotty \
  -p 8082:8080 \
  -v ~/.config/squix:/root/.config/squix \
  -v ./workspace:/workspace \
  demo-gotty
```

## What's Included

- **GoTTY v1.6.0**: Web terminal access
- **Squix v0.3.0-beta**: SQL query manager
- **Presenterm v0.16.1**: Terminal presentation tool
- **Vim**: Text editor
- **SQLite3**: Database

## Files

### workspace/
- **squirrels.db**: SQLite database with 10 squirrels, 20 sightings
- **squirrels.sql**: SQL schema and data (source)
- **squix-tutorial.md**: Interactive Squix tutorial

## Usage

### Run the Tutorial Presentation

In the browser terminal:
```bash
cd /workspace
presenterm squix-tutorial.md
```

Use arrow keys to navigate, press Enter on code blocks to execute.

### Explore the Database

```bash
# Connect to database
squix init --name squirrels --type sqlite --conn "/workspace/squirrels.db"

# Run queries
squix query "SELECT * FROM squirrels"
squix query "SELECT * FROM sightings"

# Create saved queries
squix add "all-squirrels" "SELECT * FROM squirrels"
squix run "all-squirrels"
```

### Edit Files

```bash
vim /workspace/squix-tutorial.md
```

## Stopping

```bash
docker compose down
# or
docker stop demo-gotty
```

## For VPS Deployment

Add Traefik labels to `docker-compose.yml` for reverse proxy routing with TLS.
