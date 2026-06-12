"""
Knife4j Vue3 + LiteStar 示例项目

启动方式:
    pip install -r requirements.txt
    uvicorn main:app --host 0.0.0.0 --port 8000 --reload

访问 Knife4j 文档:
    http://localhost:8000/api/doc.html
"""

from litestar import Litestar, get, Router
from litestar.response import Response
from litestar.static_files import StaticFiles
from litestar.openapi import OpenAPIConfig
from litestar.openapi.spec import Server
from litestar.config.cors import CORSConfig
from litestar.static_files.base import FileSystemAdapter
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

# Knife4j 入口页面
@get("/doc.html", include_in_schema=False)
async def doc_html() -> Response:
    """Knife4j 入口页面"""
    doc_path = os.path.join(STATIC_DIR, "doc.html")
    with open(doc_path, "r", encoding="utf-8") as f:
        content = f.read()
    return Response(content=content, media_type="text/html")

# Swagger 配置端点
@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    """Knife4j 所需的 Swagger 配置端点"""
    return {
        "urls": [{"url": "/openapi.json", "name": "default"}],
        "configUrl": "/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

# OpenAPI Schema
@get("/openapi.json", include_in_schema=False)
async def openapi_schema() -> Response:
    """OpenAPI Schema"""
    schema_path = os.path.join(STATIC_DIR, "openapi.json")
    if os.path.exists(schema_path):
        with open(schema_path, "r", encoding="utf-8") as f:
            content = f.read()
        return Response(content=content, media_type="application/json")
    return Response(content="{}", media_type="application/json")

# 创建 API 路由
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

# 静态文件路由
webjars_router = Router(
    path="/webjars",
    route_handlers=[
        StaticFiles(
            directories=[os.path.join(STATIC_DIR, "webjars")],
            file_system=FileSystemAdapter(),
            is_html_mode=False,
        )
    ],
)

oauth_router = Router(
    path="/oauth",
    route_handlers=[
        StaticFiles(
            directories=[os.path.join(STATIC_DIR, "oauth")],
            file_system=FileSystemAdapter(),
            is_html_mode=False,
        )
    ],
)

app = Litestar(
    path=ROOT_PATH,
    route_handlers=[
        doc_html,
        swagger_config,
        openapi_schema,
        webjars_router,
        oauth_router,
        api_router,
    ],
    openapi_config=OpenAPIConfig(
        title="Knife4j Vue3 LiteStar 示例",
        version="1.0.0",
        description="这是一个 Knife4j Vue3 + LiteStar 的示例项目",
        path="/schema",
        servers=[Server(url=ROOT_PATH, description="API 服务")],
    ),
    cors_config=CORSConfig(allow_origins=["*"]),
)
