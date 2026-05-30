"""
Knife4j Vue3 + LiteStar 示例项目

启动方式:
    pip install -r requirements.txt
    uvicorn main:app --host 0.0.0.0 --port 8000 --reload

访问 Knife4j 文档:
    http://localhost:8000/doc.html
"""

from litestar import Litestar, get, post, put, delete
from litestar.response import Response
from litestar.openapi import OpenAPIConfig
from litestar.middleware.cors import CORSConfig
from pydantic import BaseModel, Field
from typing import List, Optional
import os

# ============================================================
# 数据模型
# ============================================================

class UserCreate(BaseModel):
    """创建用户请求体"""
    name: str = Field(description="用户名", examples=["张三"])
    email: str = Field(description="邮箱地址", examples=["zhangsan@example.com"])
    role: str = Field(description="用户角色", default="user", examples=["admin"])

class User(BaseModel):
    """用户响应模型"""
    id: int = Field(description="用户ID", examples=[1])
    name: str = Field(description="用户名", examples=["张三"])
    email: str = Field(description="邮箱地址", examples=["zhangsan@example.com"])
    role: str = Field(description="用户角色", examples=["admin"])

# ============================================================
# 模拟数据库
# ============================================================

users_db: dict[int, User] = {}
next_id = 1

# 初始化示例数据
def init_data():
    global next_id
    for name, email, role in [("张三", "zhangsan@example.com", "admin"), ("李四", "lisi@example.com", "user")]:
        user = User(id=next_id, name=name, email=email, role=role)
        users_db[next_id] = user
        next_id += 1

init_data()

# ============================================================
# Knife4j 需要的 swagger-config 端点
# ============================================================

@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    """Knife4j 需要的 swagger-config 端点"""
    return {
        "urls": [
            {
                "url": "/schema/openapi.json",
                "name": "default"
            }
        ],
        "configUrl": "/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

# ============================================================
# Knife4j 前端静态文件托管
# ============================================================

STATIC_DIR = os.path.join(os.path.dirname(__file__), "static")

@get("/doc.html", include_in_schema=False)
async def knife4j_ui() -> Response:
    """Knife4j 入口页面"""
    static_path = os.path.join(STATIC_DIR, "doc.html")
    if os.path.isfile(static_path):
        with open(static_path, "r", encoding="utf-8") as f:
            content = f.read()
        return Response(content=content, media_type="text/html")

    # 如果没有静态文件，返回安装指引
    guide = """
    <html>
    <head><title>Knife4j Vue3 + LiteStar</title></head>
    <body style="font-family: sans-serif; padding: 40px;">
        <h1>Knife4j Vue3 + LiteStar 示例</h1>
        <p>请先将 Knife4j 前端编译产物复制到 static 目录：</p>
        <pre>
# 1. 编译前端
cd knife4j-vue3
pnpm install && pnpm build

# 2. 复制产物
cp -r dist/* examples/litestar/static/

# 3. 重启服务
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
        </pre>
        <p>然后访问 <a href="/doc.html">/doc.html</a></p>
    </body>
    </html>
    """
    return Response(content=guide, media_type="text/html")

@get("/{path:path}", include_in_schema=False)
async def serve_static(path: str) -> Response:
    """静态资源回退路由"""
    static_path = os.path.join(STATIC_DIR, path)
    if os.path.isfile(static_path):
        media_type = "text/plain"
        if path.endswith(".js"):
            media_type = "application/javascript"
        elif path.endswith(".css"):
            media_type = "text/css"
        elif path.endswith(".html"):
            media_type = "text/html"
        elif path.endswith(".json"):
            media_type = "application/json"
        elif path.endswith(".svg"):
            media_type = "image/svg+xml"
        elif path.endswith(".png"):
            media_type = "image/png"
        elif path.endswith(".ico"):
            media_type = "image/x-icon"

        with open(static_path, "rb") as f:
            content = f.read()
        return Response(content=content, media_type=media_type)

    # 回退到 doc.html
    return await knife4j_ui()

# ============================================================
# API 路由
# ============================================================

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

@delete("/api/users/{user_id:int}", tags=["用户管理"], summary="删除用户")
async def delete_user(user_id: int) -> dict:
    """根据ID删除用户"""
    if user_id not in users_db:
        return Response(content={"error": "User not found"}, status_code=404)
    del users_db[user_id]
    return {"message": "deleted"}

@get("/api/health", tags=["系统"], summary="健康检查")
async def health_check() -> dict:
    """健康检查接口"""
    return {"status": "ok", "service": "knife4j-vue3-litestar-example"}

# ============================================================
# 创建 LiteStar 应用
# ============================================================

app = Litestar(
    route_handlers=[
        # Knife4j 静态文件路由（必须在 API 路由之前）
        knife4j_ui,
        serve_static,
        swagger_config,
        # API 路由
        list_users,
        get_user,
        create_user,
        update_user,
        delete_user,
        health_check,
    ],
    openapi_config=OpenAPIConfig(
        title="Knife4j Vue3 LiteStar 示例",
        version="1.0.0",
        description="这是一个 Knife4j Vue3 + LiteStar 的示例项目，展示如何集成 API 文档界面。",
    ),
    cors_config=CORSConfig(
        allow_origins=["*"],
        allow_methods=["*"],
        allow_headers=["*"],
        allow_credentials=True,
    ),
)
