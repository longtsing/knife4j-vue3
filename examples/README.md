# Knife4j Vue3 示例项目

本目录包含 Knife4j Vue3 与 **5 个不同后端框架** 集成的完整示例。

## 总览

| 示例 | 语言 | 框架 | 端口 | 文档地址 | OpenAPI 端点 | 前端产物路径 |
|------|------|------|------|----------|------------|------------|
| [java-springboot](./java-springboot/) | Java 17 | Spring Boot 3.2 | 8080 | http://localhost:8080/api/doc.html | `/api/v3/api-docs` | `src/main/resources/static/`（Maven 自动复制） |
| [fastapi](./fastapi/) | Python 3.10+ | FastAPI | 8000 | http://localhost:8000/doc.html | `/api/openapi.json` | `static/`（手动复制） |
| [litestar](./litestar/) | Python 3.10+ | LiteStar 2.0+ | 8000 | http://localhost:8000/api/doc.html | `/api/openapi.json` | `static/`（手动复制） |
| [go-gin](./go-gin/) | Go 1.22+ | Gin | 8080 | http://localhost:8080/doc.html | `/swagger/doc.json` | `static/`（`go:embed` 嵌入） |
| [go-stdlib](./go-stdlib/) | Go 1.18+ | 标准库（零依赖） | 8080 | http://localhost:8080/doc.html | `/swagger/doc.json` | `static/`（`go:embed` 嵌入） |

> **注意**：Java 项目把 `static/` 放在 `src/main/resources/` 下（Spring Boot 标准位置），其他项目放在项目根目录。

## 前置条件：编译 Knife4j 前端

所有示例都需要先编译 Knife4j Vue3 前端产物：

```bash
# 在项目根目录执行
cd /path/to/knife4j-vue3
pnpm install
pnpm build
```

编译完成后 `dist/` 目录会生成：

```
dist/
├── doc.html        # Knife4j 入口页面
├── favicon.ico
├── robots.txt
├── webjars/        # JS/CSS 静态资源
└── oauth/          # OAuth2 授权页面
```

## 一键复制脚本（Windows）

```powershell
# 在 knife4j-vue3 根目录执行
$dist = "$PWD\dist"
foreach ($proj in @("fastapi", "litestar", "java-springboot")) {
    $dest = "$PWD\examples\$proj\"
    if ($proj -eq "java-springboot") { $dest = "$dest\src\main\resources\" }
    $dest = $dest + "static"
    if (Test-Path $dest) { Remove-Item $dest -Recurse -Force }
    New-Item -ItemType Directory -Path $dest -Force | Out-Null
    Copy-Item -Path "$dist\*" -Destination $dest -Recurse -Force
    Write-Host "[$proj] copied to $dest"
}
```

## 各示例运行方式

### Java Spring Boot

```bash
cd examples/java-springboot
# 前置：根目录 pnpm build 生成 dist/，并复制到 src/main/resources/static/
mvn clean spring-boot:run
```

访问：http://localhost:8080/api/doc.html

> Maven 的 `maven-resources-plugin` 会自动从 `src/main/resources/static/` 复制到 `target/classes/static/api/`。

### Python FastAPI

```bash
cd examples/fastapi
# 前置：根目录 pnpm build 生成 dist/，并复制到 static/
python -m venv venv
venv\Scripts\activate          # Windows
# source venv/bin/activate    # Linux/Mac
pip install -r requirements.txt
python main.py
```

访问：http://localhost:8000/doc.html

> 文档页面需要 Basic Auth 认证：`admin / admin12345`

### Python LiteStar

```bash
cd examples/litestar
# 前置：根目录 pnpm build 生成 dist/，并复制到 static/
python -m venv venv
venv\Scripts\activate          # Windows
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api
```

访问：http://localhost:8000/api/doc.html

### Go Gin

```bash
cd examples/go-gin
# 前置：根目录 pnpm build 生成 dist/，并复制到 static/
go mod tidy
go run main.go
```

访问：http://localhost:8080/doc.html

> 前端通过 `//go:embed all:static` 嵌入二进制，部署时只需一个文件。

### Go 标准库（零依赖）

```bash
cd examples/go-stdlib
# 前置：根目录 pnpm build 生成 dist/，并复制到 static/
go run main.go
```

