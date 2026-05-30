# Knife4j Vue3 示例项目

本目录包含 Knife4j Vue3 与不同后端框架集成的完整示例。

## 目录结构

```
examples/
├── java-springboot/    # Java Spring Boot 示例
├── fastapi/            # Python FastAPI 示例
├── litestar/           # Python LiteStar 示例
└── README.md           # 本文件
```

## 快速开始

### 前提条件

先编译 Knife4j Vue3 前端：

```bash
cd knife4j-vue3
pnpm install
pnpm build
```

### 选择你的后端框架

#### Java Spring Boot

```bash
cd examples/java-springboot

# 将前端产物复制到 static 目录
cp -r ../../dist/* src/main/resources/static/

# 启动服务
mvn spring-boot:run
```

访问 http://localhost:8080/doc.html

#### Python FastAPI

```bash
cd examples/fastapi

# 安装依赖
pip install -r requirements.txt

# 将前端产物复制到 static 目录
cp -r ../../dist/* static/

# 启动服务
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

访问 http://localhost:8000/doc.html

#### Python LiteStar

```bash
cd examples/litestar

# 安装依赖
pip install -r requirements.txt

# 将前端产物复制到 static 目录
cp -r ../../dist/* static/

# 启动服务
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

访问 http://localhost:8000/doc.html

## 示例功能

所有示例都包含相同的 API：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 根据ID获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/health | 健康检查 |

## 注意事项

1. **CORS 配置**：所有示例都已配置 CORS，允许跨域请求
2. **Swagger Config**：每个示例都提供了 Knife4j 需要的 `/v3/api-docs/swagger-config` 端点
3. **静态文件托管**：示例默认将前端产物放在 `static/` 目录，由后端统一托管
4. **开发环境**：开发时可以使用 `vite.config.js` 的代理配置，无需复制静态文件
