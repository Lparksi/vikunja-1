# Vikunja Frontend with TDesign Vue Next

This is a new frontend for Vikunja built with Vue 3, TypeScript, and TDesign Vue Next components. It provides a modern, responsive interface that is fully compatible with the Python FastAPI backend.

## Features

- 🎨 **Modern UI**: Built with TDesign Vue Next components
- ⚡ **Vue 3 + TypeScript**: Full type safety and modern Vue features
- 🔐 **Authentication**: Complete login/register flow with JWT tokens
- 📊 **Project Management**: Create, edit, and manage projects
- ✅ **Task Management**: Full CRUD operations for tasks
- 🔄 **Real-time Updates**: Reactive state management with Pinia
- 📱 **Responsive Design**: Works on desktop and mobile devices

## Development Setup

### Prerequisites

- Node.js 16+ with Yarn package manager
- Python backend running on `http://127.0.0.1:3456`

### Installation

```bash
# Install dependencies
yarn install

# Start development server
yarn dev

# The frontend will be available at http://localhost:4173
```

### Environment Configuration

Create a `.env.local` file to override default settings:

```bash
# Backend proxy URL (defaults to http://127.0.0.1:3456)
DEV_PROXY=http://127.0.0.1:3456

# API base URL (defaults to /api/v1)
VITE_API_BASE_URL=/api/v1
```

## Backend Compatibility

This frontend is designed to work seamlessly with the Python FastAPI backend. It supports:

- **Authentication**: `/api/v1/login` and `/api/v1/register` endpoints
- **Projects**: Full CRUD operations via `/api/v1/projects/*`
- **Tasks**: Task management via `/api/v1/projects/{id}/tasks` and `/api/v1/tasks/*`
- **Admin Features**: First registered user automatically becomes admin

## Architecture

### Directory Structure

```
src/
├── components/        # Reusable Vue components
│   └── AppLayout.vue  # Main application layout
├── views/            # Page components
│   ├── Home.vue      # Dashboard page
│   ├── Login.vue     # Authentication
│   ├── Register.vue  # User registration
│   ├── Projects.vue  # Project listing
│   └── ProjectDetail.vue # Task management
├── stores/           # Pinia state management
│   ├── auth.ts       # Authentication state
│   ├── projects.ts   # Project state
│   └── tasks.ts      # Task state
├── services/         # API communication
│   └── api.ts        # API service layer
├── types/            # TypeScript type definitions
│   └── index.ts      # All type interfaces
└── router/           # Vue Router configuration
    └── index.ts      # Route definitions
```

### State Management

Uses Pinia for reactive state management:

- **Auth Store**: User authentication and session management
- **Project Store**: Project CRUD operations and state
- **Task Store**: Task management and project-specific tasks

### API Integration

The API service (`src/services/api.ts`) provides:

- Automatic JWT token handling
- Request/response interceptors
- Error handling and automatic logout on 401
- TypeScript interfaces for all API responses

## Development Commands

```bash
# Start development server
yarn dev

# Build for production
yarn build

# Preview production build
yarn preview

# Type checking
yarn type-check
```

## Features Overview

### Authentication
- User registration with automatic admin assignment for first user
- JWT token-based authentication
- Automatic token refresh and logout on expiry
- "Remember me" functionality for extended sessions

### Dashboard
- Welcome message with user information
- Quick access to recent projects
- Project and task statistics
- Quick navigation to main features

### Project Management
- Create, edit, and delete projects
- Project color coding for visual organization
- Project descriptions and metadata
- Grid-based project overview

### Task Management
- Create tasks within projects
- Mark tasks as complete/incomplete
- Task descriptions and due dates
- Real-time task updates
- Task deletion and editing

## Customization

### Theming
TDesign Vue Next supports theming. You can customize colors and styles by modifying the TDesign configuration in `src/main.ts`.

### Components
All components are built with TDesign Vue Next components for consistency. Custom components follow the same design patterns.

### API Integration
To modify API endpoints or add new features, update the `src/services/api.ts` file and corresponding TypeScript types in `src/types/index.ts`.

## Production Deployment

```bash
# Build for production
yarn build

# The built files will be in the `dist` directory
# Deploy the contents to your web server
```

For proxy configuration in production, ensure your web server properly proxies `/api/*` requests to the Python backend.

## Differences from Original Frontend

This new frontend offers several improvements over the original:

1. **Modern Component Library**: TDesign Vue Next provides better accessibility and mobile support
2. **Full TypeScript**: Enhanced type safety and developer experience
3. **Simplified State Management**: Cleaner Pinia stores vs complex Vuex modules
4. **Better Performance**: Optimized bundle size and load times
5. **Enhanced UX**: More intuitive navigation and responsive design

## Contributing

When adding new features:

1. Add TypeScript types to `src/types/index.ts`
2. Update API service in `src/services/api.ts`
3. Create/update Pinia stores for state management
4. Build UI components with TDesign Vue Next
5. Follow the existing file structure and naming conventions

## License

This project follows the same license as the main Vikunja project (AGPL-3.0-or-later).