"""Models package."""

# Import all models to make them available
from vikunja.models.base import Base, TimestampMixin
from vikunja.models.user import User
from vikunja.models.project import Project, ProjectUser, Team, TeamMember, TeamProject
from vikunja.models.project_view import ProjectView
from vikunja.models.task import Task, TaskAssginee, Label, LabelTask, TaskComment, TaskAttachment, Bucket

# Add missing relationships to User model
from sqlalchemy.orm import relationship

# Add these to the User class
User.owned_projects = relationship("Project", back_populates="owner")

__all__ = [
    "Base", 
    "TimestampMixin",
    "User",
    "Project", 
    "ProjectUser", 
    "Team", 
    "TeamMember", 
    "TeamProject",
    "ProjectView",
    "Task", 
    "TaskAssginee", 
    "Label", 
    "LabelTask", 
    "TaskComment", 
    "TaskAttachment", 
    "Bucket",
]