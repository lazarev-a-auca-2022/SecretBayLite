# Stage 1: Build the Flutter web app
FROM ubuntu:22.04 AS builder

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -ms /bin/bash developer
USER developer
WORKDIR /home/developer

# Download and setup Flutter
RUN wget https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.16.9-stable.tar.xz \
    && tar xf flutter_linux_3.16.9-stable.tar.xz \
    && rm flutter_linux_3.16.9-stable.tar.xz

ENV PATH="/home/developer/flutter/bin:${PATH}"

# Configure Git for the Flutter installation
RUN git config --global --add safe.directory /home/developer/flutter

# Pre-download development binaries and precache Flutter
RUN flutter precache && \
    flutter doctor && \
    flutter config --enable-web

# Create and set the app directory
WORKDIR /home/developer/app

# Copy the app files with correct ownership
COPY --chown=developer:developer . .

# Get app dependencies and build
RUN flutter pub get && \
    flutter build web --release --web-renderer html

# Stage 2: Serve the app using nginx
FROM nginx:alpine

# Copy the built app
COPY --from=builder /home/developer/app/build/web /usr/share/nginx/html/

# Create required directories with proper permissions
RUN mkdir -p /etc/nginx/conf.d && \
    mkdir -p /etc/letsencrypt && \
    mkdir -p /var/www/certbot && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    chown -R nginx:nginx /var/www/certbot

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN chmod 644 /etc/nginx/conf.d/default.conf

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]