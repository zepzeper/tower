#!/bin/sh
# Start Vite in the background
npm install -D tailwindcss postcss autoprefixer
cd /app/ui && npm run tailwind &
TAILWIND_PID=$!
# Start Vite in the background
cd /app/ui && npm run dev &
VITE_PID=$!

# Start Air in the foreground
cd /app && air -c .air.toml &
AIR_PID=$!

# Handle termination
trap "kill $TAILWIND_PID $VITE_PID $AIR_PID; exit" SIGINT SIGTERM

# Keep the container running
wait
