# Multi-stage build for React application
FROM node:20-alpine AS base

# Install dependencies only when needed
FROM base AS deps
WORKDIR /app

# Copy package files
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm ci

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app

# Accept build arguments
ARG VITE_ROX_API_KEY
ARG NGINX_CONFIG=nginx.conf
ENV NGINX_CONFIG=${NGINX_CONFIG}
# Pass FM key to Vite build (optional - app works without it)
ENV VITE_ROX_API_KEY=${VITE_ROX_API_KEY}

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build with default BASE_PATH=/
# Runtime base path detection will handle environment-specific paths
RUN npm run build

# Production image with nginx
FROM nginx:alpine AS runner
WORKDIR /app

# Accept nginx config arg in runner stage
ARG NGINX_CONFIG=nginx.conf

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration (use build arg to choose which config)
COPY ${NGINX_CONFIG} /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
