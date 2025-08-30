"""Task model and related schemas."""

from datetime import datetime
from typing import Optional, List
from enum import IntEnum
from sqlalchemy import String, Boolean, DateTime, Text, Integer, ForeignKey, Float
from sqlalchemy.orm import Mapped, mapped_column, relationship
from pydantic import BaseModel

from vikunja.models.base import Base, TimestampMixin


class TaskRepeatMode(IntEnum):
    """Task repeat mode enum."""
    DEFAULT = 0
    MONTH = 1
    FROM_CURRENT_DATE = 2


class Task(Base, TimestampMixin):
    """Task model."""
    
    __tablename__ = "tasks"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(500), nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    
    # Status
    done: Mapped[bool] = mapped_column(Boolean, default=False)
    done_at: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    
    # Dates
    due_date: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    start_date: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    end_date: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    
    # Reminders (JSON stored as text)
    reminders: Mapped[Optional[str]] = mapped_column(Text)
    
    # Priority and position
    priority: Mapped[int] = mapped_column(Integer, default=0)
    position: Mapped[float] = mapped_column(Float, default=0.0)
    
    # Repeat settings
    repeat_after: Mapped[int] = mapped_column(Integer, default=0)
    repeat_mode: Mapped[int] = mapped_column(Integer, default=TaskRepeatMode.DEFAULT)
    
    # Colors and styling
    hex_color: Mapped[Optional[str]] = mapped_column(String(6))
    
    # Bucket for kanban
    bucket_id: Mapped[Optional[int]] = mapped_column(Integer, ForeignKey("buckets.id"))
    
    # Project
    project_id: Mapped[int] = mapped_column(Integer, ForeignKey("projects.id"), nullable=False)
    
    # Creator
    created_by_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Task hierarchy
    parent_task_id: Mapped[Optional[int]] = mapped_column(Integer, ForeignKey("tasks.id"))
    
    # Index for sorting
    kanban_position: Mapped[float] = mapped_column(Float, default=0.0)
    
    # Relationships
    project = relationship("Project")
    created_by = relationship("User")
    parent_task = relationship("Task", remote_side=[id])
    bucket = relationship("Bucket", back_populates="tasks")


class TaskAssginee(Base, TimestampMixin):
    """Task assignee model."""
    
    __tablename__ = "task_assignees"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    task_id: Mapped[int] = mapped_column(Integer, ForeignKey("tasks.id"), nullable=False)
    user_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    task = relationship("Task")
    user = relationship("User")


class Label(Base, TimestampMixin):
    """Label model."""
    
    __tablename__ = "labels"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(250), nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    hex_color: Mapped[Optional[str]] = mapped_column(String(6))
    
    # Creator
    created_by_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    created_by = relationship("User")


class LabelTask(Base, TimestampMixin):
    """Label-Task relationship model."""
    
    __tablename__ = "label_tasks"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    task_id: Mapped[int] = mapped_column(Integer, ForeignKey("tasks.id"), nullable=False)
    label_id: Mapped[int] = mapped_column(Integer, ForeignKey("labels.id"), nullable=False)
    
    # Relationships
    task = relationship("Task")
    label = relationship("Label")


class TaskComment(Base, TimestampMixin):
    """Task comment model."""
    
    __tablename__ = "task_comments"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    comment: Mapped[str] = mapped_column(Text, nullable=False)
    task_id: Mapped[int] = mapped_column(Integer, ForeignKey("tasks.id"), nullable=False)
    author_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    task = relationship("Task")
    author = relationship("User")


class TaskAttachment(Base, TimestampMixin):
    """Task attachment model."""
    
    __tablename__ = "task_attachments"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    task_id: Mapped[int] = mapped_column(Integer, ForeignKey("tasks.id"), nullable=False)
    file_id: Mapped[int] = mapped_column(Integer, nullable=False)  # Reference to file storage
    created_by_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    task = relationship("Task")
    created_by = relationship("User")


class Bucket(Base, TimestampMixin):
    """Bucket model for kanban boards."""
    
    __tablename__ = "buckets"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(250), nullable=False)
    position: Mapped[float] = mapped_column(Float, default=0.0)
    limit: Mapped[int] = mapped_column(Integer, default=0)  # 0 = no limit
    
    # Project and view
    project_id: Mapped[int] = mapped_column(Integer, ForeignKey("projects.id"), nullable=False)
    project_view_id: Mapped[int] = mapped_column(Integer, ForeignKey("project_views.id"), nullable=False)
    
    # Creator
    created_by_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    project = relationship("Project")
    created_by = relationship("User")
    tasks = relationship("Task", back_populates="bucket")


# Pydantic schemas
class TaskBase(BaseModel):
    """Base task schema."""
    title: str
    description: Optional[str] = None
    done: bool = False
    due_date: Optional[datetime] = None
    start_date: Optional[datetime] = None
    end_date: Optional[datetime] = None
    priority: int = 0
    repeat_after: int = 0
    repeat_mode: int = TaskRepeatMode.DEFAULT
    hex_color: Optional[str] = None
    bucket_id: Optional[int] = None
    parent_task_id: Optional[int] = None


class TaskCreate(TaskBase):
    """Task creation schema."""
    project_id: int


class TaskUpdate(BaseModel):
    """Task update schema."""
    title: Optional[str] = None
    description: Optional[str] = None
    done: Optional[bool] = None
    due_date: Optional[datetime] = None
    start_date: Optional[datetime] = None
    end_date: Optional[datetime] = None
    priority: Optional[int] = None
    repeat_after: Optional[int] = None
    repeat_mode: Optional[int] = None
    hex_color: Optional[str] = None
    bucket_id: Optional[int] = None
    parent_task_id: Optional[int] = None
    position: Optional[float] = None


class TaskResponse(TaskBase):
    """Task response schema."""
    id: int
    project_id: int
    created_by_id: int
    done_at: Optional[datetime] = None
    position: float
    kanban_position: float
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class LabelBase(BaseModel):
    """Base label schema."""
    title: str
    description: Optional[str] = None
    hex_color: Optional[str] = None


class LabelCreate(LabelBase):
    """Label creation schema."""
    pass


class LabelUpdate(BaseModel):
    """Label update schema."""
    title: Optional[str] = None
    description: Optional[str] = None
    hex_color: Optional[str] = None


class LabelResponse(LabelBase):
    """Label response schema."""
    id: int
    created_by_id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class TaskCommentBase(BaseModel):
    """Base task comment schema."""
    comment: str


class TaskCommentCreate(TaskCommentBase):
    """Task comment creation schema."""
    pass


class TaskCommentUpdate(BaseModel):
    """Task comment update schema."""
    comment: Optional[str] = None


class TaskCommentResponse(TaskCommentBase):
    """Task comment response schema."""
    id: int
    task_id: int
    author_id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class BucketBase(BaseModel):
    """Base bucket schema."""
    title: str
    position: float = 0.0
    limit: int = 0


class BucketCreate(BucketBase):
    """Bucket creation schema."""
    pass


class BucketUpdate(BaseModel):
    """Bucket update schema."""
    title: Optional[str] = None
    position: Optional[float] = None
    limit: Optional[int] = None


class BucketResponse(BucketBase):
    """Bucket response schema."""
    id: int
    project_id: int
    project_view_id: int
    created_by_id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True