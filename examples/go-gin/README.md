# Go + Gin + Knife4j Vue3 示例

基于 Gin 框架集成 Knife4j Vue3 API 文档界面，使用 `go:embed` 将前端静态文件嵌入到 Go 二进制中，**编译后单文件部署，无需额外静态文件目录**。

## 前提条件

- Go 1.22+
- 已编译的 Knife4j Vue3 前端（项目根目录的 `dist/` 目录）

## 快速开始

> **注意**：`static/` 目录在 `.gitignore` 中被排除，git 仓库不包含前端静态文件，需要从项目根目录的 `dist/` 复制。

### 1. 编译前端并复制静态文件

```bash
# 在项目根目录编译前端
cd knife4j-vue3
pnpm install && pnpm build

# 复制前端产物到 Go 示例的 static 目录
cp -r dist/* examples/go-gin/static/
```

> Windows 用户也可使用 `examples/setup.bat` 一键复制。

### 2. 初始化 Go 模块

```bash
cd examples/go-gin
go mod tidy
```

### 3. 编译（将前端嵌入二进制）

```bash
# 编译后会生成包含前端的自包含二进制
go build -o server main.go
```

### 4. 启动服务

```bash
# 方式一：直接运行
go run main.go

# 方式二：运行编译后的二进制（可拷贝到任意机器运行，无需额外文件）
./server
```

### 5. 访问文档

打开浏览器访问：http://localhost:8080/doc.html

OpenAPI JSON：http://localhost:8080/swagger/doc.json

## 项目结构

```
go-gin/
├── main.go           # 主程序（Gin 路由 + go:embed 嵌入静态文件）
├── go.mod            # Go 模块定义
├── go.sum            # 依赖校验
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
| GET | /api/users/:id | 根据 ID 获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |
| GET | /api/health | 健康检查 |

## 技术栈

- **Go 1.22+**
- **Gin v1.10.0** - Web 框架
- **embed** - 前端静态文件嵌入
- **swaggo** - Swagger 注解（仅编译时依赖）
- **Knife4j Vue3** - API 文档 UI
