"""Info routes."""

from fastapi import APIRouter
from vikunja.config import settings

router = APIRouter()


@router.get("/info")
async def get_info() -> dict:
    """Get application info."""
    return {
        "version": "0.1.0",
        "frontend_url": settings.service_public_url,
        "motd": settings.service_motd,
        "link_sharing_enabled": settings.service_enable_link_sharing,
        "max_file_size": "20971520",  # 20MB in bytes
        "registration_enabled": settings.service_enable_registration,
        "available_migrators": [],  # TODO: Add migrators
        "task_attachments_enabled": settings.service_enable_task_attachments,
        "enabled_backgrounds": {
            "upload": settings.backgrounds_upload_enabled,
            "unsplash": settings.backgrounds_unsplash_enabled,
        },
        "totp_enabled": settings.service_enable_totp,
        "legal": {
            "imprint_url": "",
            "privacy_policy_url": "",
        },
        "caldav_enabled": settings.service_enable_caldav,
        "user_deletion_enabled": settings.service_enable_user_deletion,
        "task_comments_enabled": settings.service_enable_task_comments,
        "email_reminders_enabled": settings.service_enable_email_reminders,
        "user_avatar_provider": "initials",
        "demo_mode_enabled": False,
        "auth": {
            "local": {"enabled": settings.auth_local_enabled},
            "ldap": {"enabled": settings.auth_ldap_enabled},
            "openid_connect": {"enabled": settings.auth_openid_enabled, "providers": []},
        },
        "webhooks_enabled": settings.webhooks_enabled,
    }