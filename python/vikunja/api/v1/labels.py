"""Label routes."""

from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from vikunja.db import get_db
from vikunja.models.user import User
from vikunja.models.task import Label, LabelCreate, LabelUpdate, LabelResponse
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.get("", response_model=List[LabelResponse])
async def list_labels(
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> List[LabelResponse]:
    """List all labels created by the current user."""
    result = await db.execute(
        select(Label).where(Label.created_by_id == current_user.id)
    )
    labels = result.scalars().all()
    
    return [LabelResponse.model_validate(label) for label in labels]


@router.get("/{label_id}", response_model=LabelResponse)
async def get_label(
    label_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> LabelResponse:
    """Get a specific label."""
    label = await _get_label_with_access(db, label_id, current_user)
    return LabelResponse.model_validate(label)


@router.put("", response_model=LabelResponse)
async def create_label(
    label_data: LabelCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> LabelResponse:
    """Create a new label."""
    new_label = Label(
        **label_data.dict(),
        created_by_id=current_user.id
    )
    
    db.add(new_label)
    await db.commit()
    await db.refresh(new_label)
    
    return LabelResponse.model_validate(new_label)


@router.post("/{label_id}", response_model=LabelResponse)
async def update_label(
    label_id: int,
    label_update: LabelUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> LabelResponse:
    """Update a label."""
    label = await _get_label_with_access(db, label_id, current_user)
    
    # Update label fields
    update_data = label_update.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(label, field, value)
    
    await db.commit()
    await db.refresh(label)
    
    return LabelResponse.model_validate(label)


@router.delete("/{label_id}")
async def delete_label(
    label_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Delete a label."""
    label = await _get_label_with_access(db, label_id, current_user)
    
    await db.delete(label)
    await db.commit()
    
    return {"message": "Label deleted successfully"}


async def _get_label_with_access(db: AsyncSession, label_id: int, user: User) -> Label:
    """Get label if user has access (created by user)."""
    result = await db.execute(
        select(Label).where(
            Label.id == label_id,
            Label.created_by_id == user.id
        )
    )
    label = result.scalar_one_or_none()
    
    if not label:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Label not found"
        )
    
    return label