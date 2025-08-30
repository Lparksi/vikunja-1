"""Team routes."""

from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, or_

from vikunja.db import get_db
from vikunja.models.user import User
from vikunja.models.project import Team, TeamCreate, TeamUpdate, TeamResponse, TeamMember
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.get("", response_model=List[TeamResponse])
async def list_teams(
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> List[TeamResponse]:
    """List all teams accessible to the current user."""
    # Get teams where user is creator or member
    result = await db.execute(
        select(Team).where(
            or_(
                Team.created_by_id == current_user.id,
                Team.id.in_(
                    select(TeamMember.team_id).where(
                        TeamMember.user_id == current_user.id
                    )
                )
            )
        )
    )
    teams = result.scalars().all()
    
    return [TeamResponse.model_validate(team) for team in teams]


@router.get("/{team_id}", response_model=TeamResponse)
async def get_team(
    team_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TeamResponse:
    """Get a specific team."""
    team = await _get_team_with_access(db, team_id, current_user)
    return TeamResponse.model_validate(team)


@router.put("", response_model=TeamResponse)
async def create_team(
    team_data: TeamCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TeamResponse:
    """Create a new team."""
    new_team = Team(
        **team_data.dict(),
        created_by_id=current_user.id
    )
    
    db.add(new_team)
    await db.commit()
    await db.refresh(new_team)
    
    return TeamResponse.model_validate(new_team)


@router.post("/{team_id}", response_model=TeamResponse)
async def update_team(
    team_id: int,
    team_update: TeamUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TeamResponse:
    """Update a team."""
    team = await _get_team_with_admin_access(db, team_id, current_user)
    
    # Update team fields
    update_data = team_update.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(team, field, value)
    
    await db.commit()
    await db.refresh(team)
    
    return TeamResponse.model_validate(team)


@router.delete("/{team_id}")
async def delete_team(
    team_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Delete a team."""
    team = await _get_team_with_admin_access(db, team_id, current_user)
    
    await db.delete(team)
    await db.commit()
    
    return {"message": "Team deleted successfully"}


async def _get_team_with_access(db: AsyncSession, team_id: int, user: User) -> Team:
    """Get team if user has access (creator or member)."""
    result = await db.execute(
        select(Team).where(
            Team.id == team_id,
            or_(
                Team.created_by_id == user.id,
                Team.id.in_(
                    select(TeamMember.team_id).where(
                        TeamMember.user_id == user.id
                    )
                )
            )
        )
    )
    team = result.scalar_one_or_none()
    
    if not team:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Team not found"
        )
    
    return team


async def _get_team_with_admin_access(db: AsyncSession, team_id: int, user: User) -> Team:
    """Get team if user has admin access (creator or admin member)."""
    result = await db.execute(
        select(Team).where(
            Team.id == team_id,
            or_(
                Team.created_by_id == user.id,
                Team.id.in_(
                    select(TeamMember.team_id).where(
                        TeamMember.user_id == user.id,
                        TeamMember.admin == True
                    )
                )
            )
        )
    )
    team = result.scalar_one_or_none()
    
    if not team:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Insufficient permissions"
        )
    
    return team