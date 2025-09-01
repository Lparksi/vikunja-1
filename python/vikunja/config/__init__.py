"""Configuration management."""

from typing import List, Optional
from pydantic_settings import BaseSettings
from pydantic import Field
import os


class Settings(BaseSettings):
    """Application settings."""
    
    # Service settings
    service_interface: str = Field(default=":3456", env="VIKUNJA_SERVICE_INTERFACE")
    service_public_url: Optional[str] = Field(default=None, env="VIKUNJA_SERVICE_PUBLICURL")
    service_root_path: str = Field(default=".", env="VIKUNJA_SERVICE_ROOTPATH")
    service_max_items_per_page: int = Field(default=50, env="VIKUNJA_SERVICE_MAXITEMSPERPAGE")
    service_enable_caldav: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLECALDAV")
    service_motd: str = Field(default="", env="VIKUNJA_SERVICE_MOTD")
    service_enable_link_sharing: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLELINKSHARING")
    service_enable_registration: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLEREGISTRATION")
    service_enable_task_attachments: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLETASKATTACHMENTS")
    service_timezone: str = Field(default="GMT", env="VIKUNJA_SERVICE_TIMEZONE")
    service_enable_task_comments: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLETASKCOMMENTS")
    service_enable_totp: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLETOTP")
    service_enable_email_reminders: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLEEMAILREMINDERS")
    service_enable_user_deletion: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLEUSERDELETION")
    service_enable_public_teams: bool = Field(default=True, env="VIKUNJA_SERVICE_ENABLEPUBLICTEAMS")
    
    # JWT settings
    jwt_secret: str = Field(default="your-secret-key", env="VIKUNJA_SERVICE_JWTSECRET")
    jwt_ttl: int = Field(default=259200, env="VIKUNJA_SERVICE_JWTTTL")  # 3 days
    jwt_ttl_long: int = Field(default=2592000, env="VIKUNJA_SERVICE_JWTTTLLONG")  # 30 days
    
    # Database settings
    database_type: str = Field(default="sqlite", env="VIKUNJA_DATABASE_TYPE")
    database_host: str = Field(default="localhost", env="VIKUNJA_DATABASE_HOST")
    database_port: int = Field(default=5432, env="VIKUNJA_DATABASE_PORT")
    database_user: str = Field(default="vikunja", env="VIKUNJA_DATABASE_USER")
    database_password: str = Field(default="", env="VIKUNJA_DATABASE_PASSWORD")
    database_database: str = Field(default="vikunja", env="VIKUNJA_DATABASE_DATABASE")
    database_path: str = Field(default="./vikunja.db", env="VIKUNJA_DATABASE_PATH")
    
    # CORS settings
    cors_enable: bool = Field(default=True, env="VIKUNJA_CORS_ENABLE")
    cors_origins: List[str] = Field(default=["*"], env="VIKUNJA_CORS_ORIGINS")
    cors_max_age: int = Field(default=86400, env="VIKUNJA_CORS_MAXAGE")
    
    # Auth settings
    auth_local_enabled: bool = Field(default=True, env="VIKUNJA_AUTH_LOCAL_ENABLED")
    auth_ldap_enabled: bool = Field(default=False, env="VIKUNJA_AUTH_LDAP_ENABLED")
    auth_openid_enabled: bool = Field(default=False, env="VIKUNJA_AUTH_OPENID_ENABLED")
    
    # Rate limiting
    rate_limit_enabled: bool = Field(default=False, env="VIKUNJA_RATELIMIT_ENABLED")
    rate_limit_kind: str = Field(default="user", env="VIKUNJA_RATELIMIT_KIND")
    rate_limit_no_auth_routes_limit: int = Field(default=10, env="VIKUNJA_RATELIMIT_NOAUTHROUTESLIMIT")
    
    # File settings
    files_base_path: str = Field(default="./files", env="VIKUNJA_FILES_BASEPATH")
    files_max_size: str = Field(default="20MB", env="VIKUNJA_FILES_MAXSIZE")
    
    # Logging
    log_enabled: bool = Field(default=True, env="VIKUNJA_LOG_ENABLED")
    log_level: str = Field(default="INFO", env="VIKUNJA_LOG_LEVEL")
    log_format: str = Field(default="json", env="VIKUNJA_LOG_FORMAT")
    
    # Webhooks
    webhooks_enabled: bool = Field(default=False, env="VIKUNJA_WEBHOOKS_ENABLED")
    
    # Backgrounds
    backgrounds_enabled: bool = Field(default=True, env="VIKUNJA_BACKGROUNDS_ENABLED")
    backgrounds_upload_enabled: bool = Field(default=True, env="VIKUNJA_BACKGROUNDS_UPLOAD_ENABLED")
    backgrounds_unsplash_enabled: bool = Field(default=False, env="VIKUNJA_BACKGROUNDS_UNSPLASH_ENABLED")
    
    # Migration settings
    migration_todoist_enable: bool = Field(default=True, env="VIKUNJA_MIGRATION_TODOIST_ENABLE")
    migration_trello_enable: bool = Field(default=True, env="VIKUNJA_MIGRATION_TRELLO_ENABLE")
    migration_microsoft_todo_enable: bool = Field(default=True, env="VIKUNJA_MIGRATION_MICROSOFTTODO_ENABLE")
    
    @property
    def database_url(self) -> str:
        """Get the database URL."""
        if self.database_type == "sqlite":
            return f"sqlite+aiosqlite:///{self.database_path}"
        elif self.database_type == "postgresql":
            return f"postgresql+asyncpg://{self.database_user}:{self.database_password}@{self.database_host}:{self.database_port}/{self.database_database}"
        elif self.database_type == "mysql":
            return f"mysql+aiomysql://{self.database_user}:{self.database_password}@{self.database_host}:{self.database_port}/{self.database_database}"
        else:
            raise ValueError(f"Unsupported database type: {self.database_type}")

    class Config:
        """Pydantic config."""
        env_file = ".env"
        env_file_encoding = "utf-8"


# Global settings instance
settings = Settings()