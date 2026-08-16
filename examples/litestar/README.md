# Knife4j Vue3 + LiteStar 示例

基于 LiteStar 框架集成 Knife4j Vue3 API 文档界面。

## 前置准备：生成前端产物

本示例的 `static/` 目录存放 Knife4j Vue3 编译产物，**不纳入版本控制**，需要先手动生成：

```bash
# 1. 在项目根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build
# 产物输出到 dist/ 目录

# 2. 复制到本示例的 static/ 目录
cp -r dist/* examples/litestar/static/
# Windows PowerShell:
# Copy-Item -Path dist\* -Destination examples\litestar\static\ -Recurse -Force
```

`static/` 目录应包含：

```
static/
├── doc.html        # Knife4j 入口页面
├── favicon.ico
├── robots.txt
├── webjars/        # JS/CSS 静态资源
└── oauth/          # OAuth2 授权页面
```

## 快速开始

```bash
# 1. 创建虚拟环境（可选）
python -m venv venv
venv\Scripts\activate        # Windows

# 2. 安装依赖
pip install -r requirements.txt

# 3. 启动服务
uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api
```

服务启动后访问：

- **Knife4j 文档**：http://localhost:8000/api/doc.html
- **OpenAPI JSON**：http://localhost:8000/api/schema/openapi.json

## 项目结构

```
litestar/
├── main.py              # 应用入口
├── routers/             # 路由模块
│   ├── user.py          # 用户管理
│   └── system.py        # 系统端点
├── models/              # 数据模型
│   └── user.py          # 用户模型
├── docs/                # 开发文档
├── requirements.txt     # 依赖
└── README.md
```

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 获取单个用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/health | 健康检查 |

## 技术栈

- **Python 3.10+**
- **LiteStar 2.0+** - ASGI 框架
- **Uvicorn** - ASGI 服务器
- **Pydantic v2** - 数据模型
- **Knife4j Vue3** - API 文档 UI

## 开发指南

详见 [docs/开发指南.md](docs/开发指南.md)
