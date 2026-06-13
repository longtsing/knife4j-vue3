"""
FastAPI 后端示例项目 - 后端入口

基于 FastAPI 构建，使用 Knife4j Vue3 作为 API 文档 UI。
"""
import os
import sys
import base64
from datetime import datetime
from typing import Optional

from fastapi import FastAPI, Request, Response, Query, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, ORJSONResponse
from starlette.staticfiles import StaticFiles

# ============================================================
# 配置常量
# ============================================================
ROOT_PATH = '/api'  # 反向代理前缀，部署时按需修改
PORT = 8000
START_TIME = datetime.now()
STATIC_DIR = os.path.join(os.path.dirname(__file__), "../../dist")

# ============================================================
# FastAPI 应用实例
# ============================================================
app = FastAPI(
    root_path=ROOT_PATH,
    title='FastAPI 后端示例项目',
    description='FastAPI 后端示例项目 API 文档',
    version='1.0.0',
    docs_url=None,       # 禁用内置 Swagger UI，使用 Knife4j 替代
    redoc_url=None,      # 禁用内置 ReDoc
    servers=[{"url": ROOT_PATH, "description": "API 服务"}],  # Swagger UI 所有端点自动加此前缀
    contact={'name': 'Admin', 'email': 'admin@example.com'},
    license_info={
        'name': 'Apache 2.0',
        'url': 'http://www.apache.org/licenses/LICENSE-2.0.html'
    },
    default_response_class=ORJSONResponse,
)

# ============================================================
# 中间件
# ============================================================

# CORS 跨域
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 文档访问认证（Basic Auth）
@app.middleware("http")
async def auth_middleware(request: Request, call_next):
    protected_paths = [
        f"{ROOT_PATH}/docs",
        f"{ROOT_PATH}/openapi.json",
        f"{ROOT_PATH}/redoc",
        f"{ROOT_PATH}/doc.html",
        f"{ROOT_PATH}/v3/api-docs/swagger-config",
    ]

    if request.url.path not in protected_paths:
        return await call_next(request)

    # 解析 Basic Auth
    username, password = "", ""
    try:
        encoded = request.headers["Authorization"]
        decoded = base64.b64decode(encoded[6:]).decode("utf-8")
        username, password = decoded.split(":")
    except Exception:
        pass

    # 验证账号密码
    credentials = {
        "admin": "admin12345",
    }
    if credentials.get(username) == password:
        return await call_next(request)

    return Response(
        content="Authorization header is missing or invalid",
        status_code=401,
        headers={"WWW-Authenticate": 'BASIC realm="You should provide Authorization header"'},
    )

# ============================================================
# Knife4j 文档相关路由
# ============================================================

@app.get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config():
    """Knife4j 所需的 Swagger 配置端点"""
    return {
        "urls": [{"url": f"{ROOT_PATH}/openapi.json", "name": "default"}],
        "configUrl": f"{ROOT_PATH}/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

@app.get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    """Knife4j 文档入口页面"""
    return FileResponse(os.path.join(STATIC_DIR, "doc.html"), media_type="text/html")

# ============================================================
# 系统端点（直接挂载到 app）
# ============================================================

@app.get("/version", tags=["系统"], summary='平台运行信息')
async def version():
    return {
        'title': app.title,
        'description': app.description,
        'version': app.version,
        'startTime': START_TIME.strftime('%Y-%m-%d %H:%M:%S'),
        'Datetime': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    }

@app.get("/health", tags=["系统"], summary="健康检查")
async def health_check():
    return {"status": "ok", "service": "sanxiaBackend"}

# ============================================================
# 业务路由（从 routers/ 目录引入）
# ============================================================
from routers.test import router as test_router

app.include_router(
    test_router,
    prefix="/tp",
)

# ============================================================
# 静态资源兜底路由（必须放在最后）
# ============================================================
app.mount("/", StaticFiles(directory=STATIC_DIR, html=True), name="static")

# ============================================================
# 启动入口
# ============================================================
if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=PORT)
