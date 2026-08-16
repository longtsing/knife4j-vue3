# Go Gin 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Go Gin 后端集成，使用 **运行时构造 OpenAPI** 的方案。

## 方案对比

| 方案 | 工具链 | 优点 | 缺点 |
|------|--------|------|------|
| **运行时构造（推荐）** | 无 | 零依赖，编译期类型安全，Schema 可复用 | 需手写一遍元信息 |
| swaggo 注解 | swag CLI | 注解贴近业务 | 需额外 CLI 步骤，需引入 3 个第三方包 |

本示例采用 **运行时构造** 方案：把 OpenAPI 规范当作一等数据结构放在 `docs/openapi.go` 里，路由注册时旁路调用 `docs.Register(...)` 声明端点。

## 前提条件

- Go 1.22+
- `github.com/gin-gonic/gin`
- 已编译的 Knife4j Vue3 前端

## 1. 安装依赖

```bash
go get github.com/gin-gonic/gin
```

## 2. 项目结构

```
my-go-app/
├── main.go              # 业务入口
├── docs/
│   └── openapi.go       # OpenAPI 规范构造器（运行时）
├── go.mod
├── go.sum
├── static/              # Knife4j 前端产物（不纳入版本控制）
│   ├── doc.html
│   ├── favicon.ico
│   ├── robots.txt
│   ├── webjars/
│   └── oauth/
└── README.md
```

## 3. 编写 OpenAPI 构造器 `docs/openapi.go`

完整文件见 [examples/go-gin/docs/openapi.go](../examples/go-gin/docs/openapi.go)。

核心设计：

```go
// 1) 用结构体表示 OpenAPI 3.0（最小可用子集）
type Spec struct {
    OpenAPI    string           `json:"openapi"`
    Info       Info             `json:"info"`
    Servers    []Server         `json:"servers"`
    Tags       []Tag            `json:"tags"`
    Components *Components      `json:"components"`
    Paths      map[string]*Path `json:"paths"`
}

// 2) 一次性构造，sync.Once 保护并发
var (once sync.Once; spec *Spec)
func Get() *Spec { once.Do(build); return spec }

// 3) 注册端点
func Register(method, path string, op *Operation) { ... }

// 4) 链式 Builder
func Op(summary string, tags ...string) *Operation { ... }
func (o *Operation) Param(name, in, desc string, required bool, schema *Schema) *Operation { ... }
func (o *Operation) JSONBody(schema *Schema, required bool) *Operation { ... }
func (o *Operation) Res(status, description string, schema *Schema) *Operation { ... }

// 5) 复用 Schema
func AddSchema(name string, s *Schema) { ... }
func Ref(name string) *Schema { ... }
```

### ⚠️ 死锁陷阱

**`registerEndpoints()` 不能放在 `build()` 内部**，否则会触发 `sync.Once` 自递归死锁。必须拆分到 `init()`：

```go
func build() {
    spec = &Spec{...}
    AddSchema("User", User)
    // 注意：这里不能调用 registerEndpoints()
}

func init() {
    once.Do(build)         // 锁在此释放
    registerEndpoints()    // 内部 Register() 的 once.Do(build) 立即返回
}
```

## 4. 声明端点和 Schema

```go
// docs/openapi.go
// 可复用 Schema
var User = Obj(map[string]*Schema{
    "id":    Int("用户ID", "int64"),
    "name":  Str("用户名", ""),
    "email": Str("邮箱", "email"),
    "role":  Str("角色", "", "admin", "user"),
}, "id", "name", "email")

var UserCreate = Obj(map[string]*Schema{
    "name":  Str("用户名", ""),
    "email": Str("邮箱", "email"),
    "role":  Str("角色", "", "admin", "user"),
}, "name", "email")

// 集中声明所有端点（与 main.go 路由保持同步）
func registerEndpoints() {
    Register("GET", "/api/users",
        Op("获取用户列表", "用户管理").
            Res("200", "成功", Arr(Ref("User"), "")))

    Register("GET", "/api/users/{id}",
        Op("根据ID获取用户", "用户管理").
            Param("id", "path", "用户ID", true, Int("", "int64")).
            Res("200", "成功", Ref("User")).
            Res("404", "用户不存在", Ref("Message")))

    Register("POST", "/api/users",
        Op("创建用户", "用户管理").
            JSONBody(Ref("UserCreate"), true).
            Res("200", "创建成功", Ref("User")))

    // ... 其他端点
}
```

## 5. 业务入口 `main.go`

