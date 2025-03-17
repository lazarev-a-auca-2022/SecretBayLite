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

# Build and start the containers
echo "Building and starting containers..."
docker-compose up --build -d

# Wait for services to start
echo "Waiting for services to start..."
sleep 10

# Initialize SSL certificates
echo "Initializing SSL certificates..."
# Remove --staging flag after testing
docker-compose run --rm certbot certonly --webroot --webroot-path=/var/www/certbot \
    --email admin@example.com --agree-tos --no-eff-email -d localhost --staging --force-renewal

# Reload nginx to apply new certificates
docker-compose exec frontend nginx -s reload

# Set up automatic certificate renewal
echo "0 */12 * * * docker-compose -f $(pwd)/docker-compose.yml run --rm certbot renew" | sudo crontab -

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