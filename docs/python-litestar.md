# Python LiteStar 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Python LiteStar ASGI 框架后端集成。

## 什么是 LiteStar？

LiteStar 是一个高性能的 Python ASGI 框架，原生支持 OpenAPI 3.0 文档生成。基于 ASGI 协议，异步吞吐优于 Flask。

## 前提条件

- Python 3.10+
- LiteStar
- uvicorn
- 编译后的 Knife4j Vue3 前端

## 1. 安装依赖

```bash
pip install litestar uvicorn pydantic
```

## 2. 项目结构

```
my-litestar-app/
├── main.py                 # LiteStar 应用入口
├── routers/                # 业务路由模块
│   ├── user.py
│   └── system.py
├── models/                 # Pydantic 模型
├── static/                 # Knife4j 前端产物（不纳入版本控制）
│   ├── doc.html
│   ├── webjars/
│   └── oauth/
├── requirements.txt
└── README.md
```

## 3. 创建 LiteStar 应用

参考 [examples/litestar/main.py](../examples/litestar/main.py)：

```python
# main.py
from litestar import Litestar, get, Router
from litestar.openapi import OpenAPIConfig
from litestar.openapi.spec import Server
from litestar.config.cors import CORSConfig
from litestar.static_files import create_static_files_router
from litestar.response import File
import os

ROOT_PATH = '/api'                                          # ASGI root_path 前缀
STATIC_DIR = os.path.join(os.path.dirname(__file__), "static")  # Knife4j 前端

# Knife4j swagger-config 端点
@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    return {
        "urls": [{"url": f"{ROOT_PATH}/openapi.json", "name": "default"}],
        "configUrl": f"{ROOT_PATH}/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

# doc.html 入口页面
@get("/doc.html", include_in_schema=False)
async def serve_doc_html() -> File:
    return File(
        path=os.path.join(STATIC_DIR, "doc.html"),
        content_disposition_type="inline",
        media_type="text/html"
    )

# webjars 静态资源
webjars_router = create_static_files_router(
    path="/webjars",
    directories=[os.path.join(STATIC_DIR, "webjars")],
    include_in_schema=False,
)

# oauth 静态资源
oauth_router = create_static_files_router(
    path="/oauth",
    directories=[os.path.join(STATIC_DIR, "oauth")],
    include_in_schema=False,
)

# 业务路由（放在 / 子路径下，Litestar 的 path="/api" 会自动加前缀）
api_router = Router(
    path="/",
    route_handlers=[
        # 在这里挂载你的 handler
    ],
)

app = Litestar(
    path=ROOT_PATH,                     # 所有路由自动加 /api 前缀
    route_handlers=[
        swagger_config,
        serve_doc_html,
        api_router,
        webjars_router,
        oauth_router,
    ],
    openapi_config=OpenAPIConfig(
        title="Knife4J Vue3 LiteStar 示例",
        version="1.0.0",
        path="/",                       # 关键：放在根路径，Litestar 会自动拼 /api
        render_plugins=[],              # 禁用 LiteStar 自带 Scalar/Swagger UI
        servers=[Server(url=ROOT_PATH, description="API 服务")],
    ),
    cors_config=CORSConfig(allow_origins=["*"]),
)
```

## 4. 启动流程

```bash
# 1. 在 knife4j-vue3 根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build

# 2. 复制前端产物到本项目
cp -r dist/* my-litestar-app/static/

# 3. 启动服务
cd my-litestar-app
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api
```

访问：http://localhost:8000/api/doc.html

## 5. Knife4j 集成原理

LiteStar 内置 OpenAPI 3.0 文档生成，会自动在 `{root_path}/schema/openapi.json` 暴露规范。本示例把它改到 `{root_path}/openapi.json`（即 `/api/openapi.json`），更贴近 SpringDoc 风格：

```
浏览器访问 /api/doc.html
  ↓
doc.html 自带脚本检测路径前缀 → apiBasePath = '/api'
  ↓
请求 /api/v3/api-docs/swagger-config → 拿到 OpenAPI JSON 地址 /api/openapi.json
  ↓
请求 /api/openapi.json → LiteStar 自动生成的规范
  ↓
Knife4j 渲染文档界面
```

> **2026-08 修复**：Knife4j 前端修复了**重复前缀 bug**。即使 swagger-config 返回 `/api/openapi.json`（含前缀）也不会再被拼接成 `/api/api/openapi.json`。

## 6. 关于 `path="/api"`

`Litestar(path="/api", ...)` 会让**所有挂在 app 上的路由**自动加上 `/api` 前缀。例如：

```python
@get("/users")                  # 实际路径 /api/users
async def list_users(): ...
```

业务 router 一般用 `Router(path="/")`，让单个 handler 自己写完整路径：

```python
api_router = Router(
    path="/",                    # 子路由不重复加前缀
    route_handlers=[
        list_users,              # 实际路径 /api/list-users
    ],
)
```

## 7. 业务路由示例

