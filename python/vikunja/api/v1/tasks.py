"""Task routes."""

from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, or_, and_
from sqlalchemy.orm import selectinload

from vikunja.db import get_db
from vikunja.models.user import User
from vikunja.models.task import Task, TaskCreate, TaskUpdate, TaskResponse
from vikunja.models.project import Project, ProjectUser
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.get("/all", response_model=List[TaskResponse])
async def list_all_tasks(
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
    done: Optional[bool] = Query(None),
    page: int = Query(1, ge=1),
    per_page: int = Query(50, ge=1, le=100)
) -> List[TaskResponse]:
    """List all tasks accessible to the current user."""
    # Get projects user has access to
    accessible_project_ids = await _get_accessible_project_ids(db, current_user)
    
    # Build query
    query = select(Task).where(Task.project_id.in_(accessible_project_ids))
    
    # Apply filters
    if done is not None:
        query = query.where(Task.done == done)
    
    # Apply pagination
    offset = (page - 1) * per_page
    query = query.offset(offset).limit(per_page)
    
    # Order by priority and creation date
    query = query.order_by(Task.priority.desc(), Task.created_at.desc())
    
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return [TaskResponse.model_validate(task) for task in tasks]


@router.get("/{task_id}", response_model=TaskResponse)
async def get_task(
    task_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TaskResponse:
    """Get a specific task."""
    task = await _get_task_with_access(db, task_id, current_user)
    return TaskResponse.model_validate(task)


@router.post("/{task_id}", response_model=TaskResponse)
async def update_task(
    task_id: int,
    task_update: TaskUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> TaskResponse:
    """Update a task."""
    task = await _get_task_with_write_access(db, task_id, current_user)
    
    # Update task fields
    update_data = task_update.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(task, field, value)
    
    await db.commit()
    await db.refresh(task)
    
    return TaskResponse.model_validate(task)


@router.delete("/{task_id}")
async def delete_task(
    task_id: int,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Delete a task."""
    task = await _get_task_with_write_access(db, task_id, current_user)
    
    await db.delete(task)
    await db.commit()
    
    return {"message": "Task deleted successfully"}


# Project-specific task routes (matching Go API structure)
@router.put("/projects/{project_id}/tasks", response_model=TaskResponse)
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


@router.get("/projects/{project_id}/tasks", response_model=List[TaskResponse])
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


async def _get_accessible_project_ids(db: AsyncSession, user: User) -> List[int]:
    """Get list of project IDs the user has access to."""
    result = await db.execute(
        select(Project.id).where(
            or_(
                Project.owner_id == user.id,
                Project.id.in_(
                    select(ProjectUser.project_id).where(
                        ProjectUser.user_id == user.id
                    )
                )
            )
        )
    )
    return [row[0] for row in result.fetchall()]


async def _get_task_with_access(db: AsyncSession, task_id: int, user: User) -> Task:
    """Get task if user has at least read access."""
    accessible_project_ids = await _get_accessible_project_ids(db, user)
    
    result = await db.execute(
        select(Task).where(
            and_(
                Task.id == task_id,
                Task.project_id.in_(accessible_project_ids)
            )
        )
    )
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Task not found"
        )
    
    return task


async def _get_task_with_write_access(db: AsyncSession, task_id: int, user: User) -> Task:
    """Get task if user has write access."""
    result = await db.execute(
        select(Task).where(Task.id == task_id)
    )
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Task not found"
        )
    
    # Check if user has write access to the project
    await _get_project_with_write_access(db, task.project_id, user)
    
    return task


async def _get_project_with_access(db: AsyncSession, project_id: int, user: User) -> Project:
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


async def _get_project_with_write_access(db: AsyncSession, project_id: int, user: User) -> Project:
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