# Go 标准库 + Knife4j Vue3 示例

**零依赖**：仅使用 Go 标准库（`net/http`、`encoding/json`），无需安装任何第三方包。

## 前提条件

- Go 1.18+
- 已编译的 Knife4j Vue3 前端（`dist/` 目录）

## 快速开始

```bash
# 1. 编译前端（在项目根目录）
cd ../..
pnpm install && pnpm build

# 2. 复制前端产物到 static 目录
cp -r dist/* examples/go-stdlib/static/

# 3. 启动服务（零依赖，直接运行）
cd examples/go-stdlib
go run main.go
```

访问：http://localhost:8080/doc.html

OpenAPI JSON：http://localhost:8080/openapi.json

## 项目结构

```
go-stdlib/
├── main.go           # 主程序（纯标准库，约 600 行）
├── README.md         # 本文件
└── static/           # Knife4j 前端静态文件（需手动复制）
    ├── doc.html
    └── webjars/
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 根据 ID 获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/health | 健康检查 |

## 特点

- ✅ **零依赖**：不使用 Gin、Echo 等任何第三方框架
- ✅ **手写路由器**：基于 `http.HandleFunc` 的简单路径匹配
- ✅ **完整 CORS**：已配置跨域支持
- ✅ **手写 OpenAPI 3.0.3**：完整规范约 200 行，无需 swag 工具
- ✅ **静态文件服务**：自动检测 `static/` 目录并托管（含 Content-Type 识别）
- ✅ **端口可配**：通过 `PORT` 环境变量修改（默认 8080）

## 配置

通过环境变量修改端口：

```bash
# Linux/Mac
PORT=9090 go run main.go

# Windows PowerShell
$env:PORT=9090; go run main.go
```
