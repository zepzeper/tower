#!/bin/sh
# Install dependencies if needed
cd /app/ui && npm install

# Build the static files first (creates dist folder)
cd /app/ui && npm run build

# Start Vite dev server in the background
cd /app/ui && npm run dev &
VITE_PID=$!

# Start Air in the foreground
cd /app && air -c .air.toml &
AIR_PID=$!

# Handle termination
trap "kill $VITE_PID $AIR_PID; exit" SIGINT SIGTERM

# Keep the container running
wait
