"""Project view models."""

from datetime import datetime
from typing import Optional
from sqlalchemy import String, Boolean, DateTime, Text, Integer, ForeignKey, Float
from sqlalchemy.orm import Mapped, mapped_column, relationship
from pydantic import BaseModel

from vikunja.models.base import Base, TimestampMixin


class ProjectView(Base, TimestampMixin):
    """Project view model (for different views like Kanban, List, etc)."""
    
    __tablename__ = "project_views"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(250), nullable=False)
    view_kind: Mapped[int] = mapped_column(Integer, default=0)  # 0=list, 1=gantt, 2=table, 3=kanban
    position: Mapped[float] = mapped_column(Float, default=0.0)
    
    # Filter and configuration (stored as JSON text)
    filter: Mapped[Optional[str]] = mapped_column(Text)
    bucket_configuration: Mapped[Optional[str]] = mapped_column(Text)
    
    # Project
    project_id: Mapped[int] = mapped_column(Integer, ForeignKey("projects.id"), nullable=False)
    
    # Relationships
    project = relationship("Project")


# Pydantic schemas
class ProjectViewBase(BaseModel):
    """Base project view schema."""
    title: str
    view_kind: int = 0
    position: float = 0.0
    filter: Optional[str] = None
    bucket_configuration: Optional[str] = None


class ProjectViewCreate(ProjectViewBase):
    """Project view creation schema."""
    pass


class ProjectViewUpdate(BaseModel):
    """Project view update schema."""
    title: Optional[str] = None
    view_kind: Optional[int] = None
    position: Optional[float] = None
    filter: Optional[str] = None
    bucket_configuration: Optional[str] = None


class ProjectViewResponse(ProjectViewBase):
    """Project view response schema."""
    id: int
    project_id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True