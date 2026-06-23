# 其他 Web 框架对接指南

本文档介绍如何将 Knife4j Vue3 前端与**任意 Web 框架**集成。不针对特定框架，而是阐述核心原理和通用方法，适用于 Flask、Django、Tornado、Express、Koa、Gin、Echo、Spring MVC、Vert.x、Rust Actix、Ruby Sinatra 等任何框架。

## 核心原理

Knife4j Vue3 是一个 OpenAPI 文档 UI，其工作原理非常简单：

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  Knife4j Vue3   │      │  swagger-config │      │  OpenAPI JSON   │
│  (前端 UI)      │─────>│  端点           │─────>│  规范文档       │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

**后端只需要提供两个 HTTP 端点：**

1. **`/v3/api-docs/swagger-config`** - 配置端点，告诉前端 OpenAPI 文档在哪里
2. **`/openapi.json`** - OpenAPI 3.0 规范文档

其余都是前端的工作。

## 必要条件

### 1. 两个 HTTP 端点

无论使用什么框架，只需要实现这两个端点：

#### 端点 1：swagger-config

```
GET /v3/api-docs/swagger-config
```

响应格式：

```json
{
  "urls": [
    {
      "url": "/openapi.json",
      "name": "default"
    }
  ],
  "configUrl": "/v3/api-docs/swagger-config",
  "validatorUrl": ""
}
```

#### 端点 2：OpenAPI 文档

```
GET /openapi.json
```

响应格式：标准的 OpenAPI 3.0 JSON 规范（见下方示例）。

### 2. CORS 跨域支持

前后端分离部署时，后端必须支持 CORS，允许前端域名访问。

### 3. JSON 响应

两个端点必须返回 `Content-Type: application/json`。

## OpenAPI 3.0 规范文档模板

以下是一个最小可用的 OpenAPI 3.0 文档模板，你可以基于此扩展：

```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "My API",
    "description": "API Documentation with Knife4j Vue3",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "/",
      "description": "Local server"
    }
  ],
  "paths": {
    "/api/users": {
      "get": {
        "summary": "获取所有用户",
        "description": "返回用户列表",
        "operationId": "listUsers",
        "tags": ["用户管理"],
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/User"
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "创建新用户",
        "description": "创建一个新用户",
        "operationId": "createUser",
        "tags": ["用户管理"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UserInput"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "创建成功",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          }
        }
      }
    },
    "/api/users/{userId}": {
      "get": {
        "summary": "根据ID获取用户",
        "operationId": "getUserById",
        "tags": ["用户管理"],
        "parameters": [
          {
            "name": "userId",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            },
            "description": "用户ID"
          }
        ],
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          },
          "404": {
            "description": "用户不存在"
          }
        }
      },
      "put": {
        "summary": "更新用户信息",
        "operationId": "updateUser",
        "tags": ["用户管理"],
        "parameters": [
          {
            "name": "userId",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UserInput"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "更新成功",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          },
          "404": {
            "description": "用户不存在"
          }
        }
      },
      "delete": {
        "summary": "删除用户",
        "operationId": "deleteUser",
        "tags": ["用户管理"],
        "parameters": [
          {
            "name": "userId",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "删除成功"
          },
          "404": {
            "description": "用户不存在"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer",
            "description": "用户ID"
          },
          "name": {
            "type": "string",
            "description": "用户名"
          },
          "email": {
            "type": "string",
            "format": "email",
            "description": "邮箱地址"
          }
        },
        "required": ["name", "email"]
      },
      "UserInput": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "description": "用户名"
          },
          "email": {
            "type": "string",
            "format": "email",
            "description": "邮箱地址"
          }
        },
        "required": ["name", "email"]
      }
    }
  }
}
```

## 各框架对接示例

### Python Flask

```python
from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

OPENAPI_SCHEMA = { ... }  # 上方的 OpenAPI 文档模板

@app.route("/openapi.json")
def openapi_schema():
    return jsonify(OPENAPI_SCHEMA)

@app.route("/v3/api-docs/swagger-config")
def swagger_config():
    return jsonify({
        "urls": [{"url": "/openapi.json", "name": "default"}],
        "configUrl": "/v3/api-docs/swagger-config",
        "validatorUrl": ""
    })

app.run(port=5000)
```

### Python Django

```python
# views.py
from django.http import JsonResponse

OPENAPI_SCHEMA = { ... }  # 上方的 OpenAPI 文档模板

def openapi_schema(request):
    return JsonResponse(OPENAPI_SCHEMA)

def swagger_config(request):
    return JsonResponse({
        "urls": [{"url": "/openapi.json", "name": "default"}],
        "configUrl": "/v3/api-docs/swagger-config",
        "validatorUrl": ""
    })
```

