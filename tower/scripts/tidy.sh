#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Tidying Go modules...${NC}"

# Run go mod tidy to update go.mod and go.sum
go mod tidy

# Verify the modules
echo -e "${YELLOW}Verifying Go modules...${NC}"
go mod verify

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Go modules are properly tidied and verified!${NC}"
else
  echo -e "${RED}There was an issue with the Go modules. Please check the errors above.${NC}"
  exit 1
fi
