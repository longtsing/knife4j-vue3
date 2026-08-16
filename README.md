# Knife4j Vue3 — 跨语言 API 文档增强平台

> **为 Java、Python、Go 等全栈开发者打造的新一代 OpenAPI 可视化与调试平台**

---

## 为什么要做这个项目？

在 Java 生态中，[Knife4j](https://github.com/xiaoymin/knife4j) 已经成为 Spring Boot 项目首选的 API 文档增强方案——美观的界面、强大的调试能力、完善的中文支持，深受国内开发者喜爱。

然而，当我们转向 **Python 生态**（FastAPI、LiteStar、Flask）或其他语言时，会发现一个尴尬的现实：

| 后端框架 | 原生文档 UI | 调试能力 | 中文支持 | 体验评分 |
|----------|------------|---------|---------|---------|
| Spring Boot + Knife4j | Knife4j | ⭐⭐⭐⭐⭐ | ✅ 完美 | ⭐⭐⭐⭐⭐ |
| Spring Boot + Swagger UI | Swagger UI | ⭐⭐⭐ | 部分 | ⭐⭐⭐ |
| FastAPI 原生 Swagger UI | Swagger UI | ⭐⭐⭐ | 部分 | ⭐⭐⭐ |
| LiteStar 原生 ReDoc | ReDoc | ⭐⭐ | 部分 | ⭐⭐⭐ |

**Python 框架的 API 文档体验，与 Java 生态存在代际差距。**

本项目将 Knife4j 的 Vue3 前端核心能力进行**跨语言标准化适配**，使其不再绑定于 Spring 生态，而是成为**任何支持 OpenAPI 3.0 规范的后端框架**都能使用的通用 API 文档增强平台。

---

## 核心能力

### 🎯 一套前端，全栈通用

无论你的后端是 Java、Python、Go、Rust 还是 Node.js，只要遵循 OpenAPI 3.0 规范，即可获得 Knife4j 级别的文档体验。

### 🛠️ 专业级 API 调试

- **在线调试面板**：无需 Postman，直接在文档中发送请求、查看响应
- **自动 cURL 生成**：一键复制请求命令，方便团队协作
- **全局参数管理**：统一配置 Authorization、自定义 Header 等公共参数（**支持修改和删除**）
- **参数缓存**：调试参数自动保存，刷新页面不丢失

### 📖 增强的文档展示

- **多分组管理**：按模块/API Tag 自动分组，支持下拉切换
- **Schema 模型可视化**：复杂数据结构一目了然
- **Markdown 文档集成**：在 API 文档中嵌入富文本说明
- **Extensions 扩展**：支持 `x-author`、`x-order` 等自定义元数据

### 🌍 国际化 & 本土化

- 完整的**中文界面**，符合国内开发者使用习惯
- 支持英文、日文多语言切换
- 符合中国网络环境的部署方案

---

## 支持的后端框架

| 框架 | 语言 | OpenAPI 版本 | OpenAPI 来源 | 集成难度 | 状态 |
|------|------|-------------|-------------|---------|------|
| **Spring Boot** (SpringDoc) | Java | 3.0 | 注解自动生成 | ⭐ 极简 | ✅ 完全支持 |
| **FastAPI** | Python | 3.0 | 装饰器自动生成 | ⭐⭐ 简单 | ✅ 完全支持 |
| **LiteStar** | Python | 3.0 | 装饰器自动生成 | ⭐⭐ 简单 | ✅ 完全支持 |
| **Gin** | Go | 3.0 | **运行时构造** | ⭐⭐ 简单 | ✅ 完全支持 |
| **Go 标准库**（零依赖） | Go | 3.0 | **运行时构造** | ⭐⭐ 简单 | ✅ 完全支持 |
| **JFinal** | Java | 2.0 / 3.0 | 注解自动生成 | ⭐ 极简 | ✅ 完全支持 |
| **通用 OpenAPI 3.0** | 任意语言 | 3.0 | 自定义 | ⭐⭐⭐ 中等 | ✅ 完全支持 |

> 💡 **通用适配原则**：只要你的后端能暴露 `/v3/api-docs/swagger-config` 端点，就能接入本项目。

---

## 功能全景

```
┌─────────────────────────────────────────────────────────┐
│                    Knife4j Vue3                          │
├──────────────┬──────────────┬──────────────┬────────────┤
│   文档展示    │   API 调试    │   数据模型    │   设置管理  │
├──────────────┼──────────────┼──────────────┼────────────┤
│ • 多分组切换  │ • 在线请求    │ • Schema 可视 │ • 全局参数  │
│ • 接口搜索    │ • 响应解析    │ • Model 展开  │ • 主题定制  │
│ • 标签排序    │ • cURL 生成   │ • 示例数据    │ • 语言切换  │
│ • Markdown   │ • 参数缓存    │ • allOf 合并  │ • 文档管理  │
│ • Extensions │ • AfterScript │ • OAS2/3 兼容 │ • Host 配置 │
└──────────────┴──────────────┴──────────────┴────────────┘
```

---

## 快速开始

### 方式一：使用预编译产物（推荐）

直接使用 [releases](https://github.com/longtsing/knife4j-vue3/releases) 中已编译的 `dist/` 目录，部署到你的后端项目。

### 方式二：从源码构建

```bash
# 1. 克隆项目
git clone <repo-url>
cd knife4j-vue3

# 2. 安装依赖
pnpm install

# 3. 编译生产版本
pnpm build
# 产物输出到 dist/ 目录

# 4. 将 dist/ 部署到后端项目
```

### 方式三：开发模式（热更新）

```bash
pnpm dev
# 访问 http://localhost:5173/doc.html
```

编辑 `vite.config.js` 中的代理配置，指向你的后端服务：

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8000',  // 你的后端地址
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '')
  }
}
```

---

## 后端集成示例

### Java Spring Boot

```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.springdoc</groupId>
    <artifactId>springdoc-openapi-starter-webmvc-api</artifactId>
    <version>2.6.0</version>
