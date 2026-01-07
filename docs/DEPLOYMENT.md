# Polygame Deployment Guide

This guide provides instructions for deploying the Polygame platform.

## 1. Prerequisites

- **Server**: A Linux server (Ubuntu 22.04 recommended).
- **Docker**: Docker and Docker Compose installed.
- **Domain**: A registered domain name.
- **PostgreSQL**: A PostgreSQL database (can be run in Docker).
- **Redis**: A Redis instance (can be run in Docker).

## 2. Backend Deployment (Golang)

### Step 1: Build the Binary

On your local machine or a build server:

```bash
# Navigate to the backend directory
cd backend

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o polygame-server cmd/server/main.go
```

### Step 2: Create a Dockerfile

Create a `Dockerfile` in the `backend` directory:

```dockerfile
# Use a minimal base image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary
COPY polygame-server .

# Copy environment file
COPY .env.production .env

# Expose the port
EXPOSE 8080

# Run the server
CMD ["./polygame-server"]
```

### Step 3: Build and Push the Docker Image

```bash
docker build -t your-docker-repo/polygame-backend:latest .
docker push your-docker-repo/polygame-backend:latest
```

## 3. Frontend Deployment (Vue)

### Step 1: Build the Static Files

```bash
# Navigate to the frontend directory
cd frontend

# Install dependencies
npm install

# Build for production
npm run build
```

This will create a `dist` directory with the static files.

### Step 2: Configure Nginx

Create an Nginx configuration file (`/etc/nginx/sites-available/polygame`):

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /var/www/polygame/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Step 3: Upload Files and Enable Site

- Upload the `dist` directory to `/var/www/polygame/frontend`.
- Enable the Nginx site: `sudo ln -s /etc/nginx/sites-available/polygame /etc/nginx/sites-enabled/`
- Reload Nginx: `sudo systemctl reload nginx`

## 4. Admin Dashboard Deployment

Follow the same steps as the frontend deployment, but use the `admin` directory and a different Nginx configuration (e.g., on a subdomain like `admin.your-domain.com`).

## 5. Docker Compose Setup

For a complete, containerized deployment, use Docker Compose.

Create a `docker-compose.yml` in the project root:

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: polygame-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_db_password
      POSTGRES_DB: polygame
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: polygame-redis
    restart: unless-stopped

  backend:
    image: your-docker-repo/polygame-backend:latest
    container_name: polygame-backend
    depends_on:
      - db
      - redis
    environment:
      - GIN_MODE=release
      - DB_HOST=db
      - DB_USER=postgres
      - DB_PASSWORD=your_db_password
      - REDIS_HOST=redis
    ports:
      - "8080:8080"
    restart: unless-stopped

  frontend:
    image: nginx:alpine
    container_name: polygame-frontend
    depends_on:
      - backend
    volumes:
      - ./frontend/dist:/usr/share/nginx/html
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
    ports:
      - "80:80"
    restart: unless-stopped

volumes:
  postgres_data:
```

### Running with Docker Compose

1. Build the backend image.
2. Build the frontend static files.
3. Create the `nginx.conf` file.
4. Run `docker-compose up -d`.

## 6. Mobile App Deployment (Flutter)

### Android

1.  **Build the App Bundle**:
    ```bash
    cd mobile
    flutter build appbundle
    ```
2.  **Upload to Google Play Store**: Upload the generated `.aab` file (`build/app/outputs/bundle/release/app-release.aab`) to the Google Play Console.

### iOS

1.  **Build the Archive**:
    ```bash
    cd mobile
    flutter build ios
    ```
2.  **Archive and Upload**: Open the `ios` folder in Xcode, archive the project, and upload it to App Store Connect using Xcode's Organizer.

## 7. Final Steps

- **DNS**: Point your domain to the server's IP address.
- **SSL**: Configure SSL using Let's Encrypt (`certbot`).
- **Monitoring**: Set up monitoring for your services (e.g., Prometheus, Grafana).
