"""
Knife4j Vue3 + LiteStar 示例项目

启动方式:
    pip install -r requirements.txt
    uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api

访问 Knife4j 文档:
    http://localhost:8000/api/doc.html
"""

from litestar import Litestar, get
from litestar.response import Response
from litestar.static_files import create_static_files_router
from litestar.openapi import OpenAPIConfig
from litestar.openapi.spec import Server
from litestar.config.cors import CORSConfig
import os

# ============================================================
# 配置常量
# ============================================================
ROOT_PATH = '/api'
# 静态资源使用 knife4j-vue3 的编译产物
STATIC_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../dist"))

# ============================================================
# Knife4j 文档路由
# ============================================================

@get("/doc.html", include_in_schema=False)
async def doc_html() -> Response:
    """Knife4j 入口页面"""
    with open(os.path.join(STATIC_DIR, "doc.html"), "r", encoding="utf-8") as f:
        content = f.read()
    return Response(content=content, media_type="text/html")


@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    """Knife4j 所需的 Swagger 配置端点"""
    return {
        "urls": [{"url": f"{ROOT_PATH}/schema/openapi.json", "name": "default"}],
        "configUrl": f"{ROOT_PATH}/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }


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
# 创建 LiteStar 应用
# ============================================================

# 创建静态文件路由
static_router = create_static_files_router(path="/", directories=[STATIC_DIR], name="static")

app = Litestar(
    route_handlers=[
        # Knife4j 入口页面
        doc_html,
        swagger_config,
        # 静态资源
        static_router,
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
        description="这是一个 Knife4j Vue3 + LiteStar 的示例项目",
        servers=[Server(url=ROOT_PATH, description="API 服务")],
    ),
    cors_config=CORSConfig(allow_origins=["*"]),
)