访问：http://localhost:8080/doc.html

> 同样的 `//go:embed` 方案，单二进制部署。

## 各示例的 OpenAPI 方案

| 示例 | OpenAPI 来源 | 工具链 |
|------|------------|--------|
| Java Spring Boot | SpringDoc 自动扫描 `@Operation` 等注解 | springdoc-openapi |
| FastAPI | FastAPI 自动扫描 `@app.get` 等装饰器 | starlette 内置 |
| LiteStar | LiteStar 自动扫描 `@get` 等装饰器 | litestar 内置 |
| Go Gin | **运行时构造**（[docs/openapi.go](./go-gin/docs/openapi.go)） | 无 |
| Go 标准库 | **运行时构造**（[docs/openapi.go](./go-stdlib/docs/openapi.go)） | 无 |

Go 示例采用 BrightVideoV2 风格的"运行时构造 OpenAPI"方案：
- 在 `docs/openapi.go` 用结构体表示 OpenAPI 3.0
- 通过 `docs.Register("GET", "/api/users", op)` 链式声明端点
- `docs.MustJSON()` 序列化输出 JSON

## 统一 API 接口

所有示例提供相同的 6 个 API 端点：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users` | 获取用户列表 |
| GET | `/api/users/{id}` | 根据 ID 获取用户 |
| POST | `/api/users` | 创建用户 |
| PUT | `/api/users/{id}` | 更新用户 |
| DELETE | `/api/users/{id}` | 删除用户 |
| GET | `/api/health` | 健康检查 |

## Knife4j 集成原理（所有示例一致）

```
浏览器访问 {prefix}/doc.html
  ↓
doc.html 自带脚本检测当前路径前缀 → apiBasePath = '/api'
  ↓
请求 {prefix}/v3/api-docs/swagger-config → 拿到 OpenAPI JSON 地址
  ↓
请求 OpenAPI JSON 端点（每个示例不同）
  ↓
Knife4j 渲染 API 文档界面
```

调试时，Knife4j 内部 ajax 会自动带上 `{prefix}` 前缀（通过 `apiBasePath`）。

> **2026-08 修复**：Knife4j 前端修复了**重复前缀 bug**。即使 swagger-config 返回的 URL 已包含完整前缀（如 `/api/openapi.json`），也不会再被拼接成 `/api/api/openapi.json`。

## 常见问题

### Q1：端口冲突？

Java Spring Boot 和 Go 示例默认都用 8080 端口，Python 用 8000。

- 启动 Java 后想启动 Go：先 `mvn spring-boot:run -Dspring-boot.run.arguments=--server.port=8090`
- 启动 Go 后想启动 Java：先修改 `examples/go-gin/main.go` 里的 `r.Run(":9090")`

### Q2：前端资源 404？

每个示例的 `static/` 目录都不在版本控制里（参见各项目 `.gitignore`）。请先按上面"前置条件"步骤编译前端并复制。

### Q3：调试面板响应区空白？

检查浏览器 Network：
- swagger-config 应返回 200 + JSON
- OpenAPI JSON 端点应返回 200 + OpenAPI 规范
- 调试请求 URL 应为 `/api/xxx`，**不应**是 `/api/api/xxx`

如果仍异常，请确认前端是最新版本（含 2026-08 的重复前缀 bug 修复）。

### Q4：CORS 报错？

所有示例都配置了 CORS 全开（`allow_origins=["*"]`）。如果仍有报错，请检查浏览器缓存或代理配置。

### Q5：Python 提示 `ModuleNotFoundError: litestar`？

```bash
pip install -r requirements.txt
# 或单独安装
pip install litestar uvicorn pydantic
```

### Q6：Go 报 `pattern all:static matched no files`？

`go:embed` 要求 `static/` 目录在编译时存在且至少包含一个文件。请先复制前端产物。

## 单独打开某个示例的文档

每个示例目录都有自己的 README，介绍该项目特有的细节：

- [examples/java-springboot/README.md](./java-springboot/README.md)
- [examples/fastapi/README.md](./fastapi/README.md)
- [examples/litestar/README.md](./litestar/README.md)
- [examples/go-gin/README.md](./go-gin/README.md)
- [examples/go-stdlib/README.md](./go-stdlib/README.md)
