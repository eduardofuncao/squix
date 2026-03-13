#!/bin/bash
# Initialize database if it doesn't exist or is empty
if [ ! -s /root/tree/squirrels.db ]; then
    echo "Initializing squirrels database..."
    cat /root/tree/squirrels.sql | sqlite3 /root/tree/squirrels.db
    echo "Database initialized!"
fi

# Run the provided command (gotty)
exec "$@"
