# Python FastAPI 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Python FastAPI 后端集成。

## 前提条件

- Python 3.8+
- FastAPI
- uvicorn

## 1. 安装依赖

```bash
pip install fastapi uvicorn
```

## 2. 创建 FastAPI 应用

```python
# main.py
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List

app = FastAPI(
    title="My API",
    description="API Documentation with Knife4j Vue3",
    version="1.0.0"
)

# 配置 CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 数据模型
class User(BaseModel):
    id: int = None
    name: str
    email: str

# 模拟数据库
users_db: List[User] = []

# API 路由
@app.get("/api/users", response_model=List[User], tags=["用户管理"])
def list_users():
    """获取所有用户"""
    return users_db

@app.post("/api/users", response_model=User, tags=["用户管理"])
def create_user(user: User):
    """创建新用户"""
    user.id = len(users_db) + 1
    users_db.append(user)
    return user

@app.get("/api/users/{user_id}", response_model=User, tags=["用户管理"])
def get_user(user_id: int):
    """根据ID获取用户"""
    for user in users_db:
        if user.id == user_id:
            return user
    return None
```

## 3. 配置 OpenAPI 端点

FastAPI 自动生成 OpenAPI 3.0 规范，但 Knife4j 需要特定的 swagger-config 端点。

### 方式一：直接使用 FastAPI 自带的 OpenAPI（推荐）

FastAPI 默认提供以下端点：
- `/openapi.json` - OpenAPI 3.0 规范
- `/docs` - Swagger UI（可选）
- `/redoc` - ReDoc（可选）

Knife4j 需要 `swagger-config` 端点，需要额外创建：

```python
# 添加到 main.py
@app.get("/v3/api-docs/swagger-config")
def swagger_config():
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

@app.get("/v3/api-docs")
def api_docs():
    """返回 OpenAPI 规范"""
    from fastapi.openapi.utils import get_openapi
    return get_openapi(
        title=app.title,
        version=app.version,
        description=app.description,
        routes=app.routes,
    )
```

### 方式二：使用自定义 OpenAPI 路径

```python
from fastapi.openapi.docs import get_swagger_ui_html

@app.get("/doc.html", include_in_schema=False)
async def custom_swagger_ui():
    """自定义 Knife4j 界面入口"""
    return get_swagger_ui_html(
        openapi_url="/v3/api-docs",
        title="API Documentation",
        swagger_js_url="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js",
        swagger_css_url="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css",
    )
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
      '/v3/api-docs': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      '/openapi.json': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
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
my-fastapi-app/
├── main.py
├── requirements.txt
└── README.md
```

## requirements.txt

```
fastapi>=0.100.0
uvicorn>=0.23.0
pydantic>=2.0
```

---

## 高级配置

### 分组 API 文档

```python
from fastapi import APIRouter

# 创建路由分组
user_router = APIRouter(prefix="/api/users", tags=["用户管理"])
order_router = APIRouter(prefix="/api/orders", tags=["订单管理"])

@user_router.get("/")
def list_users():
    return []

@order_router.get("/")
def list_orders():
    return []

# 注册路由
app.include_router(user_router)
app.include_router(order_router)
```

### 添加认证

```python
from fastapi import Depends, HTTPException
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

security = HTTPBearer()

@app.get("/api/protected")
def protected_endpoint(credentials: HTTPAuthorizationCredentials = Depends(security)):
    """需要认证的接口"""
    if credentials.credentials != "valid-token":
        raise HTTPException(status_code=401, detail="Invalid token")
    return {"message": "success"}
```

在 Knife4j 界面中，通过「文档管理」→「全局参数设置」添加 Authorization Header。

### 自定义 Schema

```python
from pydantic import Field
from enum import Enum

class UserRole(str, Enum):
    admin = "admin"
    user = "user"
    guest = "guest"

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

### 2. 方式一：FastAPI 直接托管静态文件（推荐）

将前端产物放入 FastAPI 的静态文件目录，由 FastAPI 统一托管：

#### 项目结构

```
my-fastapi-app/
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
from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse
import os

app = FastAPI()

# 挂载静态文件目录
static_dir = os.path.join(os.path.dirname(__file__), "static")
app.mount("/static", StaticFiles(directory=static_dir), name="static")

# Knife4j 入口页面
@app.get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    return FileResponse(os.path.join(static_dir, "doc.html"))

# 其他静态资源的回退路由
@app("/{full_path:path}", include_in_schema=False)
async def serve_static(full_path: str):
    file_path = os.path.join(static_dir, full_path)
    if os.path.isfile(file_path):
        return FileResponse(file_path)
    return FileResponse(os.path.join(static_dir, "doc.html"))
```

#### 复制前端产物

```bash
# 复制编译产物到 FastAPI 静态目录
cp -r knife4j-vue3/dist/* my-fastapi-app/static/
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

    location /openapi.json {
        proxy_pass http://127.0.0.1:8000/openapi.json;
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

#### Dockerfile

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
COPY my-fastapi-app/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 复制后端代码
COPY my-fastapi-app/ ./

# 复制前端编译产物到静态目录
COPY --from=frontend /app/dist ./static/

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]
```

#### 构建与运行

```bash
# 构建镜像
docker build -t my-fastapi-knife4j .

# 运行容器
docker run -d -p 8000:8000 --name fastapi-app my-fastapi-knife4j
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

---

## 常见问题

### Q: Knife4j 显示 "No API definitions found"？

**A:** 检查 swagger-config 端点是否正确返回：

```bash
curl http://localhost:8000/v3/api-docs/swagger-config
```

应该返回类似：
```json
{
  "urls": [{"url": "/openapi.json", "name": "default"}]
}
```

### Q: 如何处理 /api 前缀？

**A:** FastAPI 的路由已经包含前缀，在 `vite.config.js` 的代理 rewrite 中去除 `/api` 前缀：

```javascript
rewrite: (path) => path.replace(/^\/api/, '')
```

### Q: 接口返回 422 错误？

**A:** 这是 Pydantic 验证错误，检查请求体是否符合模型定义。Knife4j 的 Debug 调试栏会显示详细的验证错误信息。

---

## 与 Swagger UI 的区别

| 特性 | Knife4j Vue3 | Swagger UI |
|------|-------------|------------|
| 界面美观度 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| 在线调试 | ✅ 增强版 | ✅ 基础版 |
| Markdown 文档 | ✅ | ❌ |
| 全局搜索 | ✅ | ❌ |
| 参数缓存 | ✅ | ❌ |
| cURL 生成 | ✅ | ❌ |
| 多语言 | ✅ | 部分 |
