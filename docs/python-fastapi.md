# Python FastAPI 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Python FastAPI 后端集成。

## 前提条件

- Python 3.10+
- FastAPI
- uvicorn
- 编译后的 Knife4j Vue3 前端

## 1. 安装依赖

```bash
pip install fastapi uvicorn[standard]
```

## 2. 项目结构

```
my-fastapi-app/
├── main.py                 # FastAPI 应用入口
├── routers/                # 业务路由模块
├── models/                 # Pydantic 模型
├── static/                 # Knife4j 前端产物（不纳入版本控制）
│   ├── doc.html
│   ├── webjars/
│   └── oauth/
├── requirements.txt
└── README.md
```

`static/` 目录需手动生成（前端编译后从 `dist/*` 复制过来）。

## 3. 创建 FastAPI 应用

参考 [examples/fastapi/main.py](../examples/fastapi/main.py)：

```python
# main.py
import os
from datetime import datetime
from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from starlette.staticfiles import StaticFiles

ROOT_PATH = '/api'                              # 反向代理前缀
STATIC_DIR = os.path.join(os.path.dirname(__file__), "static")  # Knife4j 前端

app = FastAPI(
    root_path=ROOT_PATH,
    title='FastAPI 后端示例项目',
    version='1.0.0',
    docs_url=None,                              # 禁用内置 Swagger UI
    redoc_url=None,                             # 禁用内置 ReDoc
    servers=[{"url": ROOT_PATH, "description": "API 服务"}],
)

# CORS（开发模式全开）
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 业务路由
from routers.test import router as test_router
app.include_router(test_router, prefix="/tp")

# Knife4j swagger-config 端点（必须是 /v3/api-docs/swagger-config）
@app.get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config():
    return {
        "urls": [{"url": f"{ROOT_PATH}/openapi.json", "name": "default"}],
        "configUrl": f"{ROOT_PATH}/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

# Knife4j 入口页面
@app.get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    return FileResponse(os.path.join(STATIC_DIR, "doc.html"), media_type="text/html")

# 静态资源兜底路由（必须放在最后）
app.mount("/", StaticFiles(directory=STATIC_DIR, html=True), name="static")

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
```

## 4. 启动流程

```bash
# 1. 在 knife4j-vue3 根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build

# 2. 复制前端产物到本项目
cp -r dist/* my-fastapi-app/static/

# 3. 启动服务
cd my-fastapi-app
pip install -r requirements.txt
python main.py
```

访问：http://localhost:8000/doc.html

## 5. Knife4j 集成原理

```
浏览器访问 /doc.html
  ↓
doc.html 自带脚本检测当前路径前缀 → apiBasePath = '/api'
  ↓
请求 /v3/api-docs/swagger-config → 拿到 OpenAPI JSON 地址 /api/openapi.json
  ↓
请求 /api/openapi.json → FastAPI 内置的 OpenAPI 规范
  ↓
Knife4j 渲染文档界面
```

调试时，Knife4j 内部 ajax 会自动带上 `/api` 前缀（通过 `apiBasePath`）。

> **2026-08 修复**：Knife4j 前端修复了**重复前缀 bug**。即使 swagger-config 返回 `/api/openapi.json`（含前缀）也不会再被拼接成 `/api/api/openapi.json`。

## 6. 关于 root_path

`root_path` 用于反向代理前缀。如果你的部署架构是：

```
nginx /api/*  →  uvicorn (root_path=/api)
```

那么所有 `app.include_router(prefix='/users')` 的端点都会暴露在 `/api/users`。OpenAPI JSON 中的 `servers[0].url` 也要设为 `/api`。

如果不使用反向代理，可省略 `root_path`：

```python
app = FastAPI()  # 端点直接是 /api/users
```

## 7. 高级配置

### 7.1 文档访问认证（Basic Auth）

```python
@app.middleware("http")
async def auth_middleware(request: Request, call_next):
    protected_paths = [
        f"{ROOT_PATH}/openapi.json",
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

    credentials = {"admin": "admin12345"}
    if credentials.get(username) == password:
        return await call_next(request)

    return Response(
        content="Authorization header is missing or invalid",
        status_code=401,
        headers={"WWW-Authenticate": 'BASIC realm="You should provide Authorization header"'},
    )
```

详见 [examples/fastapi/main.py](../examples/fastapi/main.py) 的完整实现。

### 7.2 Bearer 认证（推荐）

```python
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

security = HTTPBearer()

async def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    if credentials.credentials != "valid-token":
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED)

@app.get("/api/protected", dependencies=[Depends(verify_token)])
async def protected_endpoint():
    return {"message": "success"}
```

在 Knife4j 界面中，通过「文档管理」→「全局参数设置」添加 Authorization Header（值填 `Bearer valid-token`）。

