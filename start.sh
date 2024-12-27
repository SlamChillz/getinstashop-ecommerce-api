#!/bin/sh

set -e

# Create Admin User
./main -email="admin@gmail.com" -password="adminpassword"

echo "start application"
exec "$@"
