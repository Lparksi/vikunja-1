"""Project routes."""

from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, or_, and_
from sqlalchemy.orm import selectinload

from vikunja.db import get_db
from vikunja.models.user import User
from vikunja.models.project import (
    Project, ProjectCreate, ProjectUpdate, ProjectResponse,
    ProjectUser, ProjectUserResponse
)
from vikunja.models.task import Task, TaskCreate, TaskResponse
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.get("", response_model=List[ProjectResponse])
async def list_projects(
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
    is_archived: Optional[bool] = Query(None),
    parent_project_id: Optional[int] = Query(None)
) -> List[ProjectResponse]:
    """List all projects accessible to the current user."""
    # Build query to get projects user has access to
    query = select(Project).where(
        or_(
            Project.owner_id == current_user.id,
            Project.id.in_(
                select(ProjectUser.project_id).where(
                    ProjectUser.user_id == current_user.id
                )
            )
        )
    )
    
    # Apply filters
    if is_archived is not None:
        query = query.where(Project.is_archived == is_archived)
    
    if parent_project_id is not None:
        query = query.where(Project.parent_project_id == parent_project_id)
    
    # Order by position and creation date
    query = query.order_by(Project.position, Project.created_at)
    
    result = await db.execute(query)
    projects = result.scalars().all()
    
    return [ProjectResponse.model_validate(project) for project in projects]


@router.get("/{project_id}", response_model=ProjectResponse)
async def get_project(
    project_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> ProjectResponse:
    """Get a specific project."""
    project = await _get_project_with_access(db, project_id, current_user)
    return ProjectResponse.model_validate(project)


@router.put("", response_model=ProjectResponse)
async def create_project(
    project_data: ProjectCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> ProjectResponse:
    """Create a new project."""
    new_project = Project(
        **project_data.dict(),
        owner_id=current_user.id
    )
    
    db.add(new_project)
    await db.commit()
    await db.refresh(new_project)
    
    return ProjectResponse.model_validate(new_project)


@router.post("/{project_id}", response_model=ProjectResponse)
async def update_project(
    project_id: int,
    project_update: ProjectUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> ProjectResponse:
    """Update a project."""
    project = await _get_project_with_write_access(db, project_id, current_user)
    
    # Update project fields
    update_data = project_update.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(project, field, value)
    
    await db.commit()
    await db.refresh(project)
    
    return ProjectResponse.model_validate(project)


@router.delete("/{project_id}")
async def delete_project(
    project_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Delete a project."""
    project = await _get_project_with_admin_access(db, project_id, current_user)
    
    await db.delete(project)
    await db.commit()
    
    return {"message": "Project deleted successfully"}


@router.get("/{project_id}/users", response_model=List[ProjectUserResponse])
async def list_project_users(
    project_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> List[ProjectUserResponse]:
    """List users with access to a project."""
    await _get_project_with_access(db, project_id, current_user)
    
    result = await db.execute(
        select(ProjectUser).where(ProjectUser.project_id == project_id)
    )
    project_users = result.scalars().all()
    
    return [ProjectUserResponse.model_validate(pu) for pu in project_users]


@router.put("/{project_id}/tasks", response_model=TaskResponse)
async def create_project_task(
    project_id: int,
    task_data: TaskCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TaskResponse:
    """Create a new task in a project."""
    # Verify user has write access to the project
    await _get_project_with_write_access(db, project_id, current_user)
    
    # Override project_id to ensure consistency
    task_data.project_id = project_id
    
    new_task = Task(
        **task_data.dict(),
        created_by_id=current_user.id
    )
    
    db.add(new_task)
    await db.commit()
    await db.refresh(new_task)
    
    return TaskResponse.model_validate(new_task)


@router.get("/{project_id}/tasks", response_model=List[TaskResponse])
async def list_project_tasks(
    project_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
    done: Optional[bool] = Query(None),
    page: int = Query(1, ge=1),
    per_page: int = Query(50, ge=1, le=100)
) -> List[TaskResponse]:
    """List tasks in a specific project."""
    # Verify user has access to the project
    await _get_project_with_access(db, project_id, current_user)
    
    # Build query
    query = select(Task).where(Task.project_id == project_id)
    
    # Apply filters
    if done is not None:
        query = query.where(Task.done == done)
    
    # Apply pagination
    offset = (page - 1) * per_page
    query = query.offset(offset).limit(per_page)
    
    # Order by position and creation date
    query = query.order_by(Task.position, Task.created_at)
    
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return [TaskResponse.model_validate(task) for task in tasks]


async def _get_project_with_access(
    db: AsyncSession, 
    project_id: int, 
    user: User
) -> Project:
    """Get project if user has at least read access."""
    result = await db.execute(
        select(Project).where(
            and_(
                Project.id == project_id,
                or_(
                    Project.owner_id == user.id,
                    Project.id.in_(
                        select(ProjectUser.project_id).where(
                            and_(
                                ProjectUser.user_id == user.id,
                                ProjectUser.right >= 0  # Read access
                            )
                        )
                    )
                )
            )
        )
    )
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Project not found"
        )
    
    return project


async def _get_project_with_write_access(
    db: AsyncSession, 
    project_id: int, 
    user: User
) -> Project:
    """Get project if user has write access."""
    result = await db.execute(
        select(Project).where(
            and_(
                Project.id == project_id,
                or_(
                    Project.owner_id == user.id,
                    Project.id.in_(
                        select(ProjectUser.project_id).where(
                            and_(
                                ProjectUser.user_id == user.id,
                                ProjectUser.right >= 1  # Write access
                            )
                        )
                    )
                )
            )
        )
    )
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Insufficient permissions"
        )
    
    return project


async def _get_project_with_admin_access(
    db: AsyncSession, 
    project_id: int, 
    user: User
) -> Project:
    """Get project if user has admin access."""
    result = await db.execute(
        select(Project).where(
            and_(
                Project.id == project_id,
                or_(
                    Project.owner_id == user.id,
                    Project.id.in_(
                        select(ProjectUser.project_id).where(
                            and_(
                                ProjectUser.user_id == user.id,
                                ProjectUser.right >= 2  # Admin access
                            )
                        )
                    )
                )
            )
        )
    )
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Insufficient permissions"
        )
    
    return project