"""Authentication routes."""

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from vikunja.db import get_db
from vikunja.models.user import User, UserLogin, UserCreate, UserResponse, Token
from vikunja.auth import verify_password, get_password_hash, create_access_token, create_long_access_token
from vikunja.auth.dependencies import get_current_user

router = APIRouter()


@router.post("/login", response_model=Token)
async def login(
    user_credentials: UserLogin,
    db: AsyncSession = Depends(get_db)
) -> Token:
    """Login user and return access token."""
    # Find user by username or email
    result = await db.execute(
        select(User).where(
            (User.username == user_credentials.username) | 
            (User.email == user_credentials.username)
        )
    )
    user = result.scalar_one_or_none()
    
    if not user or not verify_password(user_credentials.password, user.password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    if not user.is_active:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Inactive user",
        )
    
    # Create token
    token_data = {"user_id": user.id, "username": user.username}
    if user_credentials.long_token:
        access_token = create_long_access_token(token_data)
    else:
        access_token = create_access_token(token_data)
    
    return Token(token=access_token)


@router.post("/register", response_model=UserResponse)
async def register(
    user_data: UserCreate,
    db: AsyncSession = Depends(get_db)
) -> UserResponse:
    """Register a new user."""
    # Check if user already exists
    result = await db.execute(
        select(User).where(
            (User.username == user_data.username) | 
            (User.email == user_data.email)
        )
    )
    existing_user = result.scalar_one_or_none()
    
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Username or email already registered"
        )
    
    # Check if this is the first user (admin)
    count_result = await db.execute(select(User))
    user_count = len(count_result.scalars().all())
    is_first_user = user_count == 0
    
    # Create new user
    hashed_password = get_password_hash(user_data.password)
    new_user = User(
        username=user_data.username,
        email=user_data.email,
        name=user_data.name,
        password=hashed_password,
        timezone=user_data.timezone,
        week_start=user_data.week_start,
        language=user_data.language,
        is_admin=is_first_user,  # First user becomes admin
    )
    
    db.add(new_user)
    await db.commit()
    await db.refresh(new_user)
    
    return UserResponse.model_validate(new_user)


@router.get("/token/test")
async def test_token(
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Test if token is valid."""
    return {"message": "Token is valid", "user": current_user.username}


@router.post("/token/test")
async def check_token(
    current_user: User = Depends(get_current_user)
) -> dict[str, str]:
    """Check if token is valid."""
    return {"message": "Token is valid", "user": current_user.username}