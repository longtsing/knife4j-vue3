# OpenAPI 加载逻辑详解

本文档详细介绍 Knife4j Vue3 如何加载 OpenAPI 规范文档，包括端点配置、URL 构建逻辑和框架适配机制。

## 概述

Knife4j Vue3 支持多种后端框架（FastAPI、LiteStar、Spring Boot、Go Gin 等），通过统一的 OpenAPI 3.0 端点加载 API 文档。

## 核心加载流程

```
┌─────────────────────────────────────────────────────────────────┐
│                      Knife4j Vue3 初始化                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. BasicLayout.vue 调用 initSwagger()                           │
│     - 读取框架配置 (KNIFE4J_FRAMEWORK)                            │
│     - 设置 springdoc 模式                                         │
│     - 配置 URL: /v3/api-docs/swagger-config                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Knife4jAsync.js SwaggerBootstrapUi 构造函数                   │
│     - 设置 this.url = '/v3/api-docs/swagger-config'              │
│     - 设置 this.fallbackUrl = '/api/v3/api-docs/swagger-config'  │
│     - 设置 this.configUrl (OpenAPI 3.0 风格)                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. 请求 swagger-config 端点                                      │
│     GET /v3/api-docs/swagger-config                               │
│     ↓ 失败时重试备用 URL                                           │
│     GET /api/v3/api-docs/swagger-config                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. 解析 swagger-config 响应                                       │
│     {                                                             │
│       "urls": [{"url": "/openapi.json", "name": "default"}],     │
│       "configUrl": "/v3/api-docs/swagger-config",                │
│       "validatorUrl": ""                                          │
│     }                                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  5. 请求 OpenAPI 规范文档                                         │
│     GET /openapi.json (或 urls 中指定的路径)                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  6. 解析 OpenAPI 文档并渲染 UI                                     │
│     - 解析 paths、components、servers 等                          │
│     - 构建 API 分组菜单                                            │
│     - 初始化调试功能                                               │
└─────────────────────────────────────────────────────────────────┘
```

## 端点配置

### 框架配置文件 (`src/config/framework.js`)

```javascript
export const FRAMEWORK_CONFIG = {
  // API 端点配置（OpenAPI 3.0 风格）
  endpoints: {
    swaggerConfig: '/v3/api-docs/swagger-config',
    swaggerConfigFallback: '/api/v3/api-docs/swagger-config',
    openApiSchema: '/openapi.json',
    apiDocs: '/v3/api-docs'
  },
  
  framework: {
    name: 'Generic',  // 可选: FastAPI, LiteStar, SpringFox, Generic
    version: '1.0.0'
  }
};
```

### 端点说明

| 端点 | 用途 | 响应格式 |
|------|------|----------|
| `/v3/api-docs/swagger-config` | Knife4j 配置端点 | `{urls: [...], configUrl: "..."}` |
| `/openapi.json` | OpenAPI 3.0 规范文档 | OpenAPI 3.x JSON |
| `/v3/api-docs` | OpenAPI 规范（可选） | OpenAPI 3.x JSON |

## URL 构建逻辑

### 初始化配置 (`src/layouts/BasicLayout.vue`)

```javascript
// 通用框架配置（支持 FastAPI、LiteStar、Spring Boot 等）
const framework = window.KNIFE4J_FRAMEWORK || 'Generic';
const isSpringFox = framework === 'SpringFox';

initSwagger({
  springdoc: !isSpringFox,           // 非 SpringFox 框架启用 SpringDoc 模式
  baseSpringFox: isSpringFox,        // SpringFox 使用 basePath 处理
  url: '/v3/api-docs/swagger-config', // 统一使用 OpenAPI 3.0 端点
  framework: framework,
  frameworkVersion: window.KNIFE4J_FRAMEWORK_VERSION || '1.0.0',
  disableBasePath: !isSpringFox
});
```

### 核心逻辑 (`src/core/Knife4jAsync.js`)

```javascript
function SwaggerBootstrapUi(options) {
  // 默认启用 springdoc 模式以支持 OpenAPI 3.0
  this.springdoc = options.springdoc || true;
  
  if (this.springdoc) {
    // OpenAPI 3.0 框架（FastAPI、LiteStar 等）
    const path = window.location.pathname;
    const index = path.lastIndexOf('/');
    const basePath = path.length == index + 1 ? path : path.substring(0, index);
    
    // 主 URL
    this.url = options.url || '/v3/api-docs/swagger-config';
    // 备用 URL（FastAPI 常用 /api 前缀）
    this.fallbackUrl = basePath + '/api/v3/api-docs/swagger-config';
  } else {
    // SpringFox 框架
    this.url = options.url || 'swagger-resources';
  }
  
  // configUrl 也根据模式选择
  this.configUrl = options.configUrl || 
    (this.springdoc ? '/v3/api-docs/swagger-config' : 'swagger-resources/configuration/ui');
}
```

## 重试机制

当主 URL 请求失败时，自动尝试备用 URL：

