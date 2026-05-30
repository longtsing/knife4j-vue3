# Go Gin 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Go Gin 后端集成。

## 前提条件

- Go 1.22+
- Gin 框架
- 已编译的 Knife4j Vue3 前端

## 方式一：手动提供 OpenAPI JSON（推荐）

不依赖 swag 工具，直接在代码中定义 OpenAPI 规范并暴露端点。

### 1. 安装依赖

```bash
go get github.com/gin-gonic/gin
```

### 2. 项目结构

```
my-go-app/
├── main.go
├── openapi.go          # OpenAPI 规范定义
├── go.mod
├── go.sum
└── static/             # Knife4j 前端静态文件
    ├── doc.html
    └── webjars/
```

### 3. 编写 OpenAPI 规范

```go
// openapi.go
package main

import "github.com/gin-gonic/gin"

// OpenAPISpec 返回 OpenAPI 3.0 规范
func OpenAPISpec() gin.H {
	return gin.H{
		"openapi": "3.0.3",
		"info": gin.H{
			"title":       "My Go API",
			"version":     "1.0.0",
			"description": "Go Gin + Knife4j Vue3 示例 API",
		},
		"servers": []gin.H{
			{"url": "/", "description": "默认服务器"},
		},
		"paths": gin.H{
			"/api/users": gin.H{
				"get": gin.H{
					"summary":     "获取用户列表",
					"description": "返回系统中所有用户的信息",
					"tags":        []string{"用户管理"},
					"operationId": "listUsers",
					"responses": gin.H{
						"200": gin.H{
							"description": "成功",
							"content": gin.H{
								"application/json": gin.H{
									"schema": gin.H{
										"type": "array",
										"items": gin.H{"$ref": "#/components/schemas/User"},
									},
								},
							},
						},
					},
				},
				"post": gin.H{
					"summary":     "创建用户",
					"description": "创建一个新的用户",
					"tags":        []string{"用户管理"},
					"operationId": "createUser",
					"requestBody": gin.H{
						"required": true,
						"content": gin.H{
							"application/json": gin.H{
								"schema": gin.H{"$ref": "#/components/schemas/UserCreate"},
							},
						},
					},
					"responses": gin.H{
						"200": gin.H{
							"description": "成功",
							"content": gin.H{
								"application/json": gin.H{
									"schema": gin.H{"$ref": "#/components/schemas/User"},
								},
							},
						},
					},
				},
			},
			"/api/users/{id}": gin.H{
				"get": gin.H{
					"summary":     "根据ID获取用户",
					"tags":        []string{"用户管理"},
					"operationId": "getUser",
					"parameters": []gin.H{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   gin.H{"type": "integer"},
						},
					},
					"responses": gin.H{
						"200": gin.H{
							"description": "成功",
							"content": gin.H{
								"application/json": gin.H{
									"schema": gin.H{"$ref": "#/components/schemas/User"},
								},
							},
						},
						"404": gin.H{"description": "用户不存在"},
					},
				},
				"put": gin.H{
					"summary":     "更新用户",
					"tags":        []string{"用户管理"},
					"operationId": "updateUser",
					"parameters": []gin.H{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   gin.H{"type": "integer"},
						},
					},
					"requestBody": gin.H{
						"required": true,
						"content": gin.H{
							"application/json": gin.H{
								"schema": gin.H{"$ref": "#/components/schemas/UserCreate"},
							},
						},
					},
					"responses": gin.H{
						"200": gin.H{
							"description": "成功",
							"content": gin.H{
								"application/json": gin.H{
									"schema": gin.H{"$ref": "#/components/schemas/User"},
								},
							},
						},
					},
				},
				"delete": gin.H{
					"summary":     "删除用户",
					"tags":        []string{"用户管理"},
					"operationId": "deleteUser",
					"parameters": []gin.H{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   gin.H{"type": "integer"},
						},
					},
					"responses": gin.H{
						"200": gin.H{"description": "删除成功"},
					},
				},
			},
		},
		"components": gin.H{
			"schemas": gin.H{
				"User": gin.H{
					"type": "object",
					"properties": gin.H{
						"id":    gin.H{"type": "integer", "example": 1},
						"name":  gin.H{"type": "string", "example": "张三"},
						"email": gin.H{"type": "string", "example": "zhangsan@example.com"},
						"role":  gin.H{"type": "string", "example": "admin"},
					},
				},
				"UserCreate": gin.H{
					"type": "object",
					"required": []string{"name", "email"},
					"properties": gin.H{
						"name":  gin.H{"type": "string", "example": "张三"},
						"email": gin.H{"type": "string", "example": "zhangsan@example.com"},
						"role":  gin.H{"type": "string", "example": "admin"},
					},
				},
			},
		},
	}
}
```

### 4. 编写主程序

