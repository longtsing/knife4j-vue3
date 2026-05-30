# Python LiteStar 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Python LiteStar ASGI 框架后端集成。

## 什么是 LiteStar？

LiteStar 是一个高性能的 Python ASGI 框架，原生支持 OpenAPI 3.0 文档生成。与 FastAPI 类似，但基于 ASGI 协议，提供更好的异步性能。

## 前提条件

- Python 3.10+
- LiteStar
- uvicorn

## 1. 安装依赖

```bash
pip install litestar uvicorn
```

## 2. 创建 LiteStar 应用

```python
# main.py
from litestar import Litestar, get, post, put, delete
from litestar.response import Response
from litestar.openapi import OpenAPIConfig
from litestar.middleware.cors import CORSConfig
from pydantic import BaseModel
from typing import List, Optional

# 数据模型
class User(BaseModel):
    id: Optional[int] = None
    name: str
    email: str

class UserCreate(BaseModel):
    name: str
    email: str

# 模拟数据库
users_db: List[User] = []
next_id = 1

# API 路由处理器
@get("/api/users", tags=["用户管理"])
async def list_users() -> List[User]:
    """获取所有用户"""
    return users_db

@post("/api/users", tags=["用户管理"])
async def create_user(data: UserCreate) -> User:
    """创建新用户"""
    global next_id
    user = User(id=next_id, **data.model_dump())
    next_id += 1
    users_db.append(user)
    return user

@get("/api/users/{user_id:int}", tags=["用户管理"])
async def get_user(user_id: int) -> Optional[User]:
    """根据ID获取用户"""
    for user in users_db:
        if user.id == user_id:
            return user
    return None

@put("/api/users/{user_id:int}", tags=["用户管理"])
async def update_user(user_id: int, data: UserCreate) -> Optional[User]:
    """更新用户信息"""
    for i, user in enumerate(users_db):
        if user.id == user_id:
            updated = User(id=user_id, **data.model_dump())
            users_db[i] = updated
            return updated
    return None

@delete("/api/users/{user_id:int}", tags=["用户管理"])
async def delete_user(user_id: int) -> bool:
    """删除用户"""
    for i, user in enumerate(users_db):
        if user.id == user_id:
            users_db.pop(i)
            return True
    return False

# 创建应用实例
app = Litestar(
    route_handlers=[
        list_users,
        create_user,
        get_user,
        update_user,
        delete_user,
    ],
    openapi_config=OpenAPIConfig(
        title="My LiteStar API",
        version="1.0.0",
        description="API Documentation with Knife4j Vue3",
    ),
    cors_config=CORSConfig(
        allow_origins=["*"],
        allow_methods=["*"],
        allow_headers=["*"],
        allow_credentials=True,
    ),
)
```

## 3. 配置 Knife4j 端点

Knife4j 需要特定的 `swagger-config` 端点。LiteStar 默认不提供此端点，需要手动创建：

### 方式一：使用 LiteStar 内置 OpenAPI（推荐）

LiteStar 自动将 OpenAPI 规范暴露在 `/schema/openapi.json`。只需添加 swagger-config 端点：

```python
from litestar import get
from litestar.response import Response

@get("/v3/api-docs/swagger-config")
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

@get("/v3/api-docs")
async def api_docs() -> dict:
    """返回 OpenAPI 规范（直接转发 LiteStar 生成的规范）"""
    # LiteStar 的 OpenAPI 规范在 /schema/openapi.json
    # 此端点作为兼容性适配
    import httpx
    async with httpx.AsyncClient() as client:
        response = await client.get("http://localhost:8000/schema/openapi.json")
        return response.json()
```

### 方式二：完整自定义

```python
from litestar import Litestar, get
from litestar.openapi import OpenAPIConfig

# 自定义 OpenAPI 生成器
def custom_openapi_schema():
    """生成自定义的 OpenAPI 规范"""
    from litestar.openapi.spec import Info, OpenAPI
    from litestar.handlers import HTTPRoute
    
    info = Info(
        title="My LiteStar API",
        version="1.0.0",
        description="Complete API documentation"
    )
    
    # 手动构建 paths
    paths = {}
    # ... 根据路由处理器构建 paths
    
    return OpenAPI(
        info=info,
        paths=paths,
        components=None
    )

@get("/v3/api-docs/swagger-config")
async def swagger_config() -> dict:
    return {
        "urls": [{"url": "/v3/api-docs", "name": "default"}],
        "configUrl": "/v3/api-docs/swagger-config"
    }

@get("/v3/api-docs")
async def api_docs() -> dict:
    return custom_openapi_schema()
```

