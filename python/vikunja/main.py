"""Main FastAPI application."""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.trustedhost import TrustedHostMiddleware

from vikunja.api.v1 import router as api_v1_router
from vikunja.config import settings


def create_app() -> FastAPI:
    """Create and configure the FastAPI application."""
    app = FastAPI(
        title="Vikunja API",
        description="The to-do app to organize your life",
        version="0.1.0",
        docs_url="/api/v1/docs",
        redoc_url="/api/v1/redoc",
        openapi_url="/api/v1/docs.json",
    )

    # Add CORS middleware
    if settings.cors_enable:
        app.add_middleware(
            CORSMiddleware,
            allow_origins=settings.cors_origins,
            allow_credentials=True,
            allow_methods=["*"],
            allow_headers=["*"],
            max_age=settings.cors_max_age,
        )

    # Add trusted host middleware
    if settings.service_public_url:
        app.add_middleware(
            TrustedHostMiddleware,
            allowed_hosts=[settings.service_public_url]
        )

    # Include API routes
    app.include_router(api_v1_router, prefix="/api/v1")

    # Health check endpoint
    @app.get("/health")
    async def health_check() -> dict[str, str]:
        """Health check endpoint."""
        return {"status": "ok"}

    # Root endpoint with instructions for development
    @app.get("/")
    async def root() -> dict[str, str]:
        """Root endpoint with development instructions."""
        return {
            "message": "Vikunja Python Backend",
            "version": "0.1.0",
            "api_docs": "/api/v1/docs",
            "health": "/health",
            "note": "This is the API backend only. For the full application, please access the frontend at http://127.0.0.1:4173 in development mode."
        }

    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "vikunja.main:app",
        host="0.0.0.0",
        port=3456,
        reload=True,
    )