# Java Spring Boot 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Java Spring Boot 后端集成。

## 方式一：SpringDoc OpenAPI 3.0（推荐）

适用于 Spring Boot 3.x + SpringDoc OpenAPI。

### 1. 添加依赖

```xml
<!-- pom.xml -->
<dependencies>
    <!-- SpringDoc OpenAPI -->
    <dependency>
        <groupId>org.springdoc</groupId>
        <artifactId>springdoc-openapi-starter-webmvc-ui</artifactId>
        <version>2.6.0</version>
    </dependency>
    
    <!-- Knife4j 增强（可选） -->
    <dependency>
        <groupId>com.github.xiaoymin</groupId>
        <artifactId>knife4j-openapi3-jakarta-spring-boot-starter</artifactId>
        <version>4.5.0</version>
    </dependency>
</dependencies>
```

### 2. 配置 application.yml

```yaml
springdoc:
  api-docs:
    path: /v3/api-docs
    groups:
      - group: default
        paths-to-match: /api/**
  swagger-ui:
    path: /doc.html
    tags-sorter: alpha
    operations-sorter: alpha

# Knife4j 配置（可选）
knife4j:
  enable: true
  setting:
    language: zh_cn
```

### 3. 配置 CORS

```java
@Configuration
public class CorsConfig implements WebMvcConfigurer {
    
    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**")
            .allowedOriginPatterns("*")
            .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
            .allowedHeaders("*")
            .allowCredentials(true);
    }
}
```

### 4. 前端代理配置

在 `vite.config.js` 中配置代理：

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/v3/api-docs': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

### 5. 启动服务

```bash
# 后端
mvn spring-boot:run

# 前端
pnpm dev
```

访问 `http://localhost:5173/doc.html`

---

## 方式二：Springfox 2.x（旧版）

适用于 Spring Boot 2.x + Springfox。

### 1. 添加依赖

```xml
<dependency>
    <groupId>io.springfox</groupId>
    <artifactId>springfox-boot-starter</artifactId>
    <version>3.0.0</version>
</dependency>
```

### 2. 配置

```java
@Configuration
@EnableOpenApi
public class SwaggerConfig {
    
    @Bean
    public Docket createRestApi() {
        return new Docket(DocumentationType.OAS_30)
            .apiInfo(apiInfo())
            .select()
            .apis(RequestHandlerSelectors.basePackage("com.example.controller"))
            .paths(PathSelectors.any())
            .build();
    }
    
    private ApiInfo apiInfo() {
        return new ApiInfoBuilder()
            .title("API Documentation")
            .version("1.0")
            .build();
    }
}
```

### 3. 前端配置

Springfox 使用不同的端点：

```javascript
// vite.config.js
proxy: {
  '/swagger-resources': {
    target: 'http://localhost:8080',
    changeOrigin: true
  },
  '/v2/api-docs': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

---

## 完整示例项目

### 目录结构

```
my-springboot-app/
├── src/main/java/com/example/
│   ├── Application.java
│   ├── config/
│   │   └── CorsConfig.java
│   └── controller/
│       └── UserController.java
├── src/main/resources/
│   └── application.yml
└── pom.xml
```

### 示例代码

```java
@RestController
@RequestMapping("/api/users")
@Tag(name = "用户管理", description = "用户的增删改查")
public class UserController {
    
    @GetMapping
    @Operation(summary = "获取用户列表")
    public List<User> list() {
        return userService.findAll();
    }
    
    @PostMapping
    @Operation(summary = "创建用户")
    public User create(@RequestBody User user) {
        return userService.save(user);
    }
}
```

### Extensions 扩展

Knife4j 支持通过 OpenAPI Extensions 为接口添加自定义元数据：

```java
@Tag(name = "用户管理", description = "用户的增删改查", extensions = {
    @Extension(name = "extensions", properties = {
        @ExtensionProperty(name = "x-author", value = "张三"),
        @ExtensionProperty(name = "x-order", value = "1000")
    })
})
```

这些扩展信息会在 Knife4j 界面中显示。

---

## 生产环境编译与部署

### 1. 编译前端项目

```bash
# 进入前端项目目录
cd knife4j-vue3

