# Vikunja Python Backend

This is the Python implementation of the Vikunja API backend using FastAPI and UV.

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