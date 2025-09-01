"""Project model and related schemas."""

from datetime import datetime
from typing import Optional, List
from sqlalchemy import String, Boolean, DateTime, Text, Integer, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship
from pydantic import BaseModel

from vikunja.models.base import Base, TimestampMixin


class Project(Base, TimestampMixin):
    """Project model."""
    
    __tablename__ = "projects"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(250), nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    identifier: Mapped[Optional[str]] = mapped_column(String(10))
    
    # Hierarchy
    parent_project_id: Mapped[Optional[int]] = mapped_column(
        Integer, ForeignKey("projects.id")
    )
    position: Mapped[float] = mapped_column(default=0.0)
    
    # Colors and styling
    hex_color: Mapped[Optional[str]] = mapped_column(String(6))
    background_file_id: Mapped[Optional[int]] = mapped_column(Integer)
    background_blur_hash: Mapped[Optional[str]] = mapped_column(String(100))
    
    # Status
    is_archived: Mapped[bool] = mapped_column(Boolean, default=False)
    is_favorite: Mapped[bool] = mapped_column(Boolean, default=False)
    
    # Owner
    owner_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    parent_project = relationship("Project", remote_side=[id])
    owner = relationship("User", back_populates="owned_projects")


class ProjectUser(Base, TimestampMixin):
    """Project user permissions."""
    
    __tablename__ = "project_users"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    project_id: Mapped[int] = mapped_column(Integer, ForeignKey("projects.id"), nullable=False)
    user_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    right: Mapped[int] = mapped_column(Integer, default=0)  # 0=read, 1=write, 2=admin
    
    # Relationships
    project = relationship("Project")
    user = relationship("User")


class Team(Base, TimestampMixin):
    """Team model."""
    
    __tablename__ = "teams"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(250), nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    
    # Owner
    created_by_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    
    # Relationships
    created_by = relationship("User")


class TeamMember(Base, TimestampMixin):
    """Team member model."""
    
    __tablename__ = "team_members"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    team_id: Mapped[int] = mapped_column(Integer, ForeignKey("teams.id"), nullable=False)
    user_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"), nullable=False)
    admin: Mapped[bool] = mapped_column(Boolean, default=False)
    
    # Relationships
    team = relationship("Team")
    user = relationship("User")


class TeamProject(Base, TimestampMixin):
    """Team project permissions."""
    
    __tablename__ = "team_projects"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    team_id: Mapped[int] = mapped_column(Integer, ForeignKey("teams.id"), nullable=False)
    project_id: Mapped[int] = mapped_column(Integer, ForeignKey("projects.id"), nullable=False)
    right: Mapped[int] = mapped_column(Integer, default=0)  # 0=read, 1=write, 2=admin
    
    # Relationships
    team = relationship("Team")
    project = relationship("Project")


# Pydantic schemas
class ProjectBase(BaseModel):
    """Base project schema."""
    title: str
    description: Optional[str] = None
    identifier: Optional[str] = None
    parent_project_id: Optional[int] = None
    position: float = 0.0
    hex_color: Optional[str] = None
    is_archived: bool = False
    is_favorite: bool = False


class ProjectCreate(ProjectBase):
    """Project creation schema."""
    pass


class ProjectUpdate(BaseModel):
    """Project update schema."""
    title: Optional[str] = None
    description: Optional[str] = None
    identifier: Optional[str] = None
    parent_project_id: Optional[int] = None
    position: Optional[float] = None
    hex_color: Optional[str] = None
    is_archived: Optional[bool] = None
    is_favorite: Optional[bool] = None


class ProjectResponse(ProjectBase):
    """Project response schema."""
    id: int
    owner_id: int
    background_file_id: Optional[int] = None
    background_blur_hash: Optional[str] = None
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class TeamBase(BaseModel):
    """Base team schema."""
    name: str
    description: Optional[str] = None


class TeamCreate(TeamBase):
    """Team creation schema."""
    pass


class TeamUpdate(BaseModel):
    """Team update schema."""
    name: Optional[str] = None
    description: Optional[str] = None


class TeamResponse(TeamBase):
    """Team response schema."""
    id: int
    created_by_id: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class ProjectUserResponse(BaseModel):
    """Project user response schema."""
    id: int
    project_id: int
    user_id: int
    right: int
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True