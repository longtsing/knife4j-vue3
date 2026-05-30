# Go 标准库对接指南

> **零依赖方案**：仅使用 Go 标准库 `net/http` + `encoding/json`，无需任何第三方包。

## 前提条件

- Go 1.18+
- Knife4j Vue3 前端编译产物

## 核心原理

Knife4j Vue3 前端需要三个端点即可工作：

```
GET /doc.html               → 返回 Knife4j 前端页面
GET /v3/api-docs/swagger-config → 返回 OpenAPI JSON 的地址
GET /openapi.json            → 返回 OpenAPI 3.0 完整规范
```

前端加载流程：

```
浏览器访问 /doc.html
  → 请求 /v3/api-docs/swagger-config
  → 获取 OpenAPI JSON 地址（/openapi.json）
  → 请求 /openapi.json 获取完整规范
  → 渲染 API 文档界面
```

---

## 最小示例

以下是一个**最简化的完整示例**，可以直接复制使用：

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", router)
	fmt.Println("🚀 服务启动: http://localhost:8080/doc.html")
	http.ListenAndServe(":8080", nil)
}

func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// CORS 预检
	if r.Method == "OPTIONS" {
		setCORS(w)
		w.WriteHeader(204)
		return
	}

	switch {
	case path == "/v3/api-docs/swagger-config":
		// Knife4j 需要的配置端点
		jsonResponse(w, map[string]interface{}{
			"urls":          []map[string]string{{"url": "/openapi.json", "name": "default"}},
			"configUrl":     "/v3/api-docs/swagger-config",
			"validatorUrl":  "",
		})

	case path == "/openapi.json":
		// OpenAPI 3.0 规范
		jsonResponse(w, openAPISpec())

	case path == "/doc.html":
		// Knife4j 前端页面
		serveFile(w, "doc.html", "text/html")

	case strings.HasPrefix(path, "/webjars/"):
		// 静态资源
		serveFile(w, path, guessMIME(path))

	case path == "/api/hello":
		// 你的 API
		jsonResponse(w, map[string]string{"message": "Hello, Knife4j!"})

	default:
		http.NotFound(w, r)
	}
}

