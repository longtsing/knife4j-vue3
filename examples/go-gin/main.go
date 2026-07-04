package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// ============================================================
// 嵌入 Knife4j 前端静态文件
// 编译前需要将 dist/ 内容复制到 static/ 目录
// ============================================================

//go:embed all:static
var staticFS embed.FS

// ============================================================
// 数据模型
// ============================================================

// User 用户模型
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"张三"`
	Email string `json:"email" example:"zhangsan@example.com"`
	Role  string `json:"role" example:"admin"`
}

// UserCreate 创建用户请求体
type UserCreate struct {
	Name  string `json:"name" binding:"required" example:"张三"`
	Email string `json:"email" binding:"required" example:"zhangsan@example.com"`
	Role  string `json:"role" example:"admin"`
}

// Message 通用消息响应
type Message struct {
	Message string `json:"message" example:"success"`
}

// ============================================================
// 模拟数据库
// ============================================================

var (
	usersDB = make(map[int]*User)
	nextID  = 1
	mu      sync.RWMutex
)

func init() {
	// 初始化示例数据
	usersDB[1] = &User{ID: 1, Name: "张三", Email: "zhangsan@example.com", Role: "admin"}
	usersDB[2] = &User{ID: 2, Name: "李四", Email: "lisi@example.com", Role: "user"}
	nextID = 3
}

// ============================================================
// API 路由处理器
// ============================================================

// ListUsers 获取用户列表
// @Summary      获取用户列表
// @Description  返回系统中所有用户的信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200  {array}   User
// @Router       /api/users [get]
func ListUsers(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]*User, 0, len(usersDB))
	for _, u := range usersDB {
		list = append(list, u)
	}
	c.JSON(http.StatusOK, list)
}

// GetUser 根据ID获取用户
// @Summary      根据ID获取用户
// @Description  根据用户ID返回单个用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  User
// @Failure      404  {object}  Message
// @Router       /api/users/{id} [get]
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	mu.RLock()
	user, exists := usersDB[id]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建一个新的用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user  body      UserCreate  true  "用户信息"
// @Success      200   {object}  User
// @Failure      400   {object}  Message
// @Router       /api/users [post]
func CreateUser(c *gin.Context) {
	var req UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Message{Message: "参数错误: " + err.Error()})
		return
	}

	mu.Lock()
	user := &User{
		ID:    nextID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
	if user.Role == "" {
		user.Role = "user"
	}
	usersDB[nextID] = user
	nextID++
	mu.Unlock()

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  根据ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "用户ID"
// @Param        user  body      UserCreate  true  "用户信息"
// @Success      200   {object}  User
// @Failure      404   {object}  Message
// @Router       /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	var req UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Message{Message: "参数错误: " + err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	user, exists := usersDB[id]
	if !exists {
		c.JSON(http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  根据ID删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  Message
// @Failure      404  {object}  Message
// @Router       /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := usersDB[id]; !exists {
		c.JSON(http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}
	delete(usersDB, id)
	c.JSON(http.StatusOK, Message{Message: "删除成功"})
}

// HealthCheck 健康检查
// @Summary      健康检查
// @Description  服务健康状态检查
// @Tags         系统
// @Produce      json
// @Success      200  {object}  Message
// @Router       /api/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Message{Message: "ok"})
}

// ============================================================
// 主函数
// ============================================================

// @title           Knife4j Vue3 Go 示例 API
// @version         1.0.0
// @description     这是一个 Knife4j Vue3 + Go Gin 的示例项目，展示如何集成 API 文档界面。
// @host            localhost:8080
// @BasePath        /
func main() {
	r := gin.Default()

	// ============================================================
	// Knife4j 需要的 swagger-config 端点
	// ============================================================
	r.GET("/v3/api-docs/swagger-config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"urls": []gin.H{
				{
					"url":  "/swagger/doc.json",
					"name": "default",
				},
			},
			"configUrl":    "/v3/api-docs/swagger-config",
			"validatorUrl": "",
		})
	})

	// ============================================================
	// Knife4j 前端静态文件托管（使用 go:embed 嵌入）
	// 编译前需要将 dist/ 内容复制到 static/ 目录
	// ============================================================
	subFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal("嵌入静态文件系统初始化失败:", err)
	}
	httpFS := http.FS(subFS)

	// /doc.html 入口页面
	r.GET("/doc.html", func(c *gin.Context) {
		c.FileFromFS("doc.html", httpFS)
	})
	// /webjars/* 静态资源（JS/CSS）
	webjarsFS, _ := fs.Sub(subFS, "webjars")
	r.StaticFS("/webjars", http.FS(webjarsFS))
	// /oauth/* OAuth2 授权页面
	oauthFS, _ := fs.Sub(subFS, "oauth")
	r.StaticFS("/oauth", http.FS(oauthFS))
	// 根目录其他静态文件
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", httpFS)
	})
	r.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("robots.txt", httpFS)
	})

	// ============================================================
	// API 路由
	// ============================================================
	api := r.Group("/api")
	{
		api.GET("/users", ListUsers)
		api.GET("/users/:id", GetUser)
		api.POST("/users", CreateUser)
		api.PUT("/users/:id", UpdateUser)
		api.DELETE("/users/:id", DeleteUser)
		api.GET("/health", HealthCheck)
	}

	// ============================================================
	// 启动服务
	// ============================================================
	r.Run(":8080")
}