# 安装依赖
pnpm install

# 编译生产版本
pnpm build
```

编译产物在 `dist/` 目录下，包含：
- `doc.html` — Knife4j 入口页面
- `webjars/` — JS/CSS 静态资源

### 2. 部署方式一：前端资源嵌入 Spring Boot（推荐）

将编译产物复制到 Spring Boot 的 `static` 目录，随 JAR 一起打包：

```bash
# 复制前端产物到后端静态资源目录
cp -r dist/* src/main/resources/static/
```

目录结构：
```
src/main/resources/static/
├── doc.html
└── webjars/
    ├── js/
    └── css/
```

后端无需额外配置，Spring Boot 默认提供 `static/` 目录下的静态资源。

启动后访问：`http://your-server:8080/doc.html`

### 3. 部署方式二：Nginx 反向代理

#### Nginx 配置

```nginx
server {
    listen 80;
    server_name docs.example.com;

    # Knife4j 前端静态资源
    location / {
        root /var/www/knife4j-vue3/dist;
        index doc.html;
        try_files $uri $uri/ /doc.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # OpenAPI 规范代理
    location /v3/api-docs/ {
        proxy_pass http://127.0.0.1:8080/v3/api-docs/;
        proxy_set_header Host $host;
    }
}
```

#### 部署步骤

```bash
# 1. 编译前端
pnpm build

# 2. 复制到 Nginx 目录
sudo cp -r dist/* /var/www/knife4j-vue3/

# 3. 测试并重载 Nginx
sudo nginx -t
sudo systemctl reload nginx
```

### 4. 部署方式三：Docker 部署

#### Dockerfile（前端 + 后端一体化）

```dockerfile
# 构建阶段
FROM node:20-alpine AS frontend
WORKDIR /app
COPY knife4j-vue3/package.json knife4j-vue3/pnpm-lock.yaml ./
RUN corepack enable && pnpm install
COPY knife4j-vue3/ ./
RUN pnpm build

# 运行阶段
FROM eclipse-temurin:21-jre-alpine AS backend
WORKDIR /app
COPY --from=frontend /app/dist ./static/
COPY springboot-app/target/*.jar ./app.jar

EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

#### 构建与运行

```bash
# 构建镜像
docker build -t my-knife4j-app .

# 运行容器
docker run -d -p 8080:8080 --name knife4j my-knife4j-app
```

### 5. 配置文件参考

#### application.yml（生产环境）

```yaml
server:
  port: 8080
  servlet:
    context-path: /api

springdoc:
  api-docs:
    path: /v3/api-docs
  swagger-ui:
    path: /doc.html

knife4j:
  enable: true
  setting:
    language: zh_cn
    enable_swagger_models: true
    enable_request_cache: true
    enable_host: false
```

### 6. 验证部署

```bash
# 检查前端页面
curl -I http://your-server:8080/doc.html

# 检查 OpenAPI 端点
curl http://your-server:8080/v3/api-docs/swagger-config

# 检查 API 是否正常
curl http://your-server:8080/api/users
```

---

## 常见问题

### Q: 页面空白，无法加载文档？

**A:** 检查以下几点：
1. 后端是否正常启动
2. `vite.config.js` 中的代理配置是否正确
3. 浏览器控制台是否有 CORS 错误

### Q: 显示 "No API definitions found"？

**A:** 确认后端的 OpenAPI 端点可访问：
```bash
curl http://localhost:8080/v3/api-docs/swagger-config
```

### Q: 如何添加认证头？

**A:** 在 Knife4j 界面的「文档管理」→「全局参数设置」中添加 Header 参数。
