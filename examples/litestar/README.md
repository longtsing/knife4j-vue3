# Knife4j Vue3 + LiteStar 示例

基于 LiteStar 框架集成 Knife4j Vue3 API 文档界面。

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

> **注意**：LiteStar 示例直接引用 `../../dist` 目录，无需手动复制前端产物。

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
