# ZulgoProxy Admin UI

A modern React-based admin interface for managing the ZulgoProxy server, built with Vite for optimal performance.

## Features

- **Dashboard**: Overview with statistics and charts
- **User Management**: Create, edit, and manage user accounts (Admin only)
- **Proxy Logs**: View and filter proxy request logs (Admin only)
- **System Health**: Monitor system health and manage resources (Admin only)
- **Settings**: User profile and password management

## Development Setup

```bash
# Install dependencies
npm install

# Start development server (with HMR)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

## Technology Stack

- **React 18** - Modern React with hooks
- **Vite** - Fast build tool with HMR
- **Tailwind CSS** - Utility-first CSS framework
- **React Router v6** - Client-side routing
- **Recharts** - Responsive chart library
- **Lucide React** - Modern icon library
- **Axios** - HTTP client with interceptors

## Default Credentials

- Username: `admin`
- Password: `admin`

## API Integration

The UI communicates with the ZulgoProxy API server running on port 8182. Vite's dev server includes proxy configuration to forward API requests seamlessly.

## Project Structure

```
src/
├── api/              # API integration modules
├── components/       # React components
├── context/          # React context providers
├── main.jsx          # Application entry point
├── App.jsx           # Main app component
└── index.css         # Global styles (Tailwind)
```

## Components

- **AuthContext**: JWT authentication state management
- **ProtectedRoute**: Route protection with role-based access
- **Layout**: Main application layout with responsive navigation
- **Dashboard**: Real-time statistics and analytics dashboard
- **Users**: Complete user management interface
- **Logs**: Advanced proxy logs viewer with filtering
- **Health**: System health monitoring and maintenance
- **Settings**: User settings and password management

## Build Process

Vite builds the React app to the `build/` directory with optimized chunks and tree-shaking. The build output is embedded in the Go binary using `go:embed` for seamless single-binary deployment.

## Performance Features

- **Hot Module Replacement (HMR)** - Instant updates during development
- **Code Splitting** - Optimized bundle chunks for faster loading
- **Tree Shaking** - Dead code elimination
- **Modern ES Modules** - Native browser module support
- **Optimized Dependencies** - Pre-bundled vendor libraries