# Go + Gin + Knife4j Vue3 示例

## 前提条件

- Go 1.22+
- 已编译的 Knife4j Vue3 前端（`dist/` 目录）

## 快速开始

### 1. 编译前端并复制静态文件

```bash
# 编译前端
cd knife4j-vue3
pnpm install && pnpm build

# 复制到 Go 示例的 static 目录
cp -r dist/* examples/go/static/
```

### 2. 初始化 Go 模块

```bash
cd examples/go
go mod tidy
```

### 3. 启动服务

```bash
go run main.go
```

### 4. 访问文档

打开浏览器访问：http://localhost:8080/doc.html

## 项目结构

```
go/
├── main.go           # 主程序（含 Swagger 注解）
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
| GET | /api/users/:id | 根据ID获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |
| GET | /api/health | 健康检查 |

## 注意事项

1. Go 示例使用 `swaggo` 风格的注解，但**不依赖 swag 工具**，直接手动定义 OpenAPI 端点
2. `/v3/api-docs/swagger-config` 端点返回 OpenAPI 文档的 URL
3. Knife4j 前端通过该端点获取文档地址并加载 API 定义
