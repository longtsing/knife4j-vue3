# Go + Gin + Knife4j Vue3 示例

基于 Gin 框架集成 Knife4j Vue3 API 文档界面。

## 前提条件

- Go 1.22+
- 已编译的 Knife4j Vue3 前端（`dist/` 目录）

## 快速开始

### 1. 编译前端并复制静态文件

```bash
# 在项目根目录编译前端
cd knife4j-vue3
pnpm install && pnpm build

# 复制前端产物到 Go 示例的 static 目录
cp -r dist/* examples/go-gin/static/
```

> Windows 用户也可使用 `examples/setup.bat` 一键复制。

### 2. 初始化 Go 模块

```bash
cd examples/go-gin
go mod tidy
```

### 3. 启动服务

```bash
go run main.go
```

### 4. 访问文档

打开浏览器访问：http://localhost:8080/doc.html

OpenAPI JSON：http://localhost:8080/swagger/doc.json

## 项目结构

```
go-gin/
├── main.go           # 主程序（Gin 路由 + 静态文件托管）
├── go.mod            # Go 模块定义
├── go.sum            # 依赖校验
├── README.md         # 本文件
└── static/           # Knife4j 前端静态文件（需手动复制）
    ├── doc.html
    └── webjars/
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/:id | 根据 ID 获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |
| GET | /api/health | 健康检查 |

## 技术栈

- **Go 1.22+**
- **Gin v1.10.0** - Web 框架
- **swaggo** - Swagger 注解（仅编译时依赖）
- **Knife4j Vue3** - API 文档 UI
