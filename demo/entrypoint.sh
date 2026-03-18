#!/bin/bash
set -e

# Initialize database if it doesn't exist or is empty
if [ ! -s /root/tree/squirrels.db ]; then
    echo "Initializing squirrels database..."
    cat /root/tree/squirrels.sql | sqlite3 /root/tree/squirrels.db
    echo "Database initialized!"
fi

# Copy templates to workspace if empty or missing critical files
if [ ! -f /home/scrat/tree/squirrels.db ] || [ ! -f /home/scrat/tree/squix-tutorial.md ]; then
    echo "Initializing workspace..."
    cp -r /root/tree-templates/* /home/scrat/tree/
    chown -R scrat:scrat /home/scrat/tree
fi

# Fix database and directory permissions for scrat user (after mount)
# Make directory writable so SQLite can create WAL/journal files
chmod 777 /home/scrat/tree
if [ -f /home/scrat/tree/squirrels.db ]; then
    chmod 666 /home/scrat/tree/squirrels.db
    echo "Fixed database and directory permissions for writable access"
fi

# Initialize squix config for scrat user
mkdir -p /home/scrat/.config/squix
cp /root/.config/squix/config.yaml /home/scrat/.config/squix/config.yaml
chown -R scrat:scrat /home/scrat/.config

# Start cron daemon for periodic database resets
echo "Starting cron daemon for hourly database resets..."
service cron start
echo "Cron daemon started. Database will reset every hour."

# Run the provided command
exec "$@"
