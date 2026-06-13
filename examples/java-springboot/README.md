# Java Spring Boot 后端示例项目

基于 Spring Boot 3.2 + Knife4j Vue3 的后端 API 示例项目。

## 快速开始

```bash
# 1. 先编译前端产物（在项目根目录）
cd ../..
pnpm install
pnpm build

# 2. 启动服务
cd examples/java-springboot
mvn clean spring-boot:run
```

`mvn spring-boot:run` 会自动：
- 将 `../../dist` 的前端产物复制到 `target/classes/static/api/`
- 编译 Java 源码
- 启动 Spring Boot 服务

服务启动后访问：

- **Knife4j 文档**：http://localhost:8080/api/doc.html
- **OpenAPI JSON**：http://localhost:8080/api/v3/api-docs
- **Swagger Config**：http://localhost:8080/api/v3/api-docs/swagger-config

## 项目结构

```
java-springboot/
├── pom.xml                              # Maven 配置（含自动复制 dist 的插件）
├── src/main/java/com/example/
│   ├── Knife4jExampleApplication.java   # Spring Boot 入口
│   ├── config/
│   │   ├── CorsConfig.java              # CORS 跨域配置
│   │   └── OpenApiConfig.java           # OpenAPI 元信息
│   ├── controller/
│   │   ├── SwaggerConfigController.java # Knife4j swagger-config 端点
│   │   ├── SystemController.java        # /api/version, /api/health
│   │   └── UserController.java          # /api/users CRUD
│   └── model/
│       └── User.java                    # 用户实体（Lombok）
├── src/main/resources/
│   └── application.yml                  # Spring Boot 配置
└── README.md
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/{id} | 根据 ID 获取用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/{id} | 更新用户 |
| DELETE | /api/users/{id} | 删除用户 |
| GET | /api/version | 版本信息 |
| GET | /api/health | 健康检查 |

## 技术栈

- **Java 17**
- **Spring Boot 3.2.5** - Web 框架
- **SpringDoc OpenAPI 2.6.0** - 自动生成 OpenAPI 3.0 JSON
- **Maven** - 构建工具
- **Knife4j Vue3** - API 文档 UI（前端产物由 Maven 插件自动复制）

## 常见问题

### 端口冲突

默认端口为 **8080**，如果和其他示例冲突，修改 `application.yml`：

```yaml
server:
  port: 8090  # 改为其他端口
```

### Maven 依赖下载失败

如果自定义仓库不可用，`pom.xml` 中已添加 Maven Central 作为备用仓库。

### 前端资源 404

确保已执行 `pnpm build`，且 `../../dist` 目录存在。Maven 的 `generate-resources` 阶段会自动复制文件。