</dependency>
```

```java
// SwaggerConfigController.java
@RestController
public class SwaggerConfigController {
    @GetMapping("/api/v3/api-docs/swagger-config")
    public Map<String, Object> swaggerConfig() {
        return Map.of(
            "urls", List.of(Map.of("url", "/api/v3/api-docs", "name", "default")),
            "configUrl", "/api/v3/api-docs/swagger-config",
            "validatorUrl", ""
        );
    }
}
```

详细：[docs/java-springboot.md](./docs/java-springboot.md)

### Python FastAPI

```python
from fastapi import FastAPI
from fastapi.responses import FileResponse

app = FastAPI(docs_url=None)  # 禁用默认 Swagger UI

@app.get("/api/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config():
    return {
        "urls": [{"url": "/api/openapi.json", "name": "default"}],
        "configUrl": "/api/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

@app.get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    return FileResponse("static/doc.html", media_type="text/html")
```

详细：[docs/python-fastapi.md](./docs/python-fastapi.md)

### Python LiteStar

```python
from litestar import Litestar, get
from litestar.response import File
from litestar.openapi import OpenAPIConfig

@get("/v3/api-docs/swagger-config", include_in_schema=False)
async def swagger_config() -> dict:
    return {
        "urls": [{"url": "/api/openapi.json", "name": "default"}],
        "configUrl": "/api/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

@get("/doc.html", include_in_schema=False)
async def knife4j_ui() -> File:
    return File(path="static/doc.html", media_type="text/html")

app = Litestar(
    path="/api",
    route_handlers=[swagger_config, knife4j_ui],
    openapi_config=OpenAPIConfig(
        path="/",
        render_plugins=[],   # 禁用 LiteStar 自带 UI
    ),
)
```

详细：[docs/python-litestar.md](./docs/python-litestar.md)

### Go Gin（运行时构造 OpenAPI）

```go
package main

import (
    "github.com/gin-gonic/gin"
    "myapp/docs"
)

func main() {
    r := gin.Default()

    // swagger-config 端点
    r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "urls": []gin.H{{"url": "/swagger/doc.json", "name": "default"}},
            "configUrl": "/v3/api-docs/swagger-config",
            "validatorUrl": "",
        })
    })

    // OpenAPI 端点
    r.GET("/swagger/doc.json", func(c *gin.Context) {
        c.Header("Content-Type", "application/json")
        c.Writer.Write(docs.MustJSON())
    })

    // 嵌入前端（编译后单二进制部署）
    r.GET("/doc.html", ...)
    r.StaticFS("/webjars", http.FS(subFS))

    r.Run(":8080")
}
```

`docs/openapi.go` 使用结构体表示 OpenAPI 3.0 + 链式 Builder：

```go
docs.Register("GET", "/api/users",
    docs.Op("获取用户列表", "用户管理").
        Res("200", "成功", docs.Arr(docs.Ref("User"), "用户列表")))
