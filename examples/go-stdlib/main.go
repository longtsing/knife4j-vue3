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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserCreate 创建用户请求体
type UserCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Message 通用消息响应
type Message struct {
	Message string `json:"message"`
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
	usersDB[1] = &User{ID: 1, Name: "张三", Email: "zhangsan@example.com", Role: "admin"}
	usersDB[2] = &User{ID: 2, Name: "李四", Email: "lisi@example.com", Role: "user"}
	nextID = 3
}

// ============================================================
// 工具函数
// ============================================================

// writeJSON 写入 JSON 响应
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// readJSON 读取 JSON 请求体
func readJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// parseID 从 URL 路径中解析 ID
func parseID(path string) (int, error) {
	// 路径格式: /api/users/{id}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[2])
}

// ============================================================
// API 处理器
// ============================================================

// handleUsers 处理 /api/users 路由
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, Message{Message: "Method not allowed"})
	}
}

// handleUser 处理 /api/users/{id} 路由
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
		writeJSON(w, http.StatusMethodNotAllowed, Message{Message: "Method not allowed"})
	}
}

// listUsers 获取用户列表
func listUsers(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]*User, 0, len(usersDB))
	for _, u := range usersDB {
		list = append(list, u)
	}
	writeJSON(w, http.StatusOK, list)
}

// getUser 根据ID获取用户
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	mu.RLock()
	user, exists := usersDB[id]
	mu.RUnlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// createUser 创建用户
func createUser(w http.ResponseWriter, r *http.Request) {
	var req UserCreate
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, Message{Message: "参数错误: " + err.Error()})
		return
	}

	if req.Name == "" || req.Email == "" {
		writeJSON(w, http.StatusBadRequest, Message{Message: "name 和 email 不能为空"})
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

	writeJSON(w, http.StatusOK, user)
}

// updateUser 更新用户
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	var req UserCreate
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, Message{Message: "参数错误: " + err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	user, exists := usersDB[id]
	if !exists {
		writeJSON(w, http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}
	writeJSON(w, http.StatusOK, user)
}

// deleteUser 删除用户
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Message{Message: "无效的用户ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := usersDB[id]; !exists {
		writeJSON(w, http.StatusNotFound, Message{Message: "用户不存在"})
		return
	}
	delete(usersDB, id)
	writeJSON(w, http.StatusOK, Message{Message: "删除成功"})
}

// healthCheck 健康检查
func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Message{Message: "ok"})
}

// ============================================================
// Knife4j 端点
// ============================================================

// swaggerConfig 返回 Knife4j 需要的 swagger-config
func swaggerConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"urls": []map[string]string{
			{
				"url":  "/openapi.json",
				"name": "default",
			},
		},
		"configUrl":    "/v3/api-docs/swagger-config",
		"validatorUrl": "",
	}
	writeJSON(w, http.StatusOK, config)
}

// openAPISpec 返回 OpenAPI 3.0 规范
func openAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "Go 标准库 API",
			"version":     "1.0.0",
			"description": "Go 原生 net/http + Knife4j Vue3 示例，无需任何第三方包",
		},
		"servers": []map[string]string{
			{"url": "/", "description": "默认服务器"},
		},
		"paths": map[string]interface{}{
			"/api/users": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "获取用户列表",
					"description": "返回系统中所有用户的信息",
					"tags":        []string{"用户管理"},
					"operationId": "listUsers",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/User",
										},
									},
								},
							},
						},
					},
				},
				"post": map[string]interface{}{
					"summary":     "创建用户",
					"description": "创建一个新的用户",
					"tags":        []string{"用户管理"},
					"operationId": "createUser",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/UserCreate",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
					},
				},
			},
			"/api/users/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "根据ID获取用户",
					"description": "根据用户ID返回单个用户信息",
					"tags":        []string{"用户管理"},
					"operationId": "getUser",
					"parameters": []map[string]interface{}{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   map[string]string{"type": "integer"},
							"example":  1,
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "用户不存在",
						},
					},
				},
				"put": map[string]interface{}{
					"summary":     "更新用户",
					"description": "根据ID更新用户信息",
					"tags":        []string{"用户管理"},
					"operationId": "updateUser",
					"parameters": []map[string]interface{}{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   map[string]string{"type": "integer"},
						},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/UserCreate",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
					},
				},
				"delete": map[string]interface{}{
					"summary":     "删除用户",
					"description": "根据ID删除用户",
					"tags":        []string{"用户管理"},
					"operationId": "deleteUser",
					"parameters": []map[string]interface{}{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema":   map[string]string{"type": "integer"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "删除成功",
						},
					},
				},
			},
			"/api/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "健康检查",
					"description": "服务健康状态检查",
					"tags":        []string{"系统"},
					"operationId": "healthCheck",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "成功",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Message",
									},
								},
							},
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"User": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id":    map[string]interface{}{"type": "integer", "example": 1},
						"name":  map[string]interface{}{"type": "string", "example": "张三"},
						"email": map[string]interface{}{"type": "string", "example": "zhangsan@example.com"},
						"role":  map[string]interface{}{"type": "string", "example": "admin"},
					},
				},
				"UserCreate": map[string]interface{}{
					"type":     "object",
					"required": []string{"name", "email"},
					"properties": map[string]interface{}{
						"name":  map[string]interface{}{"type": "string", "example": "张三"},
						"email": map[string]interface{}{"type": "string", "example": "zhangsan@example.com"},
						"role":  map[string]interface{}{"type": "string", "example": "admin"},
					},
				},
				"Message": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message": map[string]interface{}{"type": "string", "example": "success"},
					},
				},
			},
		},
	}
	writeJSON(w, http.StatusOK, spec)
}

