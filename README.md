# CreatorSync

CreatorSync is a comprehensive content management platform designed specifically for streamers. It streamlines the process of managing, scheduling, and analyzing content across multiple platforms, helping creators focus on what they do best - creating amazing content.

## 🚀 Mission

Our mission is to empower creators by providing intuitive tools that simplify content management, enhance audience engagement, and provide actionable insights to grow their digital presence.

## ✨ Features

- **Unified Dashboard** - Manage all your content from one place
- **Cross-Platform Publishing** - Schedule and publish to multiple platforms
- **Audience Analytics** - Track performance and engagement metrics
- **Content Calendar** - Visual planning and scheduling
- **Collaboration Tools** - Work with teams and assistants seamlessly

## 🛠 Tech Stack

### Frontend

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4 with CSS Variables
- **UI Components**: shadcn/ui
- **Authentication**: Clerk
- **Analytics**: PostHog
- **State Management**: React Server Components + Context API

### Backend

- **Language**: Go (Golang)
- **Web Framework**: Fiber
- **Database**: PostgreSQL
- **ORM**: sqlc
- **Email**: Resend

### Infrastructure

- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Hosting**: TBD (Vercel/Netlify for frontend, AWS/GCP for backend)

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (for containerized development)

### Local Development

1. **Clone the repository**

   ```bash
   git clone git@github.com:baldybuilds/creatorsync.git
   cd creatorsync
   ```

2. **Set up environment variables**

   - Copy `.env.example` to `.env`
   - Update the environment variables with your local configuration

3. **Start the development environment**

   ```bash
   # Start the database
   make docker-run

   # Install frontend dependencies
   cd frontend
   npm install

   # Start the development servers (in separate terminal windows)
   make run  # Starts both backend and frontend with hot reloading
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## 📦 Available Scripts

### Development

```bash
# Run all tests
make test

# Run integration tests
make itest

# Start development servers with hot reloading
make run

# Watch for file changes and rebuild
make watch
```

### Building

```bash
# Build both frontend and backend
make build

# Build and run production build
make prod
```

### Database

```bash
# Start database container
make docker-run

# Stop database container
make docker-down

# Run database migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) to get started.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📬 Contact

Have questions or feedback? Open an issue or reach out to us at [your-email@example.com](mailto:your-email@example.com)

## 🏗️ **Environment Architecture**

### **Local Development**
- **Purpose**: Independent development environment
- **Database**: Local PostgreSQL instance
- **Clerk**: Development environment with test users
- **Domain**: `localhost:3000` (frontend) + `localhost:8080` (backend)

### **Staging (QA Environment)**
- **Purpose**: Production mirror for quality assurance testing
- **Database**: **Production database** (shared with production)
- **Clerk**: **Production Clerk environment** (same user IDs as production)
- **Domain**: `dev.creatorsync.app` (frontend) + `api-dev.creatorsync.app` (backend)
- **Users**: Real production users can test new features before release

### **Production**
- **Purpose**: Live environment for end users
- **Database**: Production PostgreSQL database
- **Clerk**: Production Clerk environment
- **Domain**: `creatorsync.app` (frontend) + `api.creatorsync.app` (backend)

## ✅ **Benefits of This Architecture**

1. **True QA Testing**: Staging uses identical data and users as production
2. **No Cross-Environment Issues**: Staging Clerk user IDs match database user IDs
3. **Real User Testing**: Production users can safely test new features on staging
4. **Simplified Deployment**: No complex data syncing between environments
5. **Consistent User Experience**: Same OAuth tokens and user profiles across environments

## 🚀 **Deployment**

### Staging Deployment
```bash
git push origin staging
```
- Deploys to Railway staging service with production database connection
- Uses production Clerk environment for authentication
- Available at staging domains for QA testing

### Production Deployment
```bash
git push origin production
```
- Deploys to Railway production service
- Uses production database and Clerk environment
- Available at production domains for end users
