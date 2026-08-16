# Go 标准库对接指南（零依赖）

> **零依赖方案**：仅使用 Go 标准库 `net/http` + `encoding/json` + `embed`，无需任何第三方包。

## 方案对比

| 方案 | 第三方依赖 | 适用场景 |
|------|-----------|---------|
| **标准库 + 运行时 OpenAPI**（本示例） | 0 | 嵌入式、CLI 工具、学习项目 |
| Gin + docs/openapi.go | 1 (gin) | 中大型 Web 服务 |
| swaggo 注解 | 3 (swag/gin-swagger/files) | 想贴近 Java Swagger 注解风格的团队 |

本示例是 `go-gin` 的零依赖版本：业务代码完全使用 `net/http`，OpenAPI 规范同样用运行时构造（共用 [examples/go-stdlib/docs/openapi.go](../examples/go-stdlib/docs/openapi.go)）。

## 前提条件

- Go 1.22+（用到了 `http.ServeMux` 的方法匹配语法）
- Knife4j Vue3 前端编译产物

## 1. 项目结构

```
my-go-app/
├── main.go              # 业务入口（含路由器 + handler）
├── docs/
│   └── openapi.go       # OpenAPI 规范构造器
├── go.mod               # 只有 module + go 行
├── static/              # Knife4j 前端产物（不纳入版本控制）
│   ├── doc.html
│   ├── favicon.ico
│   ├── robots.txt
│   ├── webjars/
│   └── oauth/
└── README.md
```

## 2. 核心原理

Knife4j Vue3 前端需要以下端点：

| 端点 | 作用 |
|------|------|
| `/v3/api-docs/swagger-config` | 返回 OpenAPI JSON 的地址 |
| `/swagger/doc.json` | OpenAPI 3.0 完整规范（运行时生成） |
| `/doc.html` | Knife4j 入口页面 |
| `/webjars/*`、`/oauth/*` | 前端 JS/CSS 静态资源 |
| `/api/*` | 业务 API |

## 3. OpenAPI 构造器 `docs/openapi.go`

与 go-gin 共用同样的设计：[examples/go-stdlib/docs/openapi.go](../examples/go-stdlib/docs/openapi.go)。

主要 API：

```go
docs.Register("GET", "/api/users", op)   // 注册端点
docs.MustJSON()                          // 序列化为 JSON
docs.AddSchema("User", userSchema)      // 注册可复用 schema
```

**关键陷阱**：`registerEndpoints()` 不能放在 `build()` 内部，否则 `sync.Once` 会自递归死锁。详见 [go-gin.md](./go-gin.md#死锁陷阱)。

## 4. 业务入口 `main.go`

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "io/fs"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"

    "my-go-app/docs"
)

// 1) 嵌入 Knife4j 前端静态文件（编译期打包进二进制）
//go:embed all:static
var staticFS embed.FS

// 2) 数据模型
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

var (
    usersDB = make(map[int]*User)
    nextID  = 1
    mu      sync.RWMutex
)

func init() {
    usersDB[1] = &User{ID: 1, Name: "张三", Email: "zhangsan@example.com", Role: "admin"}
    usersDB[2] = &User{ID: 2, Name: "李四", Email: "lisi@example.com", Role: "user"}
    nextID = 3
}

// 3) 工具函数
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func readJSON(r *http.Request, v any) error {
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(v)
}

func parseID(path string) (int, error) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) < 3 {
        return 0, fmt.Errorf("invalid path")
    }
    return strconv.Atoi(parts[2])
}

// 4) API handler
func handleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        listUsers(w, r)
    case http.MethodPost:
        createUser(w, r)
    case http.MethodOptions:
        w.WriteHeader(http.StatusNoContent)
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "Method not allowed"})
    }
}

func handleUser(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUser(w, r)
    case http.MethodPut:
        updateUser(w, r)
    case http.MethodDelete:
        deleteUser(w, r)
    case http.MethodOptions:
        w.WriteHeader(http.StatusNoContent)
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "Method not allowed"})
    }
}

func listUsers(w http.ResponseWriter, _ *http.Request) {
    mu.RLock()
    defer mu.RUnlock()
    list := make([]*User, 0, len(usersDB))
    for _, u := range usersDB {
        list = append(list, u)
    }
    writeJSON(w, http.StatusOK, list)
}

// 省略 getUser / createUser / updateUser / deleteUser / healthCheck ...

// 5) Knife4j swagger-config 端点
func swaggerConfig(w http.ResponseWriter, _ *http.Request) {
    writeJSON(w, http.StatusOK, map[string]any{
        "urls": []map[string]string{
            {"url": "/swagger/doc.json", "name": "default"},
        },
        "configUrl":   "/v3/api-docs/swagger-config",
        "validatorUrl": "",
    })
}

// 6) 嵌入 FS 子文件系统
var subStaticFS fs.FS

