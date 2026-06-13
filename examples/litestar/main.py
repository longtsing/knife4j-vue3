"""
Knife4j Vue3 + LiteStar 示例项目

启动方式:
    pip install -r requirements.txt
    python main.py

访问 Knife4j 文档:
    http://localhost:8000/api/doc.html
"""

from litestar import Litestar, get, Router
from litestar.openapi import OpenAPIConfig
from litestar.openapi.spec import Server
from litestar.config.cors import CORSConfig
from litestar.static_files import create_static_files_router
import os

# ============================================================
# 配置常量
# ============================================================
ROOT_PATH = '/api'
# 静态资源使用 knife4j-vue3 的编译产物
STATIC_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../dist"))

# ============================================================
# 导入路由
# ============================================================
from routers import (
    list_users,
    get_user,
    create_user,
    update_user,
    delete_user,
    health_check,
)

# ============================================================
# 创建 Knife4j 相关路由
# ============================================================

# Swagger 配置端点
@get(f"/{ROOT_PATH}/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    """Knife4j 所需的 Swagger 配置端点"""
    return {
        "urls": [{"url": "/api/schema/openapi.json", "name": "default"}],
        "configUrl": "/api/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

# ============================================================
# 创建路由
# ============================================================

# 创建 API 路由
# handlers 中的路径不含 /api 前缀（如 /users），app 的 path="/api" 会自动添加前缀。
# 实际访问路径示例：/api/users
api_router = Router(
    path="/",
    route_handlers=[
        list_users,
        get_user,
        create_user,
        update_user,
        delete_user,
        health_check,
    ],
)

# 使用 create_static_files_router 统一管理所有静态文件
# 这样 Knife4j 的所有静态资源（doc.html、webjars、oauth、img 等）都能正确服务
static_router = create_static_files_router(
    path="/",
    directories=[STATIC_DIR],
    html_mode=False,
    include_in_schema=False,
)

app = Litestar(
    path=ROOT_PATH,
    route_handlers=[
        swagger_config,
        api_router,
        static_router,  # 静态文件路由放在最后，作为兜底
    ],
    openapi_config=OpenAPIConfig(
        title="Knife4j Vue3 LiteStar 示例",
        version="1.0.0",
        description="这是一个 Knife4j Vue3 + LiteStar 的示例项目",
        path="/schema",
        render_plugins=[],  # 禁用 LiteStar 自带的 Scalar/Swagger UI，使用 Knife4j 替代
        servers=[Server(url=ROOT_PATH, description="API 服务")],
    ),
    cors_config=CORSConfig(allow_origins=["*"]),
)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
