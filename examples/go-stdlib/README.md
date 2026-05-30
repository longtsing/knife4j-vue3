# Go 标准库 + Knife4j Vue3 示例

**零依赖**：仅使用 Go 标准库（`net/http`、`encoding/json`），无需安装任何第三方包。

## 前提条件

- Go 1.18+
- 已编译的 Knife4j Vue3 前端（`dist/` 目录）

## 快速开始

```bash
# 1. 编译前端
cd knife4j-vue3
pnpm install && pnpm build

# 2. 复制静态文件
cp -r dist/* examples/go-stdlib/static/

# 3. 启动服务
cd examples/go-stdlib
go run main.go
```

访问 http://localhost:8080/doc.html

## 项目结构

```
go-stdlib/
├── main.go           # 主程序（纯标准库）
├── README.md         # 本文件
└── static/           # Knife4j 前端静态文件（需手动复制）
    ├── doc.html
    └── webjars/
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 根据ID获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/health | 健康检查 |

## 特点

- ✅ **零依赖**：不使用 Gin、Echo 等任何第三方框架
- ✅ **手写路由器**：基于 `http.HandleFunc` 的简单路由匹配
- ✅ **完整 CORS**：已配置跨域支持
- ✅ **静态文件服务**：自动检测 `static/` 目录并托管
- ✅ **OpenAPI 3.0**：手动构建完整规范，无需 swag 工具
