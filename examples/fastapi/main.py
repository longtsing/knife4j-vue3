"""
Knife4j Vue3 + FastAPI 示例项目

启动方式:
    pip install -r requirements.txt
    uvicorn main:app --host 0.0.0.0 --port 8000 --reload

访问 Knife4j 文档:
    http://localhost:8000/doc.html

集成步骤:
    1. 编译前端: cd knife4j-vue3 && pnpm build
    2. 复制产物: cp -r dist/* examples/fastapi/static/
    3. 启动服务: uvicorn main:app --host 0.0.0.0 --port 8000 --reload
"""

from fastapi import FastAPI, HTTPException, APIRouter
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from pydantic import BaseModel, Field
from typing import List
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

def init_data():
    global next_id
    for name, email, role in [("张三", "zhangsan@example.com", "admin"), ("李四", "lisi@example.com", "user")]:
        user = User(id=next_id, name=name, email=email, role=role)
        users_db[next_id] = user
        next_id += 1

init_data()

# ============================================================
# 创建 FastAPI 应用
# ============================================================

app = FastAPI(
    title="Knife4j Vue3 FastAPI 示例",
    description="这是一个 Knife4j Vue3 + FastAPI 的示例项目，展示如何集成 API 文档界面。",
    version="1.0.0",
    docs_url="/docs",   # 保留默认 Swagger UI 作为备用
    redoc_url=None,
)

# 配置 CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ============================================================
# Knife4j 需要的 swagger-config 端点
# ============================================================

@app.get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config():
    """Knife4j 需要的 swagger-config 端点"""
    return {
        "urls": [
            {
                "url": "/openapi.json",
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

if os.path.exists(STATIC_DIR) and os.path.isfile(os.path.join(STATIC_DIR, "doc.html")):

    @app.get("/doc.html", include_in_schema=False)
    async def knife4j_ui():
        """Knife4j 入口页面"""
        return FileResponse(os.path.join(STATIC_DIR, "doc.html"), media_type="text/html")

    @app.get("/{path:path}", include_in_schema=False)
    async def serve_static(path: str):
        """静态资源回退路由"""
        file_path = os.path.join(STATIC_DIR, path)
        if os.path.isfile(file_path):
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
            return FileResponse(file_path, media_type=media_type)
        # 回退到 doc.html（支持 SPA 路由）
        return FileResponse(os.path.join(STATIC_DIR, "doc.html"), media_type="text/html")

else:
    @app.get("/doc.html", include_in_schema=False)
    async def knife4j_ui_missing():
        """未找到前端静态文件时的提示"""
        return {
            "message": "请先将 knife4j-vue3 编译产物复制到 examples/fastapi/static/ 目录",
            "steps": [
                "1. cd knife4j-vue3 && pnpm install && pnpm build",
                "2. cp -r dist/* examples/fastapi/static/",
                "3. 重启 FastAPI 服务"
            ]
        }

# ============================================================
# API 路由
# ============================================================

@app.get("/api/users", response_model=List[User], tags=["用户管理"], summary="获取用户列表")
async def list_users():
    """返回系统中所有用户的信息"""
    return list(users_db.values())

@app.get("/api/users/{user_id}", response_model=User, tags=["用户管理"], summary="根据ID获取用户")
async def get_user(user_id: int):
    """根据用户ID返回单个用户信息"""
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    return users_db[user_id]

@app.post("/api/users", response_model=User, tags=["用户管理"], summary="创建用户")
async def create_user(data: UserCreate):
    """创建一个新的用户"""
    global next_id
    user = User(id=next_id, **data.model_dump())
    next_id += 1
    users_db[user.id] = user
    return user

@app.put("/api/users/{user_id}", response_model=User, tags=["用户管理"], summary="更新用户")
async def update_user(user_id: int, data: UserCreate):
    """根据ID更新用户信息"""
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    user = User(id=user_id, **data.model_dump())
    users_db[user_id] = user
    return user

@app.delete("/api/users/{user_id}", tags=["用户管理"], summary="删除用户")
async def delete_user(user_id: int):
    """根据ID删除用户"""
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    del users_db[user_id]
    return {"message": "deleted"}

@app.get("/api/health", tags=["系统"], summary="健康检查")
async def health_check():
    """健康检查接口"""
    return {"status": "ok", "service": "knife4j-vue3-fastapi-example"}
