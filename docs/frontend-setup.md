# Frontend Architecture

## 🎯 Design Principles

The frontend prioritizes user experience with a clean, responsive interface that makes video editing accessible to streamers who aren't video professionals. Built as a React SPA, it integrates seamlessly with the Go backend through Go-Blueprint's build system.

## 🧩 Component Structure

### Authentication Flow

Clerk handles the entire authentication experience, from social logins to token management. Users can sign in with Twitch, Google, or Discord without the frontend handling any OAuth complexity. The auth state is globally available throughout the app.

### Dashboard Experience

The main dashboard gives streamers an overview of their Twitch clips with easy import functionality. Visual thumbnails and metadata make it simple to browse and select content for editing.

### In-Browser Editor

React Video Editor provides a timeline-based interface for trimming clips without requiring desktop software. The editor sends parameters (not video files) to the backend, keeping the browser responsive and reducing bandwidth usage.

### Real-Time Status Updates

Render job status is polled regularly, showing progress indicators and download buttons when complete. This keeps users informed about processing without requiring complex WebSocket infrastructure.

## 🎨 Visual Design Strategy

### Tailwind CSS v4

Utility-first styling allows rapid prototyping and consistent design. The Twitch-inspired purple color scheme creates familiarity for streamers while maintaining a professional appearance.

### Responsive Layout

The interface works across desktop and mobile devices, though video editing functionality is optimized for larger screens where timeline manipulation is more practical.

## 🔄 State Management Approach

### Custom Hooks Pattern

Business logic is abstracted into custom hooks (useAuth, useClips, useRenderJob) that manage both local state and API interactions. This keeps components focused on presentation while centralizing data logic.

### API Service Layer

All backend communication flows through a dedicated service layer that handles authentication, error handling, and request formatting. This abstraction makes it easy to modify API interactions without changing components.

## 📊 Analytics and Feature Flags

### PostHog Integration

PostHog is integrated into the frontend for comprehensive product analytics and robust feature flag management. This allows us to:

- Track user interactions and feature adoption to understand user behavior.
- Conduct A/B tests and gradually roll out new features.
- Quickly toggle features on or off without requiring new deployments.

The PostHog React SDK is used to capture events and identify users. Feature flags control access to new functionalities, enabling targeted releases and experimentation.

## 🚀 Integration Benefits

**Go-Blueprint Synergy**: The React build integrates directly with the Go server, eliminating the need for separate deployment processes or CORS issues.

**Clerk Simplification**: Frontend never sees OAuth tokens or manages authentication state directly—Clerk handles everything.

**Performance Focus**: By sending only editing parameters (not video data) to the backend, the frontend remains fast even with large video files.
