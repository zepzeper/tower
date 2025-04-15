#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Make script exit when a command fails
set -e

echo -e "${YELLOW}Setting up development environment...${NC}"

# Verify that .air.toml exists
if [ ! -f ".air.toml" ]; then
  echo -e "${RED}Error: .air.toml file not found in the project root.${NC}"
  exit 1
fi

# Create tmp directory for air if it doesn't exist
if [ ! -d "tmp" ]; then
  echo -e "${YELLOW}Creating tmp directory for Air...${NC}"
  mkdir -p tmp
fi

echo -e "${YELLOW}Stopping any existing containers...${NC}"
docker-compose -f docker-compose.dev.yml down

echo -e "${YELLOW}Building and starting the development environment...${NC}"
docker-compose -f docker-compose.dev.yml up --build -d

echo -e "${GREEN}Development environment is up and running!${NC}"
echo -e "${YELLOW}Following logs from API container...${NC}"

# Follow logs
docker-compose -f docker-compose.dev.yml logs -f api