### 7.3 自定义 Schema

```python
from pydantic import BaseModel, Field
from enum import Enum

class UserRole(str, Enum):
    admin = "admin"
    user = "user"

class User(BaseModel):
    """用户模型"""
    id: int = Field(description="用户ID", examples=[1])
    name: str = Field(description="用户名", min_length=1, max_length=50)
    email: str = Field(description="邮箱地址", examples=["user@example.com"])
    role: UserRole = Field(description="用户角色", default=UserRole.user)

    class Config:
        json_schema_extra = {
            "example": {
                "id": 1,
                "name": "张三",
                "email": "zhangsan@example.com",
                "role": "admin"
            }
        }
```

### 7.4 自定义 OpenAPI 元信息

```python
app = FastAPI(
    title="My API",
    description="API Documentation with Knife4j Vue3",
    version="1.0.0",
    contact={"name": "Support", "email": "support@example.com"},
    license_info={"name": "Apache 2.0", "url": "https://..."},
    openapi_tags=[
        {"name": "用户管理", "description": "用户 CRUD 接口"},
        {"name": "订单管理", "description": "订单相关接口"},
    ],
)
```

## 8. 路由分组

```python
from fastapi import APIRouter

user_router = APIRouter(prefix="/api/users", tags=["用户管理"])
order_router = APIRouter(prefix="/api/orders", tags=["订单管理"])

@user_router.get("/")
async def list_users():
    return []

@user_router.get("/{user_id}")
async def get_user(user_id: int):
    return {}

app.include_router(user_router)
app.include_router(order_router)
```

## 9. 生产环境部署

### 9.1 嵌入式（推荐）

直接拷贝 `static/` 随应用一起发布：

```bash
# 编译前端
cd knife4j-vue3 && pnpm build

# 复制
cp -r dist/* my-fastapi-app/static/

# 启动
uvicorn main:app --host 0.0.0.0 --port 8000 --workers 4
```

### 9.2 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name docs.example.com;

    # Knife4j 前端
    location / {
        root /var/www/knife4j-vue3/dist;
        index doc.html;
        try_files $uri $uri/ /doc.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8000/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # OpenAPI 端点
    location /openapi.json {
        proxy_pass http://127.0.0.1:8000/openapi.json;
    }
}
```

### 9.3 Docker 一体化

```dockerfile
# 构建阶段
FROM node:20-alpine AS frontend
WORKDIR /app
COPY knife4j-vue3/package.json knife4j-vue3/pnpm-lock.yaml ./
RUN corepack enable && pnpm install
COPY knife4j-vue3/ ./
RUN pnpm build

# 运行阶段
FROM python:3.12-slim
WORKDIR /app
COPY my-fastapi-app/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY my-fastapi-app/ ./
COPY --from=frontend /app/dist ./static/

EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]
```

### 9.4 Gunicorn + Uvicorn（生产推荐）

```bash
pip install gunicorn
gunicorn main:app \
  -k uvicorn.workers.UvicornWorker \
  -w 4 \
  --bind 0.0.0.0:8000 \
  --access-logfile - \
  --error-logfile -
```

### 9.5 验证部署

```bash
curl -I http://your-server:8000/doc.html
curl http://your-server:8000/v3/api-docs/swagger-config
curl http://your-server:8000/api/users
```

## 10. 完整示例

[examples/fastapi/](../examples/fastapi/) 提供开箱即用的可运行示例：

```bash
cd examples/fastapi
# 前置：根目录 pnpm build 生成 dist/
cp -r ../dist/* static/
pip install -r requirements.txt
python main.py
# 访问 http://localhost:8000/doc.html（需 Basic Auth: admin/admin12345）
```

## 11. 常见问题

### Q1：调试面板响应区空白？

A：检查浏览器 Network：

- `/api/openapi.json` 应返回 200 + JSON
- 调试请求的 URL 应为 `/api/xxx`，**不应**是 `/api/api/xxx`

Knife4j 前端已在 2026-08 修复重复前缀 bug，swagger-config 返回 `/api/openapi.json`（含前缀）不会再被拼接。

### Q2：Knife4j 显示 "No API definitions found"？

A：检查 swagger-config 端点：

```bash
curl http://localhost:8000/v3/api-docs/swagger-config
```

应返回包含 `urls` 数组的 JSON。`url` 字段指向 `/api/openapi.json`。

### Q3：端口冲突？

A：FastAPI 默认 8000；Java Spring Boot 用 8080；Go 用 8080。修改 `uvicorn.run(... port=9000)` 或在 `application.yml` 中修改 `server.port`。

### Q4：接口返回 422 错误？

A：Pydantic 验证失败。检查请求体是否符合模型定义。Knife4j 调试面板会显示详细错误。