```javascript
ajax.request(requestConfig).then(data => {
  success(data);
}).catch(err => {
  // OpenAPI 3.0 框架重试逻辑
  if (this.springdoc && !isRetry && this.fallbackUrl && 
      config.url && config.url.includes('swagger-config')) {
    console.warn(this.framework + ' Primary URL failed, trying fallback:', this.fallbackUrl);
    performRequest(this.fallbackUrl, true);
  } else {
    error(err);
  }
});
```

## 框架适配

### 框架类型检测

```javascript
// 通过全局变量设置框架类型
window.KNIFE4J_FRAMEWORK = 'FastAPI';     // 或 'LiteStar', 'Generic', 'SpringFox'
window.KNIFE4J_FRAMEWORK_VERSION = '0.100.0';
```

### 不同框架的端点差异

| 框架 | swagger-config 端点 | OpenAPI 文档端点 |
|------|---------------------|------------------|
| FastAPI | `/v3/api-docs/swagger-config` (需手动创建) | `/openapi.json` |
| LiteStar | `/schema/swagger-config` | `/schema/openapi.json` |
| Spring Boot (SpringDoc) | `/v3/api-docs/swagger-config` | `/v3/api-docs` |
| Spring Boot (SpringFox) | `/swagger-resources` | `/v2/api-docs` |
| Go Gin | `/v3/api-docs/swagger-config` (需手动创建) | `/openapi.json` |

### basePath 处理

OpenAPI 3.0 使用 `servers` 字段而非 `basePath`，Knife4j 会自动处理：

```javascript
// OpenAPI 3.0 框架特殊处理：如果检测到 /api/doc.html，不添加 /api 前缀（避免重复）
if (tempPath === '/api' && that.framework !== 'SpringFox') {
  console.log(that.framework + ' detected: Skipping /api basePath to avoid duplication');
  tempPath = '';
}
```

## FastAPI 后端配置示例

FastAPI 需要手动创建 `swagger-config` 端点：

```python
from fastapi import FastAPI
from fastapi.openapi.utils import get_openapi

app = FastAPI(title="My API", version="1.0.0")

@app.get("/v3/api-docs/swagger-config")
def swagger_config():
    """Knife4j 需要的 swagger-config 端点"""
    return {
        "urls": [
            {
                "url": "/openapi.json",
                "name": "default"
            }
        ],
        "configUrl": "/v3/api-docs/swagger-config",
        "validatorUrl": ""
    }

@app.get("/openapi.json")
def openapi_schema():
    """返回 OpenAPI 规范"""
    return get_openapi(
        title=app.title,
        version=app.version,
        description=app.description,
        routes=app.routes,
    )
```

## 调试技巧

### 1. 检查端点响应

```bash
# 检查 swagger-config 端点
curl http://localhost:8000/v3/api-docs/swagger-config

# 检查 OpenAPI 文档
curl http://localhost:8000/openapi.json
```

### 2. 浏览器控制台日志

Knife4j 会在控制台输出加载日志：

```
OpenAPI 3.0 Mode - Primary URL: /v3/api-docs/swagger-config
OpenAPI 3.0 Mode - Fallback URL: /api/v3/api-docs/swagger-config
FastAPI API Response: /v3/api-docs/swagger-config {...}
```

### 3. 常见问题排查

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 请求 `/swagger-resources/configuration/ui` | 未启用 springdoc 模式 | 检查 `initSwagger({ springdoc: true })` |
| API 路径缺少 `/api` 前缀 | servers 配置未正确解析 | 检查 OpenAPI 的 `servers` 字段 |
| "Cannot read properties of undefined" | OpenAPI 文档格式问题 | 检查 `components` 是否存在 |

## 架构图

```
┌──────────────────────────────────────────────────────────────────────┐
│                           Knife4j Vue3 前端                          │
├──────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐  │
│  │  BasicLayout    │  │  Knife4jAsync   │  │  Debug.vue          │  │
│  │  - initSwagger  │→ │  - URL 构建     │→ │  - baseURL 解析     │  │
│  │  - 框架配置     │  │  - 重试机制     │  │  - servers 处理     │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ HTTP 请求
┌──────────────────────────────────────────────────────────────────────┐
│                           后端 API 服务                               │
├──────────────────────────────────────────────────────────────────────┤
│  FastAPI / LiteStar / Spring Boot / Go Gin                           │
│                                                                      │
│  端点:                                                               │
│  - GET /v3/api-docs/swagger-config  → 返回 Knife4j 配置              │
│  - GET /openapi.json               → 返回 OpenAPI 3.0 规范           │
│  - GET /api/*                      → 业务 API                       │
└──────────────────────────────────────────────────────────────────────┘
```

## 相关文件

- [src/config/framework.js](../src/config/framework.js) - 框架配置
- [src/layouts/BasicLayout.vue](../src/layouts/BasicLayout.vue) - 初始化逻辑
- [src/core/Knife4jAsync.js](../src/core/Knife4jAsync.js) - 核心加载逻辑
- [src/views/api/Debug.vue](../src/views/api/Debug.vue) - 调试请求 URL 构建