```

详细：[docs/go-gin.md](./docs/go-gin.md)

### Go 标准库（零依赖）

与 go-gin 完全相同的设计，业务代码只使用 `net/http`，编译后只依赖 Go 运行时。

详细：[docs/go-stdlib.md](./docs/go-stdlib.md)

> 📚 完整集成指南请参阅 [docs/](./docs/) 目录

---

## 示例项目

`examples/` 目录包含 **5 个**开箱即用的完整示例：

```bash
# 一键复制 Knife4j 前端产物到各示例（PowerShell）
foreach ($proj in @("fastapi", "litestar", "java-springboot")) {
    $dest = "$PWD\examples\$proj\"
    if ($proj -eq "java-springboot") { $dest = "$dest\src\main\resources\" }
    $dest = $dest + "static"
    if (Test-Path $dest) { Remove-Item $dest -Recurse -Force }
    New-Item -ItemType Directory -Path $dest -Force | Out-Null
    Copy-Item -Path "$PWD\dist\*" -Destination $dest -Recurse -Force
}
```

| 示例 | 启动命令 | 访问地址 |
|------|---------|---------|
| [Java Spring Boot](./examples/java-springboot) | `mvn spring-boot:run` | http://localhost:8080/api/doc.html |
| [FastAPI](./examples/fastapi) | `python main.py` | http://localhost:8000/doc.html |
| [LiteStar](./examples/litestar) | `uvicorn main:app --root-path /api` | http://localhost:8000/api/doc.html |
| [Go Gin](./examples/go-gin) | `go run main.go` | http://localhost:8080/doc.html |
| [Go 标准库](./examples/go-stdlib) | `go run main.go`（零依赖） | http://localhost:8080/doc.html |

> 每个示例项目都有独立的 README，详见 [examples/README.md](./examples/README.md)

---

## 项目结构

```
knife4j-vue3/
├── src/
│   ├── core/                        # 核心引擎
│   │   ├── Knife4jAsync.js          # OpenAPI/Swagger 规范解析引擎（7000+ 行）
│   │   ├── utils.js                 # 工具函数库
│   │   ├── json5.js                 # JSON5 解析器
│   │   └── oas3/                    # OpenAPI 3.0 专用解析器
│   │       ├── OAS3BaseModel.js     # OAS3 数据模型
│   │       └── OAS3ResponseExampleReader.js  # 响应示例读取器
│   │
│   ├── components/                  # UI 组件层
│   │   ├── GlobalHeader/            # 顶部导航栏
│   │   ├── GlobalHeaderTab/         # 多标签页管理
│   │   ├── SiderMenu/               # 侧边菜单（支持三级嵌套）
│   │   ├── HeaderSearch/            # 全局搜索组件
│   │   ├── GlobalFooter/            # 底部信息栏
│   │   ├── Markdown/                # Markdown 渲染器
│   │   └── common/                  # 通用组件
│   │       ├── ContextMenu.vue      # 右键菜单（关闭标签页）
│   │       └── MethodApi.vue        # HTTP 方法标识
│   │
│   ├── views/                       # 页面视图层
│   │   ├── api/                     # API 文档核心页面
│   │   │   ├── index.vue            # API 详情容器
│   │   │   ├── Document.vue         # 文档展示页
│   │   │   ├── Debug.vue            # 在线调试面板（2000+ 行）
│   │   │   ├── OpenApi.vue          # OpenAPI 原始规范展示
│   │   │   ├── EditorDebugShow.vue  # 调试响应编辑器
│   │   │   └── ScriptView.vue       # 脚本视图
│   │   ├── index/                   # 首页
│   │   └── settings/                # 设置管理
│   │       ├── GlobalParameters.vue # 全局参数配置（支持修改/删除）
│   │       ├── Authorize.vue        # 认证管理
│   │       └── Settings.vue         # 界面设置
│   │
│   ├── store/                       # 状态管理（Pinia）
│   │   ├── constants.js             # 常量定义
│   │   └── modules/                 # 模块化 Store
│   │       ├── global.js            # 全局状态
│   │       └── header.js            # 标签页状态
│   │
│   ├── lang/                        # 国际化
│   │   ├── zh.js                    # 中文
│   │   ├── en.js                    # 英文
│   │   └── jp.js                    # 日文
│   │
│   ├── layouts/                     # 布局引擎
│   │   ├── BasicLayout.vue          # 主布局（标签页 + 右键菜单）
│   │   └── menu.js                  # 菜单构建器
│   │
│   ├── router/                      # 路由配置
│   ├── assets/                      # 静态资源
│   └── style/                       # 全局样式
│
├── examples/                        # 示例项目
│   ├── java-springboot/             # Spring Boot 完整示例
│   ├── fastapi/                     # FastAPI 完整示例
│   ├── litestar/                    # LiteStar 完整示例
│   ├── go-gin/                      # Go Gin 完整示例（含 docs/openapi.go）
│   ├── go-stdlib/                   # Go 标准库示例（零依赖，含 docs/openapi.go）
│   └── README.md                    # 示例项目总览
│
├── docs/                            # 集成文档
│   ├── java-springboot.md           # Java 集成指南
│   ├── python-fastapi.md            # FastAPI 集成指南
│   ├── python-litestar.md           # LiteStar 集成指南
│   ├── go-gin.md                    # Go Gin 集成指南
│   ├── go-stdlib.md                 # Go 标准库集成指南（零依赖）
│   ├── openapi-loading.md           # OpenAPI 加载机制
│   └── other.md                     # 其他框架指南
│
├── doc.html                         # 构建入口 HTML
├── vite.config.js                   # Vite 构建配置
├── package.json                     # 项目依赖
└── README.md                        # 项目说明
```

---

## 技术架构

```
┌─────────────────────────────────────────────────────────┐
│                     浏览器                               │
│  ┌─────────────────────────────────────────────────────┐│
│  │              Knife4j Vue3 前端                       ││
│  │  ┌─────────┐  ┌──────────┐  ┌──────────────────┐  ││
│  │  │ UI 层   │  │ 状态管理  │  │   核心解析引擎     │  ││
│  │  │(Ant     │  │ (Pinia)  │  │  (Knife4jAsync)  │  ││
│  │  │ Design) │  │          │  │                  │  ││
│  │  └────┬────┘  └────┬─────┘  └────────┬─────────┘  ││
│  │       │            │                 │             ││
│  └───────┼────────────┼─────────────────┼─────────────┘│
│          │            │                 │               │
│          ▼            ▼                 ▼               │
│  ┌─────────────────────────────────────────────────────┐│
│  │              HTTP 请求                               ││
│  │  GET {prefix}/v3/api-docs/swagger-config            ││
│  │  GET {prefix}/openapi.json（或 /v3/api-docs 等）      ││
│  │  GET {prefix}/api/*（业务 API）                       ││
│  └─────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   后端服务                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │ Spring   │  │ FastAPI  │  │ LiteStar │  │  Go    │ │
│  │ Boot     │  │ (Python) │  │ (Python) │  │ (Gin)  │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘ │
│       │              │              │            │      │
│       └──────────────┴──────────────┴────────────┘      │
│                        │                                 │
│              OpenAPI 3.0 规范输出                         │
└─────────────────────────────────────────────────────────┘
```

---

## 环境变量

### .env 文件说明

Vite 使用 `.env` 文件来管理环境变量，不同文件在不同场景下生效：

| 文件名 | 生效时机 | 用途 | 是否需要保留 |
|--------|----------|------|-------------|
| `.env` | 所有模式 | 通用配置 | ✅ 可选 |
| `.env.development` | `npm run dev` 开发模式 | 开发专用配置 | ✅ **必须保留** |
| `.env.production` | `npm run build` 构建时 | 生产构建配置 | ✅ 可选 |

#### `.env.development` 文件的价值

`.env.development` 文件在**本地开发模式**下至关重要：

1. **API 代理路径**：`VITE_APP_BASE_API` 定义了开发时的 API 请求基础路径，配合 `vite.config.js` 中的代理配置，实现前后端联调
2. **框架特性开关**：`VITE_RELEASE_APP_TYPE` 和 `VITE_FRAMEWORK` 用于启用特定框架的功能特性
3. **调试模式**：`VITE_DEBUG = true` 启用详细的调试日志输出

```
开发流程：
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│  浏览器请求   │ ──▶  │  Vite 代理   │ ──▶  │  后端服务    │
│ /api/users  │      │ (dev server)│      │ :8000       │
└─────────────┘      └─────────────┘      └─────────────┘
                            │
                     读取 .env.development
                     VITE_APP_BASE_API=/api
```

> ⚠️ **重要提示**：运行 `pnpm build` 构建生产版本时，环境变量会被**编译打包**到 JS 文件中。构建后的 `dist/` 目录不再依赖任何 `.env` 文件。

### 环境变量列表

| 变量名 | 说明 | 可选值 | 默认值 |
|--------|------|--------|--------|
| `VITE_APP_BASE_API` | API 基础路径 | — | `/api` |
| `VITE_RELEASE_APP_TYPE` | 发行版本类型 | 见下表 | `SpringDocOpenApi` |
| `VITE_FRAMEWORK` | 目标框架标识 | `LiteStar` 等 | — |
| `VITE_FRAMEWORK_VERSION` | 框架版本号 | — | — |
| `VITE_DEBUG` | 调试模式 | `true` / `false` | `false` |

### 发行版本类型

| 值 | 适用场景 |
|----|---------|
| `SpringDocOpenApi` | Spring Boot + SpringDoc OpenAPI 3.0（默认） |
| `Knife4jSpringUi` | Spring Boot + Knife4j Spring UI |
| `Knife4jJFinal` | JFinal 框架 |
| `Knife4jFront` | 纯前端模式 |
| `LiteStarOpenApi` | Python LiteStar ASGI 框架 |

---

## 故障排查（FAQ）

### 1. 调试面板响应区空白？

最常见原因：**重复路径前缀**。

```bash
# 期望
GET /api/users        → 200

# 实际（修复前可能出现的）
GET /api/api/users    → 404
```

Knife4j 前端已在 2026-08 修复此 bug，确保使用最新版本（`pnpm build` 后再 `cp -r dist/*`）。

### 2. "No API definitions found"？

检查后端 swagger-config 端点：

```bash
curl http://localhost:8080/api/v3/api-docs/swagger-config
```

应返回：

```json
{
  "urls": [{"url": "/api/v3/api-docs", "name": "default"}],
  "configUrl": "/api/v3/api-docs/swagger-config"
}
```

### 3. 全局参数添加后实际请求未携带？

进入「文档管理」→「全局参数设置」，**确认新添加的参数左侧复选框是勾选状态**。未勾选的参数不会在请求中发送。

### 4. 前端资源 404？

每个示例的 `static/` 目录都不在版本控制里，需要先编译前端并复制：

```bash
pnpm build
cp -r dist/* examples/<project>/static/    # 大多数项目
cp -r dist/* examples/java-springboot/src/main/resources/static/   # Java 项目
```

### 5. 端口冲突？

各示例默认端口：

| 示例 | 默认 |
|------|------|
| Java Spring Boot | 8080 |
| FastAPI | 8000 |
| LiteStar | 8000 |
| Go Gin | 8080 |
| Go 标准库 | 8080 |

启动两个都用 8080 的项目时，先启动的项目用 `mvn spring-boot:run -Dspring-boot.run.arguments=--server.port=8090` 或修改 `r.Run(":9090")`。

---

## 与同类工具对比

| 能力 | Knife4j Vue3 | Swagger UI | Redoc | Apifox |
|------|-------------|------------|-------|--------|
| **开源** | ✅ Apache 2.0 | ✅ Apache 2.0 | ✅ MIT | ❌ 商业 |
| **自托管** | ✅ | ✅ | ✅ | ❌ |
| **在线调试** | ✅ 增强版 | ✅ 基础版 | ❌ | ✅ |
| **cURL 生成** | ✅ | ❌ | ❌ | ✅ |
| **全局搜索** | ✅ | ❌ | ❌ | ✅ |
| **参数缓存** | ✅ | ❌ | ❌ | ✅ |
| **多标签页** | ✅ | ❌ | ❌ | ✅ |
| **Markdown 文档** | ✅ | ❌ | ✅ | ✅ |
| **多语言** | ✅ 中/英/日 | 部分 | ✅ | ✅ |
| **跨语言后端** | ✅ 通用 | ✅ 通用 | ✅ 通用 | ✅ |
| **国内生态适配** | ✅ 完美 | 一般 | 一般 | ✅ |
| **零依赖部署** | ✅（Go 标准库示例） | ❌ | ❌ | ❌ |

---

## 构建与部署

### 编译

```bash
pnpm install
pnpm build
```

产物输出到 `dist/` 目录，包含：

- `doc.html` — Knife4j 入口页面
- `webjars/` — JS/CSS 静态资源
- `oauth/` — OAuth2 授权页面

### 部署方案

| 方案 | 适用场景 | 说明 |
|------|---------|------|
| **嵌入后端** | 单体应用 | 将 `dist/` 复制到后端的 static 目录 |
| **Nginx 代理** | 微服务架构 | Nginx 托管前端，反向代理后端 API |
| **Docker** | 容器化部署 | 多阶段构建，前端编译 + 后端运行一体化 |
| **单二进制**（仅 Go） | CLI、嵌入式工具 | Go 用 `go:embed` 把前端打入二进制 |

> 📚 详细部署指南请参阅各 [后端集成文档](./docs/)

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'feat: add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 提交 Pull Request

---

## 开源协议

[Apache License 2.0](LICENSE)

基于 [Knife4j](https://github.com/xiaoymin/knife4j) 开源项目扩展，感谢原作者 xiaoymin 的贡献。
