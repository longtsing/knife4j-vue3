package com.example.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Map;

/**
 * 系统端点控制器
 */
@RestController
@Tag(name = "系统", description = "系统相关接口")
public class SystemController {

    private static final String START_TIME = LocalDateTime.now()
            .format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));

    @GetMapping("/api/version")
    @Operation(summary = "平台运行信息")
    public Map<String, String> version() {
        return Map.of(
                "title", "Spring Boot 后端示例项目",
                "description", "Spring Boot 后端示例项目 API 文档",
                "version", "1.0.0",
                "startTime", START_TIME,
                "Datetime", LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))
        );
    }

    @GetMapping("/api/health")
    @Operation(summary = "健康检查")
    public Map<String, String> healthCheck() {
        return Map.of(
                "status", "ok",
                "service", "java-springboot-backend"
        );
    }
}