```python
# urls.py
from django.urls import path
from . import views

urlpatterns = [
    path('openapi.json', views.openapi_schema),
    path('v3/api-docs/swagger-config', views.swagger_config),
]
```

Django 还需要配置 CORS，推荐使用 `django-cors-headers`：

```bash
pip install django-cors-headers
```

```python
# settings.py
INSTALLED_APPS = [
    ...
    'corsheaders',
]

MIDDLEWARE = [
    'corsheaders.middleware.CorsMiddleware',
    ...
]

CORS_ALLOW_ALL_ORIGINS = True
```

### Python Tornado

```python
import tornado.web
import tornado.ioloop

OPENAPI_SCHEMA = { ... }  # 上方的 OpenAPI 文档模板

class BaseHandler(tornado.web.RequestHandler):
    def set_default_headers(self):
        self.set_header("Access-Control-Allow-Origin", "*")
        self.set_header("Access-Control-Allow-Methods", "GET, OPTIONS")
        self.set_header("Access-Control-Allow-Headers", "Content-Type")

    def options(self):
        self.set_status(204)
        self.finish()

class OpenAPIHandler(BaseHandler):
    def get(self):
        self.set_header("Content-Type", "application/json")
        self.write(OPENAPI_SCHEMA)

class SwaggerConfigHandler(BaseHandler):
    def get(self):
        self.set_header("Content-Type", "application/json")
        self.write({
            "urls": [{"url": "/openapi.json", "name": "default"}],
            "configUrl": "/v3/api-docs/swagger-config",
            "validatorUrl": ""
        })

app = tornado.web.Application([
    (r"/openapi.json", OpenAPIHandler),
    (r"/v3/api-docs/swagger-config", SwaggerConfigHandler),
])

app.listen(5000)
tornado.ioloop.IOLoop.current().start()
```

### Node.js Express

```javascript
const express = require('express');
const cors = require('cors');

const app = express();
app.use(cors());

const OPENAPI_SCHEMA = { ... }; // 上方的 OpenAPI 文档模板

app.get('/openapi.json', (req, res) => {
  res.json(OPENAPI_SCHEMA);
});

app.get('/v3/api-docs/swagger-config', (req, res) => {
  res.json({
    urls: [{ url: '/openapi.json', name: 'default' }],
    configUrl: '/v3/api-docs/swagger-config',
    validatorUrl: ''
  });
});

app.listen(5000);
```

### Node.js Koa

```javascript
const Koa = require('koa');
const Router = require('@koa/router');
const cors = require('@koa/cors');

const app = new Koa();
const router = new Router();

app.use(cors());

const OPENAPI_SCHEMA = { ... }; // 上方的 OpenAPI 文档模板

router.get('/openapi.json', (ctx) => {
  ctx.body = OPENAPI_SCHEMA;
});

router.get('/v3/api-docs/swagger-config', (ctx) => {
  ctx.body = {
    urls: [{ url: '/openapi.json', name: 'default' }],
    configUrl: '/v3/api-docs/swagger-config',
    validatorUrl: ''
  };
});

app.use(router.routes());
app.listen(5000);
```

### Go Gin

```go
package main

import (
    "net/http"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.Use(cors.Default())

    r.GET("/openapi.json", func(c *gin.Context) {
        c.JSON(http.StatusOK, openapiSchema) // openapiSchema 为 map 类型
    })

    r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "urls":       []gin.H{{"url": "/openapi.json", "name": "default"}},
            "configUrl":  "/v3/api-docs/swagger-config",
            "validatorUrl": "",
        })
    })

    r.Run(":5000")
}
```

### Go Echo

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()
    e.Use(middleware.CORS())

    e.GET("/openapi.json", func(c echo.Context) error {
        return c.JSON(http.StatusOK, openapiSchema)
    })

    e.GET("/v3/api-docs/swagger-config", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]interface{}{
            "urls":         []map[string]string{{"url": "/openapi.json", "name": "default"}},
            "configUrl":    "/v3/api-docs/swagger-config",
            "validatorUrl": "",
        })
    })

    e.Start(":5000")
}
```

### Java Servlet

```java
@WebServlet(urlPatterns = {"/openapi.json", "/v3/api-docs/swagger-config"})
public class OpenAPIServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp)
            throws IOException {
        resp.setContentType("application/json");
        resp.setHeader("Access-Control-Allow-Origin", "*");

        if (req.getRequestURI().endsWith("swagger-config")) {
            resp.getWriter().write(
                "{\"urls\":[{\"url\":\"/openapi.json\",\"name\":\"default\"}]," +
                "\"configUrl\":\"/v3/api-docs/swagger-config\"," +
                "\"validatorUrl\":\"\"}"
            );
        } else {
            resp.getWriter().write(getOpenAPISchema());
        }
    }
}
```

### Rust Actix-web

```rust
use actix_web::{web, App, HttpServer, HttpResponse, middleware};
use serde_json::json;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .wrap(middleware::Cors::permissive())
            .route("/openapi.json", web::get().to(|| async {
                HttpResponse::Ok().json(json!({ "openapi": "3.0.0", ... }))
            }))
            .route("/v3/api-docs/swagger-config", web::get().to(|| async {
                HttpResponse::Ok().json(json!({
                    "urls": [{"url": "/openapi.json", "name": "default"}],
                    "configUrl": "/v3/api-docs/swagger-config",
                    "validatorUrl": ""
                }))
            }))
    })
    .bind("0.0.0.0:5000")?
    .run()
    .await
}
```

### Ruby Sinatra

```ruby
require 'sinatra'
require 'json'

