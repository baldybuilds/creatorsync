# CreatorSync

CreatorSync is a comprehensive content management platform designed specifically for digital creators. It streamlines the process of managing, scheduling, and analyzing content across multiple platforms, helping creators focus on what they do best - creating amazing content.

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
