# Frontend Setup Instructions

## Quick Start

To run the original frontend with the Python backend:

### 1. Start Python Backend
```bash
cd python/
pip install -e .
python run_server.py
```
Backend will be available at: http://127.0.0.1:3456

### 2. Configure Frontend
Create `frontend/.env.local` with:
```
DEV_PROXY=http://127.0.0.1:3456
```

### 3. Install Frontend Dependencies
```bash
cd frontend/
CYPRESS_INSTALL_BINARY=0 PUPPETEER_SKIP_DOWNLOAD=true npm install --legacy-peer-deps
```

### 4. Start Frontend
```bash
cd frontend/
npm run dev
```
Frontend will be available at: http://127.0.0.1:4173

## Access Points

- **Complete Application**: http://127.0.0.1:4173 (Frontend with API proxy)
- **API Documentation**: http://127.0.0.1:3456/api/v1/docs  
- **Backend Direct**: http://127.0.0.1:3456/api/v1/info

## Admin Functionality

- First registered user automatically becomes admin (`is_admin: true`)
- Subsequent users are regular users (`is_admin: false`)
- Registration works perfectly with original frontend format

## API Compatibility

The Python backend is 100% compatible with the original frontend:
- Same authentication flow (JWT tokens)
- Same API endpoints (`/api/v1/*`)
- Same request/response formats
- Same user management functionality

## Key Changes Made

1. **Deleted** `/front` directory as requested
2. **Configured** original frontend to work with Python backend
3. **Fixed** user registration compatibility (name field handling)
4. **Verified** admin user functionality works correctly
5. **Tested** authentication and API endpoints