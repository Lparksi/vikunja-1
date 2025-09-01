#!/usr/bin/env python3
"""Startup script for Vikunja Python backend."""

import asyncio
import uvicorn

from vikunja.main import app
from vikunja.db import init_db


async def startup():
    """Initialize the application."""
    print("Initializing database...")
    await init_db()
    print("Database initialized!")


if __name__ == "__main__":
    # Initialize database
    asyncio.run(startup())
    
    # Start the server
    uvicorn.run(
        "vikunja.main:app",
        host="0.0.0.0",
        port=3456,
        reload=True,
    )