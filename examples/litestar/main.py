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
from litestar.response import File
import os

# ============================================================
# 配置常量
# ============================================================
ROOT_PATH = '/api'
# 静态资源使用 knife4j-vue3 的编译产物
STATIC_DIR = os.path.join(os.path.dirname(__file__), "static")

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
    echo_get,
    echo_post,
)

# ============================================================
# 创建 Knife4j 相关路由
# ============================================================

# Swagger 配置端点 - 路径必须与前端请求路径一致
# 注意：Litestar(path=ROOT_PATH) 会自动添加 /api 前缀，所以这里只需要写 /v3/api-docs/swagger-config
@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    """Knife4j 所需的 Swagger 配置端点"""
    return {
        "urls": [{"url": f"{ROOT_PATH}/openapi.json", "name": "default"}],
        "configUrl": f"{ROOT_PATH}/v3/api-docs/swagger-config",
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
        echo_get,
        echo_post,
    ],
)

# 使用多个静态文件路由分别管理 Knife4j 的静态资源
# doc.html 需要单独处理
@get("/doc.html", include_in_schema=False)
async def serve_doc_html() -> File:
    """服务 Knife4j 的 doc.html"""
    return File(
        path=os.path.join(STATIC_DIR, "doc.html"),
        content_disposition_type="inline",
        media_type="text/html"
    )

# webjars 目录下的静态资源
webjars_router = create_static_files_router(
    path="/webjars",
    directories=[os.path.join(STATIC_DIR, "webjars")],
    include_in_schema=False,
)

# oauth 目录
oauth_router = create_static_files_router(
    path="/oauth",
    directories=[os.path.join(STATIC_DIR, "oauth")],
    include_in_schema=False,
)

# favicon
@get("/favicon.ico", include_in_schema=False)
async def serve_favicon() -> File:
    """服务 favicon"""
    return File(path=os.path.join(STATIC_DIR, "favicon.ico"))

app = Litestar(
    path=ROOT_PATH,
    route_handlers=[
        swagger_config,
        serve_doc_html,
        serve_favicon,
        api_router,
        webjars_router,
        oauth_router,
    ],
    openapi_config=OpenAPIConfig(
        title="Knife4j Vue3 LiteStar 示例",
        version="1.0.0",
        description="这是一个 Knife4j Vue3 + LiteStar 的示例项目",
        path="/",
        render_plugins=[],  # 禁用 LiteStar 自带的 Scalar/Swagger UI，使用 Knife4j 替代
        servers=[Server(url=ROOT_PATH, description="API 服务")],
    ),
    cors_config=CORSConfig(allow_origins=["*"]),
)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