// ============================================================
// 路由器
// ============================================================

// router 简单的路径路由器
func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// CORS 预检请求
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
	case path == "/openapi.json":
		openAPISpec(w, r)

	// API 路由
	case path == "/api/users":
		handleUsers(w, r)
	case strings.HasPrefix(path, "/api/users/"):
		handleUser(w, r)
	case path == "/api/health":
		healthCheck(w, r)

	// Knife4j 静态文件（从 embed.FS 读取）
	case path == "/doc.html" || path == "/doc.html/":
		serveFile(w, r, "doc.html", "text/html")
	case strings.HasPrefix(path, "/webjars/"):
		serveStatic(w, r, path)
	case strings.HasPrefix(path, "/oauth/"):
		serveStatic(w, r, path)
	case path == "/favicon.ico":
		serveFile(w, r, "favicon.ico", "image/x-icon")
	case path == "/robots.txt":
		serveFile(w, r, "robots.txt", "text/plain")

	// 默认
	default:
		http.NotFound(w, r)
	}
}

// ============================================================
// 静态文件服务（从 embed.FS 读取）
// ============================================================

// subStaticFS 嵌入的 static 子文件系统
var subStaticFS fs.FS

func init() {
	var err error
	subStaticFS, err = fs.Sub(staticFS, "static")
	if err != nil {
		log.Println("警告: 嵌入静态文件系统初始化失败:", err)
	}
}

// notFoundHTML 当嵌入的前端文件不存在时的提示页面
const notFoundHTML = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Knife4j Vue3 + Go</title></head>
<body style="font-family:sans-serif;padding:40px;max-width:700px;margin:0 auto;">
<h1>Knife4j Vue3 + Go 标准库示例</h1>
<p>未检测到嵌入的前端文件，请重新编译：</p>
<pre style="background:#f5f5f5;padding:16px;border-radius:8px;">
# 1. 编译前端
cd knife4j-vue3
pnpm install &amp;&amp; pnpm build

# 2. 复制产物
cp -r dist/* examples/go-stdlib/static/

# 3. 重新编译 Go（将前端嵌入到二进制）
cd examples/go-stdlib
go build -o server main.go
./server
</pre>
<p>然后访问 <a href="/doc.html">/doc.html</a></p>
</body>
</html>`

// serveFile 从嵌入的 FS 提供单个文件
func serveFile(w http.ResponseWriter, r *http.Request, name, contentType string) {
	data, err := fs.ReadFile(subStaticFS, name)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, notFoundHTML)
		return
	}
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.Write(data)
}

// serveStatic 从嵌入的 FS 提供静态资源（如 /webjars/js/xxx.js）
func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	// 去掉前导 / 以匹配 embed.FS 的相对路径
	relPath := strings.TrimPrefix(path, "/")
	data, err := fs.ReadFile(subStaticFS, relPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 根据扩展名设置 Content-Type
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
	case strings.HasSuffix(path, ".woff"):
		contentType = "font/woff"
	case strings.HasSuffix(path, ".woff2"):
		contentType = "font/woff2"
	case strings.HasSuffix(path, ".ttf"):
		contentType = "font/ttf"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// ============================================================
// 主函数
// ============================================================

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Knife4j Vue3 + Go 标准库示例启动\n")
	fmt.Printf("📖 文档地址: http://localhost:%s/doc.html\n", port)
	fmt.Printf("🔗 API 地址: http://localhost:%s/api/users\n", port)
	fmt.Printf("📋 OpenAPI:  http://localhost:%s/openapi.json\n", port)
	fmt.Printf("\n")

	if err := http.ListenAndServe(":"+port, http.HandlerFunc(router)); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
