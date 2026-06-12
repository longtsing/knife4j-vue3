"""
用户管理路由
"""
from litestar import get, post, put, delete
from litestar.response import Response
from typing import List

from models.user import User, UserCreate

# ============================================================
# 模拟数据库
# ============================================================

users_db: dict[int, User] = {}
next_id = 1


def init_data():
    """初始化示例数据"""
    global next_id
    for name, email, role in [("张三", "zhangsan@example.com", "admin"), ("李四", "lisi@example.com", "user")]:
        user = User(id=next_id, name=name, email=email, role=role)
        users_db[next_id] = user
        next_id += 1


init_data()


@get("/api/users", tags=["用户管理"], summary="获取用户列表")
async def list_users() -> List[User]:
    """返回系统中所有用户的信息"""
    return list(users_db.values())


@get("/api/users/{user_id:int}", tags=["用户管理"], summary="根据ID获取用户")
async def get_user(user_id: int) -> User:
    """根据用户ID返回单个用户信息"""
    if user_id not in users_db:
        return Response(content={"error": "User not found"}, status_code=404)
    return users_db[user_id]


@post("/api/users", tags=["用户管理"], summary="创建用户")
async def create_user(data: UserCreate) -> User:
    """创建一个新的用户"""
    global next_id
    user = User(id=next_id, **data.model_dump())
    next_id += 1
    users_db[user.id] = user
    return user


@put("/api/users/{user_id:int}", tags=["用户管理"], summary="更新用户")
async def update_user(user_id: int, data: UserCreate) -> User:
    """根据ID更新用户信息"""
    if user_id not in users_db:
        return Response(content={"error": "User not found"}, status_code=404)
    user = User(id=user_id, **data.model_dump())
    users_db[user_id] = user
    return user


@delete("/api/users/{user_id:int}", tags=["用户管理"], summary="删除用户", status_code=200)
async def delete_user(user_id: int) -> dict:
    """根据ID删除用户"""
    if user_id not in users_db:
        return Response(content={"error": "User not found"}, status_code=404)
    del users_db[user_id]
    return {"message": "deleted"}
