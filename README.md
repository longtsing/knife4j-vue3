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
- **全局参数管理**：统一配置 Authorization、自定义 Header 等公共参数
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

| 框架 | 语言 | OpenAPI 版本 | 集成难度 | 状态 |
|------|------|-------------|---------|------|
| **Spring Boot** (SpringDoc / Springfox) | Java | 2.0 / 3.0 | ⭐ 极简 | ✅ 完全支持 |
| **FastAPI** | Python | 3.0 | ⭐⭐ 简单 | ✅ 完全支持 |
| **LiteStar** | Python | 3.0 | ⭐⭐ 简单 | ✅ 完全支持 |
| **JFinal** | Java | 2.0 / 3.0 | ⭐ 极简 | ✅ 完全支持 |
| **通用 OpenAPI 3.0** | 任意语言 | 3.0 | ⭐⭐⭐ 中等 | ✅ 完全支持 |

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

### 方式一：直接使用（推荐）

将编译产物部署到你的后端项目中，即可拥有专业级 API 文档界面。

```bash
# 1. 获取项目
git clone <repo-url>
cd knife4j-vue3

# 2. 安装依赖
pnpm install

# 3. 编译生产版本
pnpm build

# 4. 将 dist/ 目录部署到后端项目
```

### 方式二：开发模式

```bash
# 启动开发服务器（热更新）
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

### Java Spring Boot（3 行配置）

```xml
<dependency>
    <groupId>com.github.xiaoymin</groupId>
    <artifactId>knife4j-openapi3-jakarta-spring-boot-starter</artifactId>
    <version>4.5.0</version>
</dependency>
```

```yaml
springdoc:
  swagger-ui:
    path: /doc.html
```

### Python FastAPI（5 行配置）

```python
from fastapi import FastAPI
from fastapi.responses import FileResponse

app = FastAPI(docs_url=None)  # 禁用默认 Swagger UI

@app.get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    return FileResponse("static/doc.html")
```

### Python LiteStar（5 行配置）

```python
from litestar import Litestar, get
from litestar.response import Response

@get("/doc.html", include_in_schema=False)
async def knife4j_ui():
    with open("static/doc.html", "r") as f:
        return Response(content=f.read(), media_type="text/html")
```

> 📚 完整集成指南请参阅 [docs/](./docs/) 目录

---

## 示例项目

`examples/` 目录包含三个**开箱即用**的完整示例：

```bash
# 一键配置（Windows）
cd examples && setup.bat

# 或手动配置
cp -r dist/* examples/fastapi/static/
cd examples/fastapi && uvicorn main:app --port 8000
```

| 示例 | 启动命令 | 访问地址 |
|------|---------|---------|
| [Java Spring Boot](./examples/java-springboot) | `mvn spring-boot:run` | http://localhost:8080/doc.html |
| [FastAPI](./examples/fastapi) | `uvicorn main:app --port 8000` | http://localhost:8000/doc.html |
| [LiteStar](./examples/litestar) | `uvicorn main:app --port 8000` | http://localhost:8000/doc.html |

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
│   │       ├── GlobalParameters.vue # 全局参数配置
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
│   └── setup.bat                    # 一键配置脚本
│
├── docs/                            # 集成文档
│   ├── java-springboot.md           # Java 集成指南
│   ├── python-fastapi.md            # FastAPI 集成指南
│   └── python-litestar.md           # LiteStar 集成指南
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
│  │  │ UI 层   │  │ 状态管理  │  │   核心解析引擎    │  ││
│  │  │(Ant     │  │ (Pinia)  │  │  (Knife4jAsync)  │  ││
│  │  │ Design) │  │          │  │                  │  ││
│  │  └────┬────┘  └────┬─────┘  └────────┬─────────┘  ││
│  │       │            │                 │             ││
│  └───────┼────────────┼─────────────────┼─────────────┘│
│          │            │                 │               │
│          ▼            ▼                 ▼               │
│  ┌─────────────────────────────────────────────────────┐│
│  │              HTTP 请求                               ││
│  │  GET /v3/api-docs/swagger-config  → OpenAPI 规范     ││
│  │  GET /api/*                       → 业务 API         ││
│  └─────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   后端服务                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │ Spring   │  │ FastAPI  │  │ LiteStar │  │  其他   │ │
│  │ Boot     │  │ (Python) │  │ (Python) │  │ 框架   │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘ │
│       │              │              │            │      │
│       └──────────────┴──────────────┴────────────┘      │
│                        │                                 │
│              OpenAPI 3.0 规范输出                         │
└─────────────────────────────────────────────────────────┘
```

---

## 环境变量

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

### 部署方案

| 方案 | 适用场景 | 说明 |
|------|---------|------|
| **嵌入后端** | 单体应用 | 将 `dist/` 复制到后端的 static 目录 |
| **Nginx 代理** | 微服务架构 | Nginx 托管前端，反向代理后端 API |
| **Docker** | 容器化部署 | 多阶段构建，前端编译 + 后端运行一体化 |

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
