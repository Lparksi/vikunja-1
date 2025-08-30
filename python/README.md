# Vikunja Python Backend

This is the Python implementation of the Vikunja API backend using FastAPI and UV.

## Quick Start

### Running the Complete Application (Backend + Frontend)

1. **Start the Python backend:**
```bash
cd python/
pip install -e .
python run_server.py
```
The backend will be available at http://127.0.0.1:3456

2. **Start the frontend (in a separate terminal):**
```bash
cd frontend/
echo "DEV_PROXY=http://127.0.0.1:3456" > .env.local
pnpm install
pnpm dev
```
The frontend will be available at http://127.0.0.1:4173

3. **Access the application:**
   - **Full Application**: http://127.0.0.1:4173 (Frontend with API proxy)
   - **API Only**: http://127.0.0.1:3456 (Backend only, for testing)
   - **API Documentation**: http://127.0.0.1:3456/api/v1/docs

### Important Notes

- The frontend at http://127.0.0.1:4173 automatically proxies API calls to the Python backend
- The backend at http://127.0.0.1:3456 serves API endpoints only
- For the complete Vikunja experience, always use the frontend URL: http://127.0.0.1:4173

## Setup

1. Install UV if you haven't already:
```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

2. Install dependencies:
```bash
uv sync
```

3. Run the development server:
```bash
uv run uvicorn vikunja.main:app --reload --port 3456
```

## Project Structure

- `vikunja/` - Main application package
  - `api/` - API routes and handlers
  - `models/` - Database models and schemas
  - `auth/` - Authentication and authorization
  - `config/` - Configuration management
  - `db/` - Database connection and utilities
  - `migrations/` - Database migrations
  - `utils/` - Utility functions

## API Compatibility

This implementation maintains full compatibility with the existing Vikunja frontend by preserving the same API endpoints and response formats.

## Testing

Run tests with:
```bash
uv run pytest
```

## Development

Format code:
```bash
uv run black .
uv run isort .
```

Type checking:
```bash
uv run mypy vikunja
```