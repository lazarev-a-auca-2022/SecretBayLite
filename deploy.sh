#!/bin/bash

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create necessary directories
mkdir -p logs

# Create directories for SSL certificates
mkdir -p certbot/www
mkdir -p certbot/conf

# Set permissions
chmod 755 logs
chmod -R 755 certbot

# Start frontend first in HTTP mode for certificate generation
echo "Starting frontend service for certificate generation..."
docker-compose up -d frontend

# Wait for frontend to be ready
echo "Waiting for frontend to start..."
sleep 10

# Generate SSL certificates
echo "Generating SSL certificates..."
docker-compose run --rm certbot certonly --webroot --webroot-path=/var/www/certbot \
    --email admin@example.com --agree-tos --no-eff-email -d localhost --staging --force-renewal

# Now start all services
echo "Starting all services..."
docker-compose up -d --force-recreate

# Set up automatic certificate renewal
echo "0 */12 * * * cd $(pwd) && docker-compose run --rm certbot renew && docker-compose exec frontend nginx -s reload" | sudo crontab -

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo "SecretBay has been successfully deployed!"
    echo "Frontend is available at: https://localhost"
    echo "Backend API is available at: https://localhost:8080"
    echo "SSL certificates have been initialized"
    echo "Automatic certificate renewal has been configured (every 12 hours)"
else
    echo "Deployment failed. Please check docker-compose logs for details."
    docker-compose logs
fi