# Knife4j Vue3 示例项目

本目录包含 Knife4j Vue3 与 5 个不同后端框架集成的完整示例。

## 总览

| 示例 | 语言 | 框架 | 端口 | 访问地址 | 前端产物处理 |
|------|------|------|------|----------|-------------|
| [java-springboot](./java-springboot/) | Java 17 | Spring Boot 3.2 | 8080 | `http://localhost:8080/api/doc.html` | Maven 插件自动复制 |
| [fastapi](./fastapi/) | Python 3.10+ | FastAPI | 8000 | `http://localhost:8000/doc.html` | 代码自动引用 `../../dist` |
| [litestar](./litestar/) | Python 3.10+ | LiteStar 2.0+ | 8000 | `http://localhost:8000/api/doc.html` | 代码自动引用 `../../dist` |
| [go-gin](./go-gin/) | Go 1.22+ | Gin | 8080 | `http://localhost:8080/doc.html` | 需手动复制到 `static/` |
| [go-stdlib](./go-stdlib/) | Go 1.18+ | 标准库（零依赖） | 8080 | `http://localhost:8080/doc.html` | 需手动复制到 `static/` |

## 前提条件

所有示例都需要先编译 Knife4j Vue3 前端产物：

```bash
# 在项目根目录执行
pnpm install
pnpm build
```

编译完成后 `dist/` 目录会生成 `doc.html`、`webjars/` 等文件。

## 一键配置（Windows）

提供了 `setup.bat` 脚本自动将前端产物复制到各示例项目：

```bash
cd examples
setup.bat
```

> **注意**：FastAPI 和 LiteStar 示例直接引用 `../../dist`，无需复制。`setup.bat` 仅复制到需要手动处理的示例。

## 各示例运行方式

### Java Spring Boot

```bash
cd examples/java-springboot                # 前端产物由 Maven 插件自动从 ../../dist 复制
mvn clean spring-boot:run                  # 编译 + 复制 + 启动
```

访问：http://localhost:8080/api/doc.html

### Python FastAPI

```bash
cd examples/fastapi
python -m venv venv                        # 创建虚拟环境（可选）
venv\Scripts\activate                      # Windows 激活
pip install -r requirements.txt            # 安装依赖
python main.py                             # 启动服务（端口 8000）
```

访问：http://localhost:8000/doc.html

> 文档页面需要 Basic Auth 认证：`admin/admin12345`

### Python LiteStar

```bash
cd examples/litestar
python -m venv venv                        # 创建虚拟环境（可选）
venv\Scripts\activate                      # Windows 激活
pip install -r requirements.txt            # 安装依赖
uvicorn main:app --host 0.0.0.0 --port 8000 --reload --root-path /api
```

访问：http://localhost:8000/api/doc.html

### Go Gin

```bash
cd examples/go-gin
cp -r ../../dist/* static/                 # 复制前端产物
go mod tidy                                # 下载依赖
go run main.go                             # 启动服务（端口 8080）
```

访问：http://localhost:8080/doc.html

### Go 标准库（零依赖）

```bash
cd examples/go-stdlib
cp -r ../../dist/* static/                 # 复制前端产物
go run main.go                             # 直接运行（零依赖，端口 8080）
```

访问：http://localhost:8080/doc.html

## 统一 API 接口

所有示例提供相同的 6 个 API 端点：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 根据 ID 获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/health | 健康检查 |

## 常见问题

1. **端口冲突**：Java Spring Boot 和 Go 示例默认均使用 8080 端口，同时运行前请先修改其中一个
2. **前端资源 404**：请确保已执行 `pnpm build`，且该示例的前端产物已就位（复制或用 Maven 插件）
3. **CORS 报错**：所有示例已配置 CORS 全开，如有问题请检查浏览器缓存
4. **Swagger Config**：每个示例都提供了 Knife4j 所需的 `/v3/api-docs/swagger-config` 端点
