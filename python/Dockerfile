FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy and install Python dependencies
COPY pyproject.toml ./
COPY vikunja/ ./vikunja/

# Install the package
RUN pip install -e .

# Create non-root user
RUN useradd -m -u 1000 vikunja && chown -R vikunja:vikunja /app
USER vikunja

# Expose port
EXPOSE 3456

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3456/health || exit 1

# Run the application
CMD ["python", "-m", "uvicorn", "vikunja.main:app", "--host", "0.0.0.0", "--port", "3456"]