```python
# routers/user.py
from litestar import get, post, put, delete
from pydantic import BaseModel

class User(BaseModel):
    id: int | None = None
    name: str
    email: str

@get("/users", tags=["用户管理"])
async def list_users() -> list[User]:
    return []

@get("/users/{user_id:int}", tags=["用户管理"])
async def get_user(user_id: int) -> User | None:
    return None

@post("/users", tags=["用户管理"])
async def create_user(data: User) -> User:
    return User(id=1, **data.model_dump())
```

## 8. 高级配置

### 8.1 认证中间件

```python
from litestar import Request
from litestar.middleware import BaseHTTPMiddleware
from litestar.response import Response

class AuthMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        auth_header = request.headers.get("Authorization")
        if not auth_header or not auth_header.startswith("Bearer "):
            return Response(
                content={"error": "Unauthorized"},
                status_code=401
            )

        token = auth_header.replace("Bearer ", "")
        if token != "valid-token":
            return Response(
                content={"error": "Invalid token"},
                status_code=401
            )

        return await call_next(request)

app = Litestar(
    route_handlers=[...],
    middleware=[AuthMiddleware],
)
```

在 Knife4j 界面中，通过「文档管理」→「全局参数设置」添加 Authorization Header（值 `Bearer valid-token`）。

### 8.2 Extensions 扩展

```python
@get(
    "/api/users",
    tags=["用户管理"],
    summary="获取用户列表",
    description="返回系统中所有用户的信息",
    openapi_extra={
        "x-author": "张三",
        "x-order": "1000"
    }
)
async def list_users():
    return []
```

### 8.3 分组 API 文档

```python
from litestar import Litestar, get
from litestar.openapi import OpenAPIConfig

@get("/api/users", tags=["用户管理"])
async def list_users():
    return []

@get("/api/orders", tags=["订单管理"])
async def list_orders():
    return []

app = Litestar(
    route_handlers=[list_users, list_orders],
    openapi_config=OpenAPIConfig(
        title="API Documentation",
        version="1.0.0",
        tags=[
            {"name": "用户管理", "description": "用户相关接口"},
            {"name": "订单管理", "description": "订单相关接口"},
        ]
    ),
)
```

## 9. 生产环境部署

### 9.1 嵌入式（推荐）

```bash
# 编译前端
cd knife4j-vue3 && pnpm build

# 复制
cp -r dist/* my-litestar-app/static/

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
    location /schema/openapi.json {
        proxy_pass http://127.0.0.1:8000/schema/openapi.json;
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
COPY my-litestar-app/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY my-litestar-app/ ./
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
curl -I http://your-server:8000/api/doc.html
curl http://your-server:8000/api/v3/api-docs/swagger-config
curl http://your-server:8000/api/users
```

## 10. 完整示例

[examples/litestar/](../examples/litestar/) 提供开箱即用的可运行示例：

```bash
cd examples/litestar
# 前置：根目录 pnpm build 生成 dist/
cp -r ../dist/* static/
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api
# 访问 http://localhost:8000/api/doc.html
```

## 11. 常见问题

### Q1：报 `Handler already registered for path '/api/openapi.json' and http method OPTIONS`？

A：Litestar 内置 `/schema/openapi.json` 与自定义端点冲突，且 Knife4j CORS 预检会触发 OPTIONS。修复：

```python
openapi_config=OpenAPIConfig(
    path="/",                # 把 OpenAPI 端点放在根（自动变 /api/openapi.json）
    render_plugins=[],       # 禁用 LiteStar 自带 UI
)
```

如果仍报错，请检查 [examples/litestar/main.py](../examples/litestar/main.py) 的最新版本。

### Q2：调试面板响应区空白？

A：检查浏览器 Network：

- `/api/openapi.json` 应返回 200 + JSON
- 调试请求的 URL 应为 `/api/xxx`，**不应**是 `/api/api/xxx`

Knife4j 前端已在 2026-08 修复重复前缀 bug，swagger-config 返回 `/api/openapi.json`（含前缀）不会再被拼接。

### Q3：Knife4j 显示 "No API definitions found"？

A：检查 swagger-config 端点：

```bash
curl http://localhost:8000/api/v3/api-docs/swagger-config
```

应返回 `urls[0].url = /api/openapi.json`。

### Q4：端口冲突？

A：LiteStar 默认 8000；Java/Go 用 8080。修改 `uvicorn ... --port 9000`。

### Q5：开发端口与部署端口？

A：开发模式用 `--reload`，生产模式用 `--workers 4`。

## 12. LiteStar vs FastAPI 对比

| 特性 | LiteStar | FastAPI |
|------|----------|---------|
| ASGI 协议 | ✅ 原生 | ✅ 原生 |
| OpenAPI 3.0 | ✅ 自动 | ✅ 自动 |
| WebSocket | ✅ | ✅ |
| 依赖注入 | ✅ | ✅ |
| 中间件 | ✅ | ✅ |
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 生态成熟度 | 成长中 | ⭐⭐⭐⭐⭐ |
| 文档质量 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |