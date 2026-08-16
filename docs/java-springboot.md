# Java Spring Boot 对接指南

本文档介绍如何将 Knife4j Vue3 前端与 Java Spring Boot 后端集成。

## 方案对比

| 方案 | 优点 | 缺点 |
|------|------|------|
| **SpringDoc OpenAPI 3.0 + Knife4j Vue3**（本示例） | 自动生成 OpenAPI、注解贴近业务、生态成熟 | 需了解 springdoc 注解 |
| 旧版 Springfox 2.x | 老项目兼容 | 已停止维护，不再推荐 |

本示例使用 **SpringDoc OpenAPI 2.6.0** + **Knife4j Vue3 前端**（不依赖 knife4j-spring-boot-starter，因为前端已经是独立项目）。

## 前提条件

- Java 17+
- Maven 3.8+
- Spring Boot 3.2+
- SpringDoc OpenAPI 2.6.0
- 编译后的 Knife4j Vue3 前端

## 1. 项目结构

```
my-springboot-app/
├── pom.xml
├── src/main/java/com/example/
│   ├── Application.java
│   ├── config/
│   │   ├── CorsConfig.java
│   │   └── OpenApiConfig.java
│   └── controller/
│       ├── SwaggerConfigController.java
│       ├── UserController.java
│       └── SystemController.java
└── src/main/resources/
    ├── application.yml
    └── static/                     # Knife4j 前端产物（不纳入版本控制）
        ├── doc.html
        ├── webjars/
        └── oauth/
```

## 2. Maven 配置

参考 [examples/java-springboot/pom.xml](../examples/java-springboot/pom.xml)：

```xml
<dependencies>
    <!-- Spring Boot Web -->
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>

    <!-- SpringDoc OpenAPI 3.0（仅生成 JSON，不含 Swagger UI） -->
    <dependency>
        <groupId>org.springdoc</groupId>
        <artifactId>springdoc-openapi-starter-webmvc-api</artifactId>
        <version>2.6.0</version>
    </dependency>
</dependencies>

<build>
    <plugins>
        <!-- 将 Knife4j 前端从 src/main/resources/static 复制到 target/classes/static/api -->
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-resources-plugin</artifactId>
            <executions>
                <execution>
                    <id>copy-knife4j-dist</id>
                    <phase>generate-resources</phase>
                    <goals>
                        <goal>copy-resources</goal>
                    </goals>
                    <configuration>
                        <outputDirectory>${project.build.directory}/classes/static/api</outputDirectory>
                        <resources>
                            <resource>
                                <directory>src/main/resources/static</directory>
                                <includes>
                                    <include>**/*</include>
                                </includes>
                            </resource>
                        </resources>
                    </configuration>
                </execution>
            </executions>
        </plugin>

        <plugin>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-maven-plugin</artifactId>
        </plugin>
    </plugins>
</build>
```

## 3. application.yml

```yaml
server:
  port: 8080

springdoc:
  api-docs:
    path: /api/v3/api-docs          # OpenAPI JSON 路径
  swagger-ui:
    enabled: false                  # 禁用 SpringDoc 自带 Swagger UI（用 Knife4j 替代）
```

## 4. 核心 Java 类

### 4.1 入口类

```java
@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

### 4.2 CORS 配置

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

### 4.3 OpenAPI 元信息

```java
@Configuration
public class OpenApiConfig {
    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
            .info(new Info()
                .title("Spring Boot 后端示例项目")
                .description("Spring Boot 后端示例项目 API 文档")
                .version("1.0.0")
                .contact(new Contact().name("Admin").email("admin@example.com")))
            .servers(List.of(
                new Server().url("/api").description("API 服务")
            ));
    }
}
```

### 4.4 Knife4j swagger-config 端点

参考 [examples/java-springboot/src/main/java/com/example/controller/SwaggerConfigController.java](../examples/java-springboot/src/main/java/com/example/controller/SwaggerConfigController.java)：

```java
@RestController
public class SwaggerConfigController {

    @GetMapping("/api/v3/api-docs/swagger-config")
    public Map<String, Object> swaggerConfig() {
        return Map.of(
            "urls", List.of(
                Map.of("url", "/api/v3/api-docs", "name", "default")
            ),
            "configUrl", "/api/v3/api-docs/swagger-config",
            "validatorUrl", ""
        );
    }
}
```

### 4.5 业务 Controller

```java
@RestController
@RequestMapping("/api/users")
@Tag(name = "用户管理", description = "用户的增删改查接口")
public class UserController {

    @GetMapping
    @Operation(summary = "获取用户列表")
    public List<User> listUsers() { /* ... */ }

    @PostMapping
    @Operation(summary = "创建用户")
    public User createUser(@RequestBody User user) { /* ... */ }
}
```

## 5. 启动流程

```bash
# 1. 在 knife4j-vue3 根目录编译前端
cd /path/to/knife4j-vue3
pnpm install
pnpm build

# 2. 复制前端产物到 Java 项目的 src/main/resources/static/
cp -r dist/* my-springboot-app/src/main/resources/static/

# 3. 启动
cd my-springboot-app
mvn clean spring-boot:run
```

访问：http://localhost:8080/api/doc.html

## 6. Knife4j 集成原理

```
浏览器访问 /api/doc.html
  ↓
doc.html 自带脚本检测路径前缀 → apiBasePath = '/api'
  ↓
请求 /api/v3/api-docs/swagger-config → 拿到 OpenAPI JSON 地址 /api/v3/api-docs
  ↓
请求 /api/v3/api-docs → SpringDoc 生成的 OpenAPI 规范
  ↓
Knife4j 渲染文档界面
```

调试时，Knife4j 内部 ajax 会自动带上 `/api` 前缀（通过 `apiBasePath`）。

> **2026-08 修复**：Knife4j 前端修复了**重复前缀 bug**。即使 swagger-config 返回 `/api/v3/api-docs`（含前缀）也不会再被拼接成 `/api/api/v3/api-docs`。

## 7. 关于 context-path

如果你的部署架构是：

```
nginx /api/*  →  Spring Boot (server.servlet.context-path=/api)
```

那么所有 `@RequestMapping("/users")` 的端点会暴露在 `/api/users`。OpenAPI JSON 中的 `servers[0].url` 也要设为 `/api`。

```yaml
server:
  port: 8080
  servlet:
    context-path: /api
```

**重要**：本示例不依赖 `context-path`。所有 API 路径都直接写成 `/api/xxx`，方便维护。

## 8. 高级配置

### 8.1 认证

#### Basic Auth（保护文档页面）

```java
@Component
public class AuthInterceptor implements HandlerInterceptor {
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        // 仅保护文档相关路径
        String path = request.getRequestURI();
        if (!path.contains("/doc.html") && !path.contains("/v3/api-docs")) {
            return true;
        }

        String auth = request.getHeader("Authorization");
        if (auth != null && auth.startsWith("Basic ")) {
            String decoded = new String(Base64.getDecoder().decode(auth.substring(6)));
            String[] parts = decoded.split(":", 2);
            if ("admin".equals(parts[0]) && "admin12345".equals(parts[1])) {
                return true;
            }
        }

        response.setHeader("WWW-Authenticate", "BASIC realm=\"API Documentation\"");
        response.setStatus(401);
        return false;
    }
}
```

#### Bearer Auth（保护业务 API）

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/**").authenticated()
                .anyRequest().permitAll()
            )
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()));
        return http.build();
    }
}
```

然后在 OpenAPI 中声明：

```java
@Bean
public OpenAPI customOpenAPI() {
    return new OpenAPI()
        .components(new Components()
            .addSecuritySchemes("bearerAuth",
                new SecurityScheme()
                    .type(SecurityScheme.Type.HTTP)
                    .scheme("bearer")
                    .bearerFormat("JWT")))
        .info(new Info().title("My API").version("1.0.0"));
}
```

在 Knife4j 界面中，通过「文档管理」→「全局参数设置」添加 Authorization Header（值 `Bearer <token>`）。

### 8.2 分组 API 文档

```yaml
springdoc:
  api-docs:
    groups:
      - group: users
        paths-to-match: /api/users/**
      - group: orders
        paths-to-match: /api/orders/**
```

swagger-config 返回多组：

```java
@GetMapping("/api/v3/api-docs/swagger-config")
public Map<String, Object> swaggerConfig() {
    return Map.of(
        "urls", List.of(
            Map.of("url", "/api/v3/api-docs?group=users", "name", "用户管理"),
            Map.of("url", "/api/v3/api-docs?group=orders", "name", "订单管理")
        ),
        "configUrl", "/api/v3/api-docs/swagger-config",
        "validatorUrl", ""
    );
}
```

### 8.3 Extensions 扩展

```java
@Tag(name = "用户管理", extensions = {
    @Extension(name = "extensions", properties = {
        @ExtensionProperty(name = "x-author", value = "张三"),
        @ExtensionProperty(name = "x-order", value = "1000")
    })
})
```

## 9. 生产环境部署

### 9.1 嵌入式（推荐）

编译后前端自动打到 JAR 包的 `BOOT-INF/classes/static/api/` 下。运行时访问 `/api/doc.html`。

```bash
mvn clean package
java -jar target/my-app-1.0.0.jar
```

### 9.2 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name docs.example.com;

    # Knife4j 前端
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

    # OpenAPI 端点
    location /v3/api-docs/ {
        proxy_pass http://127.0.0.1:8080/v3/api-docs/;
    }
}
```

### 9.3 Docker 一体化

```dockerfile
# 构建阶段
FROM node:20-alpine AS frontend
WORKDIR /app
COPY knife4j-vue3/package.json knife4j-vue3/pnpm-lock.yaml ./
RUN corepack enable && pnpm install
COPY knife4j-vue3/ ./
RUN pnpm build

# 构建后端
FROM maven:3.9-eclipse-temurin-17 AS backend
WORKDIR /app
COPY my-springboot-app/pom.xml .
RUN mvn dependency:go-offline
COPY my-springboot-app/src ./src
COPY --from=frontend /app/dist ./src/main/resources/static
RUN mvn clean package -DskipTests

# 运行
FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
COPY --from=backend /app/target/*.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

### 9.4 验证部署

```bash
curl -I http://your-server:8080/api/doc.html
curl http://your-server:8080/api/v3/api-docs/swagger-config
curl http://your-server:8080/api/users
```

## 10. 完整示例

[examples/java-springboot/](../examples/java-springboot/) 提供开箱即用的可运行示例：

```bash
cd examples/java-springboot
# 前置：根目录 pnpm build 生成 dist/
cp -r ../dist/* src/main/resources/static/
mvn clean spring-boot:run
# 访问 http://localhost:8080/api/doc.html
```

## 11. 常见问题

### Q1：调试面板响应区空白？

A：检查浏览器 Network：

- `/api/v3/api-docs/swagger-config` 应返回 200 + JSON
- `/api/v3/api-docs` 应返回 200 + OpenAPI JSON
- 调试请求的 URL 应为 `/api/xxx`，**不应**是 `/api/api/xxx`

Knife4j 前端已在 2026-08 修复重复前缀 bug，swagger-config 返回的 `/api/v3/api-docs`（含前缀）不会再被拼接。

### Q2：Knife4j 显示 "No API definitions found"？

A：检查 swagger-config 端点：

```bash
curl http://localhost:8080/api/v3/api-docs/swagger-config
```

应返回 `urls[0].url = /api/v3/api-docs`。

### Q3：端口冲突？

A：Spring Boot 默认 8080；Python 用 8000；Go 用 8080。修改 `application.yml` 的 `server.port` 或 `mvn spring-boot:run -Dspring-boot.run.arguments=--server.port=9090`。

### Q4：Maven 依赖下载失败？

A：检查 `pom.xml` 中的 repositories 配置。生产环境建议配置企业内网 Nexus 仓库。

### Q5：前端资源 404？

A：确保执行了 `pnpm build` 并把 `dist/*` 复制到 `src/main/resources/static/`。Maven 的 `generate-resources` 阶段会自动从这个目录复制到 `target/classes/static/api/`。

### Q6：POM 修改后未生效？

A：执行 `mvn clean` 后重新 `mvn spring-boot:run`，确保 `target/classes/static/api/` 是最新的。