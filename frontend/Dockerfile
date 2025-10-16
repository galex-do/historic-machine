# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built app to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create cache directory for map tiles
RUN mkdir -p /var/cache/nginx/tiles && \
    chown -R nginx:nginx /var/cache/nginx/tiles && \
    chmod -R 755 /var/cache/nginx/tiles

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]