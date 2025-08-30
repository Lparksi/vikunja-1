"""API v1 routes."""

from fastapi import APIRouter

from vikunja.api.v1 import auth, users, projects, tasks, labels, teams, info

router = APIRouter()

# Include all route modules
router.include_router(auth.router, tags=["auth"])
router.include_router(info.router, tags=["info"])  
router.include_router(users.router, prefix="/user", tags=["users"])
router.include_router(projects.router, prefix="/projects", tags=["projects"])
router.include_router(tasks.router, prefix="/tasks", tags=["tasks"])
router.include_router(labels.router, prefix="/labels", tags=["labels"])
router.include_router(teams.router, prefix="/teams", tags=["teams"])