```go
// main.go
package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserCreate 创建用户请求
type UserCreate struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
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

func main() {
	r := gin.Default()

	// Knife4j swagger-config 端点
	r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"urls": []gin.H{
				{"url": "/openapi.json", "name": "default"},
			},
			"configUrl":   "/v3/api-docs/swagger-config",
			"validatorUrl": "",
		})
	})

	// OpenAPI 规范端点
	r.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, OpenAPISpec())
	})

	// Knife4j 静态文件
	r.Static("/static", "./static")
	r.GET("/doc.html", func(c *gin.Context) {
		c.File("./static/doc.html")
	})
	r.Static("/webjars", "./static/webjars")

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/users", func(c *gin.Context) {
			mu.RLock()
			defer mu.RUnlock()
			list := make([]*User, 0, len(usersDB))
			for _, u := range usersDB {
				list = append(list, u)
			}
			c.JSON(http.StatusOK, list)
		})

		api.GET("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			mu.RLock()
			user, exists := usersDB[id]
			mu.RUnlock()
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
				return
			}
			c.JSON(http.StatusOK, user)
		})

		api.POST("/users", func(c *gin.Context) {
			var req UserCreate
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			mu.Lock()
			user := &User{ID: nextID, Name: req.Name, Email: req.Email, Role: req.Role}
			if user.Role == "" {
				user.Role = "user"
			}
			usersDB[nextID] = user
			nextID++
			mu.Unlock()
			c.JSON(http.StatusOK, user)
		})

		api.DELETE("/users/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			mu.Lock()
			defer mu.Unlock()
			if _, exists := usersDB[id]; !exists {
				c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
				return
			}
			delete(usersDB, id)
			c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
		})
	}

	r.Run(":8080")
}
```

---

## 方式二：使用 swag 工具自动生成

如果你更喜欢使用注解自动生成 OpenAPI 规范，可以使用 [swaggo/swag](https://github.com/swaggo/swag)。

### 1. 安装 swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 添加依赖

```bash
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 3. 在 main.go 中添加注解

```go
// @title           My Go API
// @version         1.0.0
// @description     Go Gin + Knife4j Vue3 示例 API
// @host            localhost:8080
// @BasePath        /

// @Summary      获取用户列表
// @Description  返回所有用户
// @Tags         用户管理
// @Produce      json
// @Success      200  {array}   User
// @Router       /api/users [get]
func ListUsers(c *gin.Context) {
	// ...
}
```

### 4. 生成 Swagger 文档

```bash
swag init
```

这会生成 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml` 文件。

### 5. 注册路由

```go
import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "my-app/docs"  // 导入生成的文档
)

func main() {
	r := gin.Default()

	// Knife4j swagger-config 端点（指向 swag 生成的 JSON）
	r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"urls": []gin.H{
				{"url": "/swagger/doc.json", "name": "default"},
			},
			"configUrl": "/v3/api-docs/swagger-config",
		})
	})

	// Knife4j 静态文件
	r.Static("/static", "./static")
	r.GET("/doc.html", func(c *gin.Context) {
		c.File("./static/doc.html")
	})

	// Swagger JSON 端点（swag 生成）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由...

	r.Run(":8080")
}
```

---

## Knife4j 集成原理

Knife4j Vue3 前端需要一个 `swagger-config` 端点来获取 OpenAPI 文档的 URL：

```
GET /v3/api-docs/swagger-config
→ 返回: { "urls": [{ "url": "/openapi.json", "name": "default" }] }

GET /openapi.json
→ 返回: OpenAPI 3.0 完整规范

GET /doc.html
→ 返回: Knife4j Vue3 前端页面
```

前端加载流程：

```
doc.html 加载
  → 请求 /v3/api-docs/swagger-config
  → 获取 OpenAPI JSON 地址
  → 请求 /openapi.json 获取完整规范
  → 渲染 API 文档界面
```

---

## CORS 配置

```go
r.Use(func(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
})
```

---

## 常见问题

### Q: Knife4j 显示 "No API definitions found"？

**A:** 检查 swagger-config 端点是否正确：

```bash
curl http://localhost:8080/v3/api-docs/swagger-config
```

应该返回包含 OpenAPI JSON 地址的配置。

### Q: 如何自定义 API 文档信息？

**A:** 修改 OpenAPI 规范中的 `info` 字段：

```go
"info": gin.H{
    "title":       "我的 API 文档",
    "version":     "1.0.0",
    "description": "项目描述",
    "contact": gin.H{
        "name":  "技术支持",
        "email": "support@example.com",
    },
},
```

### Q: 如何添加认证？

**A:** 在 OpenAPI 规范中添加 securitySchemes，Knife4j 界面会自动显示认证配置入口：

```go
"components": gin.H{
    "securitySchemes": gin.H{
        "bearerAuth": gin.H{
            "type":         "http",
            "scheme":       "bearer",
            "bearerFormat": "JWT",
        },
    },
},
```

---

## 完整示例

参阅 [examples/go/](../examples/go/) 目录获取完整可运行的示例项目。