## 4. 运行服务

```bash
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

## 5. 前端代理配置

在 `vite.config.js` 中配置代理：

```javascript
export default defineConfig({
  server: {
    proxy: {
      // Knife4j swagger-config 端点
      '/v3/api-docs': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      // LiteStar OpenAPI 规范
      '/schema/openapi.json': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      // API 路由
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
```

## 6. 启动前端

```bash
pnpm dev
```

访问 `http://localhost:5173/doc.html`

---

## 完整目录结构

```
my-litestar-app/
├── main.py
├── requirements.txt
└── README.md
```

## requirements.txt

```
litestar>=2.0
uvicorn>=0.23.0
pydantic>=2.0
httpx>=0.24.0  # 用于内部 HTTP 请求（可选）
```

---

## 高级配置

### 分组 API 文档

```python
from litestar import Litestar, get
from litestar.tags import Tag

# 使用 tags 分组
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
            Tag(name="用户管理", description="用户相关接口"),
            Tag(name="订单管理", description="订单相关接口"),
        ]
    ),
)
```

### 添加认证中间件

```python
from litestar import Litestar, Request
from litestar.middleware import BaseHTTPMiddleware
from litestar.response import Response

class AuthMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        # 检查认证头
        auth_header = request.headers.get("Authorization")
        if not auth_header or not auth_header.startswith("Bearer "):
            return Response(
                content={"error": "Unauthorized"},
                status_code=401
            )
        
        # 验证 token
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

在 Knife4j 界面中，通过「文档管理」→「全局参数设置」添加 Authorization Header。

### 使用 Extensions 扩展

LiteStar 支持 OpenAPI Extensions，可以在路由处理器中添加自定义元数据：

```python
from litestar import get
from litestar.openapi.spec import OpenAPI

@get(
    "/api/users",
    tags=["用户管理"],
    summary="获取用户列表",
    description="返回系统中所有用户的信息",
    # 使用 extensions 添加自定义属性
    openapi_extra={
        "x-author": "张三",
        "x-order": "1000"
    }
)
async def list_users():
    return []
```

---

## LiteStar vs FastAPI 对比

| 特性 | LiteStar | FastAPI |
|------|----------|---------|
| ASGI 协议 | ✅ 原生支持 | ✅ 原生支持 |
| OpenAPI 3.0 | ✅ 自动生成 | ✅ 自动生成 |
| WebSocket | ✅ | ✅ |
| 依赖注入 | ✅ | ✅ |
| 中间件 | ✅ | ✅ |
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 生态系统 | 🔄 成长中 | ⭐⭐⭐⭐⭐ |
| 文档质量 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## 常见问题

### Q: Knife4j 显示 "No API definitions found"？

**A:** 检查 swagger-config 端点是否正确：

```bash
curl http://localhost:8000/v3/api-docs/swagger-config
```

应该返回：
```json
{
  "urls": [{"url": "/schema/openapi.json", "name": "default"}]
}
```

### Q: 如何自定义 OpenAPI 规范？

**A:** 使用 `OpenAPIConfig` 的 `openapi_config` 参数：

```python
from litestar.openapi import OpenAPIConfig
from litestar.openapi.spec import Info, Contact, License

app = Litestar(
    openapi_config=OpenAPIConfig(
        info=Info(
            title="My API",
            version="1.0.0",
            description="API Documentation",
            contact=Contact(name="Support", email="support@example.com"),
            license=License(name="MIT", url="https://opensource.org/licenses/MIT")
        )
    )
)
```

### Q: 如何添加请求示例？

**A:** 在 Pydantic 模型中使用 `examples`：

```python
from pydantic import BaseModel, Field

class User(BaseModel):
    name: str = Field(description="用户名", examples=["张三"])
    email: str = Field(description="邮箱", examples=["zhangsan@example.com"])
    
    class Config:
        json_schema_extra = {
            "example": {
                "name": "张三",
                "email": "zhangsan@example.com"
            }
        }
```

### Q: 开发环境使用哪个端口？

**A:** LiteStar 默认使用 uvicorn 的 8000 端口，前端 vite 开发服务器使用 5173 端口。确保 `vite.config.js` 中的代理配置指向正确的后端端口。

---

## 生产环境编译与部署

### 1. 编译前端项目

```bash
# 进入前端项目目录
cd knife4j-vue3

# 安装依赖
pnpm install

# 编译生产版本
pnpm build
```

编译产物在 `dist/` 目录下，包含：
- `doc.html` — Knife4j 入口页面
- `webjars/` — JS/CSS 静态资源

### 2. 方式一：LiteStar 直接托管静态文件（推荐）

将前端产物放入 LiteStar 的静态文件目录，由 LiteStar 统一托管。

#### 项目结构

```
my-litestar-app/
├── main.py
├── static/
│   ├── doc.html
│   └── webjars/
│       ├── js/
│       └── css/
├── requirements.txt
└── Dockerfile
```

#### 代码配置

```python
from litestar import Litestar, get
from litestar.response import Response
from litestar.handlers import HTTPRoute
from litestar.static_files import StaticFiles
import os

# 静态文件目录
STATIC_DIR = os.path.join(os.path.dirname(__file__), "static")

@get("/doc.html", include_in_schema=False)
async def knife4j_ui() -> Response:
    """Knife4j 入口页面"""
    file_path = os.path.join(STATIC_DIR, "doc.html")
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()
    return Response(
        content=content,
        media_type="text/html",
    )

# 其他静态资源回退
@get("/{path:path}", include_in_schema=False)
async def serve_static(path: str) -> Response:
    """静态资源回退路由"""
    file_path = os.path.join(STATIC_DIR, path)
    if os.path.isfile(file_path):
        # 根据扩展名确定 media_type
        media_type = "text/plain"
        if path.endswith(".js"):
            media_type = "application/javascript"
        elif path.endswith(".css"):
            media_type = "text/css"
        elif path.endswith(".html"):
            media_type = "text/html"
        elif path.endswith(".json"):
            media_type = "application/json"
        
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()
        return Response(content=content, media_type=media_type)
    
    # 回退到 doc.html
    return await knife4j_ui()

app = Litestar(
    route_handlers=[
        # ... 你的 API 路由
        knife4j_ui,
        serve_static,
    ],
    openapi_config=OpenAPIConfig(
        title="My LiteStar API",
        version="1.0.0",
    ),
)
```

#### 复制前端产物

```bash
# 复制编译产物到 LiteStar 静态目录
cp -r knife4j-vue3/dist/* my-litestar-app/static/
```

启动后访问：`http://your-server:8000/doc.html`

### 3. 方式二：Nginx 反向代理

#### Nginx 配置

```nginx
server {
    listen 80;
    server_name docs.example.com;

    # Knife4j 前端静态资源
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

    # OpenAPI 规范代理
    location /v3/api-docs/ {
        proxy_pass http://127.0.0.1:8000/v3/api-docs/;
        proxy_set_header Host $host;
    }

    location /schema/openapi.json {
        proxy_pass http://127.0.0.1:8000/schema/openapi.json;
        proxy_set_header Host $host;
    }
}
```

#### 部署步骤

```bash
# 1. 编译前端
pnpm build

# 2. 复制到 Nginx 目录
sudo cp -r dist/* /var/www/knife4j-vue3/

# 3. 测试并重载 Nginx
sudo nginx -t
sudo systemctl reload nginx
```

### 4. 方式三：Docker 部署

#### Dockerfile（前端 + 后端一体化）

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

# 安装 Python 依赖
COPY my-litestar-app/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 复制后端代码
COPY my-litestar-app/ ./

# 复制前端编译产物到静态目录
COPY --from=frontend /app/dist ./static/

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]
```

#### 构建与运行

```bash
# 构建镜像
docker build -t my-litestar-knife4j .

# 运行容器
docker run -d -p 8000:8000 --name litestar-app my-litestar-knife4j
```

### 5. 使用 Gunicorn + Uvicorn 部署（推荐生产环境）

```bash
# 安装 gunicorn
pip install gunicorn

# 启动（4 个 worker）
gunicorn main:app \
  -k uvicorn.workers.UvicornWorker \
  -w 4 \
  --bind 0.0.0.0:8000 \
  --access-logfile - \
  --error-logfile -
```

### 6. 验证部署

```bash
# 检查前端页面
curl -I http://your-server:8000/doc.html

# 检查 OpenAPI 端点
curl http://your-server:8000/v3/api-docs/swagger-config

# 检查 API 是否正常
curl http://your-server:8000/api/users
```
