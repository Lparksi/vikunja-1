"""User model and related schemas."""

from datetime import datetime
from typing import Optional
from sqlalchemy import String, Boolean, DateTime, Text, Integer
from sqlalchemy.orm import Mapped, mapped_column
from pydantic import BaseModel, EmailStr

from vikunja.models.base import Base, TimestampMixin


class User(Base, TimestampMixin):
    """User model."""
    
    __tablename__ = "users"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    username: Mapped[str] = mapped_column(String(250), unique=True, nullable=False)
    password: Mapped[str] = mapped_column(String(250), nullable=False)
    email: Mapped[str] = mapped_column(String(250), unique=True, nullable=False)
    name: Mapped[str] = mapped_column(String(250), nullable=False)
    
    # Settings
    timezone: Mapped[Optional[str]] = mapped_column(String(255))
    week_start: Mapped[int] = mapped_column(Integer, default=0)  # 0 = Sunday, 1 = Monday
    language: Mapped[Optional[str]] = mapped_column(String(5))
    
    # Status fields
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)
    password_reset_token: Mapped[Optional[str]] = mapped_column(String(450))
    email_confirm_token: Mapped[Optional[str]] = mapped_column(String(450))
    is_email_confirmed: Mapped[bool] = mapped_column(Boolean, default=False)
    
    # Deletion
    deletion_scheduled_at: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    deletion_last_reminder_sent: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    
    # TOTP
    totp_secret: Mapped[Optional[str]] = mapped_column(String(32))
    totp_enabled: Mapped[bool] = mapped_column(Boolean, default=False)
    
    # Export
    export_request_date: Mapped[Optional[datetime]] = mapped_column(DateTime(timezone=True))
    
    # Avatar
    avatar_provider: Mapped[str] = mapped_column(String(255), default="initials")
    avatar_file_id: Mapped[Optional[int]] = mapped_column(Integer)
    
    # Settings as JSON text (for complex settings)
    settings: Mapped[Optional[str]] = mapped_column(Text)


# Pydantic schemas
class UserBase(BaseModel):
    """Base user schema."""
    username: str
    email: EmailStr
    name: Optional[str] = None  # Make name optional for compatibility with frontend
    timezone: Optional[str] = None
    week_start: int = 0
    language: Optional[str] = None


class UserCreate(UserBase):
    """User creation schema."""
    password: str


class UserUpdate(BaseModel):
    """User update schema."""
    name: Optional[str] = None
    email: Optional[EmailStr] = None
    timezone: Optional[str] = None
    week_start: Optional[int] = None
    language: Optional[str] = None


class UserResponse(UserBase):
    """User response schema."""
    id: int
    is_active: bool
    is_admin: bool
    is_email_confirmed: bool
    totp_enabled: bool
    avatar_provider: str
    avatar_file_id: Optional[int] = None
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class UserLogin(BaseModel):
    """User login schema."""
    username: str
    password: str
    long_token: bool = False


class Token(BaseModel):
    """Token response schema."""
    token: str
    type: str = "Bearer"