// jsonResponse 写入 JSON 响应
func jsonResponse(w http.ResponseWriter, data interface{}) {
	setCORS(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// setCORS 设置 CORS 头
func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// serveFile 提供静态文件
func serveFile(w http.ResponseWriter, name, mime string) {
	path := filepath.Join("./static", name)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprint(w, "<h1>请先复制 static/ 目录</h1>")
		return
	}
	setCORS(w)
	w.Header().Set("Content-Type", mime)
	w.Write(data)
}

// guessMIME 根据扩展名猜 MIME 类型
func guessMIME(path string) string {
	switch {
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".html"):
		return "text/html"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

// openAPISpec 返回 OpenAPI 3.0 规范（简化版）
func openAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":   "我的 API",
			"version": "1.0.0",
		},
		"paths": map[string]interface{}{
			"/api/hello": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "问候接口",
					"tags":    []string{"示例"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"message": map[string]interface{}{
												"type":    "string",
												"example": "Hello, Knife4j!",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
```

---

## 完整 CRUD 示例

参阅 [examples/go-stdlib/](../examples/go-stdlib/) 获取包含用户增删改查的完整示例。

---

## 自定义路由

### 方式一：简单字符串匹配（上面的示例）

适合小型项目，路由数量 < 20 个。

### 方式二：前缀匹配 + 方法分发

```go
func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// API 路由分发
	if strings.HasPrefix(path, "/api/users") {
		switch r.Method {
		case "GET":
			if path == "/api/users" {
				listUsers(w, r)
			} else {
				getUser(w, r)
			}
		case "POST":
			createUser(w, r)
		case "PUT":
			updateUser(w, r)
		case "DELETE":
			deleteUser(w, r)
		}
		return
	}

	// 其他路由...
}
```

### 方式三：使用 ServeMux（Go 1.22+）

Go 1.22 增强了 `http.ServeMux`，支持方法匹配和路径参数：

```go
mux := http.NewServeMux()

// 方法 + 路径匹配
mux.HandleFunc("GET /api/users", listUsers)
mux.HandleFunc("POST /api/users", createUser)
mux.HandleFunc("GET /api/users/{id}", getUser)
mux.HandleFunc("PUT /api/users/{id}", updateUser)
mux.HandleFunc("DELETE /api/users/{id}", deleteUser)

// Knife4j 端点
mux.HandleFunc("GET /v3/api-docs/swagger-config", swaggerConfig)
mux.HandleFunc("GET /openapi.json", openAPISpec)
mux.HandleFunc("GET /doc.html", serveDoc)

http.ListenAndServe(":8080", mux)
```

---

## CORS 完整配置

```go
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// 使用
http.HandleFunc("/", corsMiddleware(router))
```

---

## 静态文件服务

### 嵌入单个文件（Go 1.16+）

如果不想依赖 `static/` 目录，可以将 `doc.html` 嵌入二进制：

```go
import "embed"

//go:embed static/doc.html
var docHTML []byte

func serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(docHTML)
}
```

### 嵌入整个目录（Go 1.16+）

```go
//go:embed static/*
var staticFiles embed.FS

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	data, err := staticFiles.ReadFile("static/" + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", guessMIME(path))
	w.Write(data)
}
```

这样编译后的二进制文件自带前端资源，无需额外部署 `static/` 目录。

---

## 项目目录结构

```
my-go-app/
├── main.go
├── go.mod
└── static/
    ├── doc.html
    └── webjars/
        ├── js/
        │   └── *.js
        └── css/
            └── *.css
```

`go.mod` 最简形式（无依赖）：

```
module my-go-app

go 1.18
```

---

## 常见问题

### Q: Knife4j 显示 "No API definitions found"？

**A:** 检查 swagger-config 端点：

```bash
curl http://localhost:8080/v3/api-docs/swagger-config
```

应返回：
```json
{
  "urls": [{"url": "/openapi.json", "name": "default"}],
  "configUrl": "/v3/api-docs/swagger-config"
}
```

### Q: 如何添加新的 API 接口？

**A:** 两步：
1. 在 `router` 函数中添加路由匹配
2. 在 `openAPISpec()` 中添加对应的 OpenAPI 路径定义

### Q: 如何处理路径参数（如 `/api/users/123`）？

**A:** 使用字符串分割：

```go
func parseID(path string) int {
	// /api/users/123 → ["api", "users", "123"]
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		id, _ := strconv.Atoi(parts[2])
		return id
	}
	return 0
}
```

### Q: 如何添加请求体解析？

**A:** 使用 `json.NewDecoder`：

```go
var user User
json.NewDecoder(r.Body).Decode(&user)
defer r.Body.Close()
```

### Q: 如何生成请求响应示例？

**A:** 在 OpenAPI 规范的 schema 中添加 `example` 字段：

```go
"properties": map[string]interface{}{
    "name": map[string]interface{}{
        "type":    "string",
        "example": "张三",
    },
},
```

Knife4j 会自动读取并在界面中显示示例数据。

---

## 与 Gin 版本对比

| 特性 | 标准库版本 | Gin 版本 |
|------|-----------|---------|
| 依赖数量 | 0 | 1 (gin) |
| 路由能力 | 手动匹配 | 完整路由树 |
| 中间件 | 手动实现 | 内置支持 |
| 路径参数 | 手动解析 | 自动解析 |
| 适用场景 | 小型/嵌入式 | 中大型项目 |
| 编译体积 | 最小 | +5MB |
| 性能 | 极高 | 极高 |

---

## 完整示例

参阅 [examples/go-stdlib/](../examples/go-stdlib/) 获取完整可运行的示例项目。
