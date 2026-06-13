# FastAPI 后端示例项目

基于 FastAPI + Knife4j Vue3 的后端 API 示例项目。

## 快速开始

```bash
# 1. 创建虚拟环境（可选）
python -m venv venv
venv\Scripts\activate        # Windows
# source venv/bin/activate   # Linux/Mac

# 2. 安装依赖
pip install -r requirements.txt

# 3. 启动服务
python main.py
```

服务启动后访问：

- **Knife4j 文档**：http://localhost:8000/doc.html
- **OpenAPI JSON**：http://localhost:8000/api/openapi.json

> **注意**：FastAPI 示例直接引用 `../../dist` 目录，无需手动复制前端产物。

## 认证

访问文档页面需要 Basic Auth 认证：

| 用户名 | 密码 |
|--------|------|
| hxgis | hxgis12345 |
| hbxqx | hbxqx168 |

## 项目结构

```
fastapi/
├── main.py              # 入口：FastAPI 应用配置、中间件、系统端点
├── routers/             # 路由模块（按业务拆分）
│   └── test.py          # 测试路由器示例
├── docs/                # 开发文档
├── requirements.txt     # Python 依赖
└── README.md
```

## 添加新端点

**方式一：直接在 main.py 添加（适合系统级端点）**

```python
@app.get("/my-endpoint", tags=["我的模块"], summary="接口说明")
async def my_endpoint(param: str = Query(description="参数")):
    return {"data": param}
```

**方式二：在 routers/ 中创建新路由器**

```python
# routers/example.py
from fastapi import APIRouter
router = APIRouter(tags=["示例"])

@router.get("/example")
async def example():
    return {"message": "hello"}
```

```python
# main.py 中引入
from routers.example import router as example_router
app.include_router(example_router)
```

详细开发文档请查看 [docs/](docs/) 目录。

## 技术栈

- **FastAPI 0.110.0** - Web 框架
- **Uvicorn 0.28.0** - ASGI 服务器
- **ORJSON 3.9.15** - 高性能 JSON
- **Knife4j Vue3** - API 文档 UI
