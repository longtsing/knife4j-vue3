package com.example.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Map;

/**
 * Knife4j Swagger 配置端点
 * 提供 Knife4j Vue3 前端所需的 swagger-config 接口
 */
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
