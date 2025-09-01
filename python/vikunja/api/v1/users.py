"""User routes."""

from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from vikunja.db import get_db
from vikunja.models.user import User, UserResponse, UserUpdate
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.get("", response_model=UserResponse)
async def get_current_user_info(
    current_user: User = Depends(get_current_user)
) -> UserResponse:
    """Get current user information."""
    return UserResponse.model_validate(current_user)


@router.get("s", response_model=List[UserResponse])
async def list_users(
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> List[UserResponse]:
    """List all users."""
    result = await db.execute(select(User).where(User.is_active == True))
    users = result.scalars().all()
    return [UserResponse.model_validate(user) for user in users]


@router.post("", response_model=UserResponse)
async def update_user(
    user_update: UserUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> UserResponse:
    """Update current user information."""
    # Update user fields
    update_data = user_update.dict(exclude_unset=True)
    
    for field, value in update_data.items():
        setattr(current_user, field, value)
    
    await db.commit()
    await db.refresh(current_user)
    
    return UserResponse.model_validate(current_user)


@router.get("/timezones")
async def get_available_timezones(
    current_user: User = Depends(get_current_user)
) -> List[str]:
    """Get available timezones."""
    # Simplified timezone list - in a real implementation, 
    # you'd want to use a proper timezone library
    return [
        "UTC",
        "America/New_York", 
        "America/Chicago",
        "America/Denver",
        "America/Los_Angeles",
        "Europe/London",
        "Europe/Paris",
        "Europe/Berlin",
        "Asia/Tokyo",
        "Asia/Shanghai",
        "Australia/Sydney",
    ]