```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
    "strconv"
    "sync"

    "github.com/gin-gonic/gin"
    "my-go-app/docs"
)

// 1) 嵌入 Knife4j 前端静态文件
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

// 3) Handler（省略 CRUD 实现）
func ListUsers(c *gin.Context) { /* ... */ }

func main() {
    r := gin.Default()

    // 4) Knife4j 所需的 swagger-config 端点
    r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "urls": []gin.H{
                {"url": "/swagger/doc.json", "name": "default"},
            },
            "configUrl":   "/v3/api-docs/swagger-config",
            "validatorUrl": "",
        })
    })

    // 5) OpenAPI 规范端点（运行时从 docs 包生成）
    r.GET("/swagger/doc.json", func(c *gin.Context) {
        c.Header("Content-Type", "application/json; charset=utf-8")
        _, _ = c.Writer.Write(docs.MustJSON())
    })

    // 6) 嵌入前端静态文件
    subFS, _ := fs.Sub(staticFS, "static")
    httpFS := http.FS(subFS)
    r.GET("/doc.html", func(c *gin.Context) { c.FileFromFS("doc.html", httpFS) })
    webjarsFS, _ := fs.Sub(subFS, "webjars")
    r.StaticFS("/webjars", http.FS(webjarsFS))

    // 7) API 路由
    api := r.Group("/api")
    {
        api.GET("/users", ListUsers)
        // ... 其他路由
    }

    log.Printf("🚀 服务启动: http://localhost:8080/doc.html")
    r.Run(":8080")
}
```

## 6. 启动

```bash
# 1. 在 knife4j-vue3 根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build

# 2. 复制前端产物到本项目
cp -r dist/* my-go-app/static/

# 3. 启动服务
cd my-go-app
go mod tidy
go run main.go
```

访问：http://localhost:8080/doc.html

## 7. Knife4j 集成原理

```
浏览器访问 /doc.html
  ↓
请求 /v3/api-docs/swagger-config → 拿到 OpenAPI JSON 地址 /swagger/doc.json
  ↓
请求 /swagger/doc.json → docs.MustJSON() 生成规范
  ↓
Knife4j 渲染 API 文档界面
```

调试时，Knife4j 内部 ajax 会自动调用 `/api/users` 等接口。前端已修复重复前缀 bug，swagger-config 返回 `/swagger/doc.json` 这样的绝对路径不会再被重复拼接。

## 8. 高级用法

### 8.1 Bearer 认证

```go
// 1) 在 OpenAPI 中声明安全方案
spec.Components.SecuritySchemes = map[string]*SecurityScheme{
    "bearerAuth": {
        Type:        "http",
        Scheme:      "bearer",
        Description: "登录后获取的不透明 token",
    },
}

// 2) 给需要鉴权的接口加 Sec()
Register("POST", "/api/users",
    Op("创建用户", "用户管理").
        Sec("bearerAuth").
        JSONBody(Ref("UserCreate"), true).
        Res("200", "成功", Ref("User")))
```

Knife4j 界面会自动显示 Authorize 按钮；也可通过「文档管理」→「全局参数设置」添加 Authorization Header。

### 8.2 重命名端点 operationId

Knife4j 默认按函数名生成 operationId。如果想自定义：

```go
// 在 Op() 后链式调用 Id()
Register("GET", "/api/users",
    Op("获取用户列表", "用户管理").Id("listUsers").
        Res("200", "成功", Arr(Ref("User"), "")))
```

（需要先在 Operation 结构体添加 `Id string` 字段并标记 `json:"operationId,omitempty"`）

### 8.3 数组 + 引用 + 必填字段

```go
Arr(Ref("User"), "用户列表")           // []User
Obj({...}, "id", "name")              // 标记 id、name 为 required
Str("用户名", "", "admin", "user")    // 带 enum 枚举
Int("年龄", "int64")                  // 带 format
```

## 9. 完整示例

[examples/go-gin/](../examples/go-gin/) 提供开箱即用的可运行示例：

```bash
cd examples/go-gin
cp -r ../../dist/* static/
go mod tidy
go run main.go
# 访问 http://localhost:8080/doc.html
```

## 10. 常见问题

### Q1：报错 `Handler already registered for path ... and http method OPTIONS`？

A：Knife4j 会发 CORS 预检请求（OPTIONS），你的路由如果用了 `Method("GET")` 之类限制方法，需要加上 `OPTIONS` 支持：

```go
api := r.Group("/api", func(c *gin.Context) {
    if c.Request.Method == http.MethodOptions {
        c.AbortWithStatus(http.StatusNoContent)
        return
    }
    c.Next()
})
```

或使用 Gin 的 `cors` 中间件统一处理。

### Q2：调试面板响应区空白？

A：检查浏览器 Network，确认请求 URL 没有重复前缀（`/api/api/...`）。如果是 `/api` 前缀项目，确保：

1. swagger-config 返回的 URL 已是完整路径（含 `/api`）
2. 后端能直接访问 `/api/<url>` 返回 200

参考 [examples/go-gin/main.go](../examples/go-gin/main.go) 的写法。

### Q3：go:embed 报 "pattern all:static matched no files"？

A：`go:embed` 要求目录在编译时存在。`static/` 是空目录或缺失都会报错。请先按上面"启动"步骤复制前端产物。

### Q4：端口冲突？

A：Java、Go 示例都用 8080，Python 用 8000。可在 `r.Run(":9090")` 修改 Go 端口。

### Q5：如何添加请求示例？

A：在 schema 里用 `example` 字段（Gin 示例见 [User 模型](../examples/go-gin/main.go)）。