set :port, 5000

before do
  headers 'Access-Control-Allow-Origin' => '*',
          'Access-Control-Allow-Methods' => 'GET'
end

get '/openapi.json' do
  content_type :json
  # openapi_schema 为 Hash
  openapi_schema.to_json
end

get '/v3/api-docs/swagger-config' do
  content_type :json
  {
    urls: [{ url: '/openapi.json', name: 'default' }],
    configUrl: '/v3/api-docs/swagger-config',
    validatorUrl: ''
  }.to_json
end
```

### PHP 原生

```php
<?php
// openapi.json
header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');

$schema = [
    'openapi' => '3.0.0',
    'info' => ['title' => 'My API', 'version' => '1.0.0'],
    // ...
];

echo json_encode($schema);
```

```php
<?php
// swagger-config.php
header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');

echo json_encode([
    'urls' => [['url' => '/openapi.json', 'name' => 'default']],
    'configUrl' => '/v3/api-docs/swagger-config',
    'validatorUrl' => ''
]);
```

## 前端配置

对接完成后，需要在 Knife4j Vue3 前端进行配置。

### 设置全局变量

在 `index.html` 中添加：

```html
<script>
  window.KNIFE4J_FRAMEWORK = 'Generic';
  window.KNIFE4J_FRAMEWORK_VERSION = '1.0.0';
</script>
```

### 开发环境代理配置

在 `vite.config.js` 中配置代理：

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/v3/api-docs': {
        target: 'http://localhost:5000',
        changeOrigin: true
      },
      '/openapi.json': {
        target: 'http://localhost:5000',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:5000',
        changeOrigin: true
      }
    }
  }
})
```

## 部署方式

### 方式一：前后端分离（推荐）

使用 nginx 反向代理：

```nginx
server {
    listen 80;
    server_name example.com;

    # 前端静态资源
    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 和文档端点
    location /api {
        proxy_pass http://localhost:5000;
    }

    location /openapi.json {
        proxy_pass http://localhost:5000;
    }

    location /v3/api-docs {
        proxy_pass http://localhost:5000;
    }
}
```

### 方式二：后端托管前端静态资源

将构建产物复制到后端的静态资源目录，由后端同时托管前端页面和 API。

```bash
npm run build
cp -r dist/* /path/to/backend/static/
```

## 常见问题

### Q1: 无法加载 OpenAPI 文档

1. 确保 `/openapi.json` 端点返回正确的 JSON 格式
2. 确保 CORS 配置正确
3. 检查浏览器控制台是否有跨域错误

### Q2: swagger-config 端点 404

确保 `/v3/api-docs/swagger-config` 端点存在且返回正确格式：

```json
{
  "urls": [{"url": "/openapi.json", "name": "default"}],
  "configUrl": "/v3/api-docs/swagger-config",
  "validatorUrl": ""
}
```

### Q3: 调试功能不工作

确保 API 路由路径与 OpenAPI 文档中定义的路径一致，包括请求方法、路径参数、查询参数等。

### Q4: 如何自定义 OpenAPI 文档路径

在 `urls` 数组中可以指定任意路径：

```json
{
  "urls": [
    {"url": "/api/v1/openapi.json", "name": "V1"},
    {"url": "/api/v2/openapi.json", "name": "V2"}
  ]
}
```

Knife4j 会加载所有 URL 并展示在界面上供切换。

## 快速验证清单

对接完成后，按以下步骤验证：

1. 浏览器访问 `http://localhost:5000/openapi.json`，确认返回合法 JSON
2. 浏览器访问 `http://localhost:5000/v3/api-docs/swagger-config`，确认返回合法 JSON
3. 确认两个端点的响应头包含 `Content-Type: application/json`
4. 确认 CORS 头正确（跨域场景下）
5. 启动 Knife4j Vue3 前端，确认文档正常加载