func init() {
    var err error
    subStaticFS, err = fs.Sub(staticFS, "static")
    if err != nil {
        log.Println("警告: 嵌入静态文件系统初始化失败:", err)
    }
}

// 7) 静态文件服务
func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
    relPath := strings.TrimPrefix(path, "/")
    data, err := fs.ReadFile(subStaticFS, relPath)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    contentType := "application/octet-stream"
    switch {
    case strings.HasSuffix(path, ".js"):
        contentType = "application/javascript"
    case strings.HasSuffix(path, ".css"):
        contentType = "text/css"
    case strings.HasSuffix(path, ".html"):
        contentType = "text/html"
    case strings.HasSuffix(path, ".json"):
        contentType = "application/json"
    case strings.HasSuffix(path, ".svg"):
        contentType = "image/svg+xml"
    case strings.HasSuffix(path, ".png"):
        contentType = "image/png"
    case strings.HasSuffix(path, ".ico"):
        contentType = "image/x-icon"
    case strings.HasSuffix(path, ".woff2"):
        contentType = "font/woff2"
    }
    w.Header().Set("Content-Type", contentType)
    w.Write(data)
}

// 8) 简单路由器
func router(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    // CORS 预检
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    switch {
    // Knife4j 端点
    case path == "/v3/api-docs/swagger-config":
        swaggerConfig(w, r)
    case path == "/swagger/doc.json":
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write(docs.MustJSON())

    // 业务 API
    case path == "/api/users":
        handleUsers(w, r)
    case strings.HasPrefix(path, "/api/users/"):
        handleUser(w, r)
    case path == "/api/health":
        writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})

    // 静态文件
    case path == "/doc.html":
        serveStatic(w, r, "doc.html")
    case strings.HasPrefix(path, "/webjars/"), strings.HasPrefix(path, "/oauth/"):
        serveStatic(w, r, path)
    case path == "/favicon.ico":
        serveStatic(w, r, "favicon.ico")

    default:
        http.NotFound(w, r)
    }
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Printf("🚀 服务启动: http://localhost:%s/doc.html\n", port)
    log.Fatal(http.ListenAndServe(":"+port, http.HandlerFunc(router)))
}
```

## 5. 启动

```bash
# 1. 在 knife4j-vue3 根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build

# 2. 复制前端产物到本项目
cp -r dist/* my-go-app/static/

# 3. 启动服务（零依赖，go run 即可）
cd my-go-app
go run main.go
```

访问：http://localhost:8080/doc.html

## 6. 进阶：Go 1.22+ ServeMux 方法匹配

Go 1.22 起 `http.ServeMux` 支持方法+路径精确匹配，可以替代手写的字符串 switch：

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
mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Write(docs.MustJSON())
})

http.ListenAndServe(":8080", mux)
```

路径参数通过 `r.PathValue("id")` 读取。

## 7. 与 Gin 版对比

| 特性 | 标准库版 | Gin 版 |
|------|---------|--------|
| 第三方依赖 | 0 | 1 (gin) |
| OpenAPI 构造 | `docs/openapi.go` 运行时构造 | 同上 |
| 路由能力 | Go 1.22+ ServeMux 或手写 | 完整路由树 |
| 中间件 | 手写包装 | 内置支持 |
| 路径参数 | Go 1.22+ `r.PathValue` 或手解析 | `c.Param` |
| 嵌入静态 | `embed.FS` | `embed.FS` |
| 编译产物体积 | ~5 MB | ~10 MB |
| 适用场景 | 嵌入式、CLI、零依赖部署 | 中大型 Web 服务 |

## 8. 完整示例

[examples/go-stdlib/](../examples/go-stdlib/) 提供开箱即用的可运行示例：

```bash
cd examples/go-stdlib
cp -r ../../dist/* static/
go run main.go
# 访问 http://localhost:8080/doc.html
```

## 9. 常见问题

### Q1：报 "pattern all:static matched no files"？

A：`go:embed` 要求目录在编译时存在，且至少有 1 个文件。请先按上面"启动"步骤把 `dist/*` 复制到 `static/`。

### Q2：能打包成单个二进制吗？

A：可以。所有静态资源已通过 `//go:embed all:static` 嵌入二进制。部署时只需要一个文件：

```bash
go build -o myserver main.go
./myserver
```

### Q3：前端 ajax 请求被浏览器拦截？

A：标准库版默认没有 CORS 中间件。本示例在 `writeJSON` 中已经设置 `Access-Control-Allow-Origin: *`，覆盖 Knife4j 调试场景。如需生产级 CORS，请用 `cors` 中间件包裹 `router`。

### Q4：调试面板响应区空白？

A：检查浏览器 Network，确认请求 URL 没有重复前缀。前端已修复重复前缀 bug，swagger-config 返回的 `/swagger/doc.json`（绝对路径）不会再被拼接成 `/swagger/doc.json`。