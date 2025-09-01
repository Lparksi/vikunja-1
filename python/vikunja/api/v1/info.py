"""Info routes."""

from fastapi import APIRouter
from vikunja.config import settings

router = APIRouter()


@router.get("/info")
async def get_info() -> dict:
    """Get application info."""
    # Build enabled background providers array to match Go backend
    enabled_background_providers = []
    if settings.backgrounds_upload_enabled:
        enabled_background_providers.append("upload")
    if settings.backgrounds_unsplash_enabled:
        enabled_background_providers.append("unsplash")
    
    return {
        "version": "0.1.0",
        "frontend_url": settings.service_public_url,
        "motd": settings.service_motd,
        "link_sharing_enabled": settings.service_enable_link_sharing,
        "max_file_size": "20971520",  # 20MB in bytes  
        "available_migrators": [],  # TODO: Add migrators
        "task_attachments_enabled": settings.service_enable_task_attachments,
        "enabled_background_providers": enabled_background_providers,
        "totp_enabled": settings.service_enable_totp,
        "legal": {
            "imprint_url": "",
            "privacy_policy_url": "",
        },
        "caldav_enabled": settings.service_enable_caldav,
        "auth": {
            "local": {
                "enabled": settings.auth_local_enabled,
                "registration_enabled": settings.auth_local_enabled and settings.service_enable_registration,
            },
            "ldap": {
                "enabled": settings.auth_ldap_enabled,
            },
            "openid_connect": {
                "enabled": settings.auth_openid_enabled, 
                "providers": []
            },
        },
        "email_reminders_enabled": settings.service_enable_email_reminders,
        "user_deletion_enabled": settings.service_enable_user_deletion,
        "task_comments_enabled": settings.service_enable_task_comments,
        "demo_mode_enabled": False,
        "webhooks_enabled": settings.webhooks_enabled,
        "public_teams_enabled": settings.service_enable_public_teams,
    }