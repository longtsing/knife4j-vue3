# Go 标准库 + Knife4j Vue3 示例

**零依赖**：仅使用 Go 标准库（`net/http`、`encoding/json`、`embed`），无需安装任何第三方包。前端静态文件通过 `go:embed` 嵌入到二进制中，**编译后单文件部署**。

## 前提条件

- Go 1.22+
- 已编译的 Knife4j Vue3 前端（项目根目录的 `dist/` 目录）

## 快速开始

> **注意**：`static/` 目录在 `.gitignore` 中被排除，git 仓库不包含前端静态文件，需要从项目根目录的 `dist/` 复制。

```bash
# 1. 编译前端（在项目根目录）
cd ../..
pnpm install && pnpm build

# 2. 复制前端产物到 static 目录
cp -r dist/* examples/go-stdlib/static/

# 3. 编译 Go（将前端嵌入到二进制）
cd examples/go-stdlib
go build -o server main.go

# 4. 运行（零依赖，单文件部署）
./server
```

访问：http://localhost:8080/doc.html

OpenAPI JSON：http://localhost:8080/openapi.json

## 项目结构

```
go-stdlib/
├── main.go           # 主程序（纯标准库，含 go:embed 嵌入）
├── go.mod            # Go 模块定义
├── README.md         # 本文件
└── static/           # Knife4j 前端静态文件（编译时嵌入到二进制）
    ├── doc.html
    ├── webjars/
    ├── oauth/
    ├── favicon.ico
    └── robots.txt
```

## embed 原理

使用 Go 1.16+ 的 `embed` 包，在编译时将 `static/` 目录的内容嵌入到二进制文件中：

```go
//go:embed all:static
var staticFS embed.FS
```

编译后的 `server` 二进制**完全自包含**，可以单独拷贝到任意机器运行，无需携带 `static/` 目录。

> ⚠️ **注意**：如果修改了前端代码，需要重新复制 `dist/*` 到 `static/` 并**重新编译** Go 二进制才能生效。

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
- ✅ **单文件部署**：通过 `go:embed` 将前端嵌入二进制，编译后无需额外文件
- ✅ **手写路由器**：基于 `http.HandleFunc` 的简单路径匹配
- ✅ **完整 CORS**：已配置跨域支持
- ✅ **手写 OpenAPI 3.0.3**：完整规范约 200 行，无需 swag 工具
- ✅ **端口可配**：通过 `PORT` 环境变量修改（默认 8080）

## 配置

通过环境变量修改端口：

```bash
# Linux/Mac
PORT=9090 ./server

# Windows PowerShell
$env:PORT=9090; .\server.exe
```
