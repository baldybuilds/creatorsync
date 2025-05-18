# Backend Architecture

## 🏗️ Architecture Philosophy

The CreatorSync backend is designed around simplicity and scalability. By using Go + Fiber with Go-Blueprint scaffolding, we achieve a fast, lightweight API that can handle concurrent video processing requests while maintaining clean code organization.

## 🧱 Core Components

### Authentication Layer

Using Clerk's JWT validation removes the complexity of managing OAuth tokens manually. The middleware validates every protected request and extracts user context, allowing the rest of the application to focus on business logic rather than auth concerns.

### Service Architecture

The backend follows a service-oriented pattern where each major feature (Twitch integration, video rendering, file storage, server-side analytics/feature flag evaluation via PostHog) is isolated into its own service or interacts with a dedicated client. This separation makes testing easier and allows different parts of the system to evolve independently.

### Database Strategy

PostgreSQL on Railway serves dual purposes: storing application data (users, clips, metadata) and managing asynchronous job states. This unified approach simplifies the architecture while Railway handles backup, scaling, and maintenance.

### Async Processing Model

Video rendering is deliberately separated from the main application flow. When a user requests a render, we create a job record and hand it off to Remotion Lambda, then track completion via webhooks. This keeps the main API responsive even during heavy processing.

## 🎯 Why This Approach

**Stateless Design**: No server-side sessions means easy horizontal scaling when needed.

**Clear Separation**: External services (Twitch, Clerk, Remotion, PostHog) are abstracted behind service interfaces or dedicated clients, making them easy to mock for testing or swap later.

**Simple Storage**: Starting with temporary files eliminates AWS complexity for MVP while keeping the interface compatible with future S3 migration.

**Railway Integration**: One-click deployments and managed database remove infrastructure headaches, letting us focus on product development.
