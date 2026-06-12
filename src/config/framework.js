/**
 * Knife4j Vue3 框架配置
 * 支持多种后端框架：FastAPI、LiteStar、Spring Boot、Go Gin 等
 */

// 框架配置常量
export const FRAMEWORK_CONFIG = {
  // API 端点配置（OpenAPI 3.0 风格）
  // 注意：所有路径已包含 /api 前缀，适应 root_path 配置
  endpoints: {
    swaggerConfig: '/api/v3/api-docs/swagger-config',
    swaggerConfigFallback: '/api/v3/api-docs/swagger-config',
    openApiSchema: '/api/openapi.json',
    apiDocs: '/api/v3/api-docs'
  },
  
  // 框架信息（可根据实际后端框架动态调整）
  framework: {
    name: 'Generic',
    version: '1.0.0',
    description: 'API Documentation with Knife4j Vue3'
  },
  
  // UI 定制配置
  ui: {
    title: 'API Documentation',
    description: '基于 Knife4j Vue3 的 API 文档界面',
    showFrameworkInfo: false,
    customFooter: 'Powered by Knife4j Vue3',
    // 禁用自动添加 basePath 功能（适用于现代框架如 FastAPI、LiteStar）
    disableBasePath: true,
    theme: {
      primaryColor: '#1890ff',
      headerBackgroundColor: '#001529',
      headerTextColor: '#ffffff'
    }
  },
  
  // 请求配置
  request: {
    timeout: 10000,
    retryCount: 2,
    retryDelay: 1000,
    enableFallback: true
  },
  
  // 调试配置
  debug: {
    enableConsoleLog: true,
    enableNetworkLog: true,
    enableErrorTracking: true
  }
};

// 获取 Knife4j 初始化选项
export const getKnife4jOptions = () => {
  // 读取部署配置的前缀
  const apiBasePath = (typeof window !== 'undefined' && window.KNIFE4J_CONFIG?.apiBasePath) || '';

  return {
    // 启用 SpringDoc 模式以支持 OpenAPI 3.0（FastAPI、LiteStar 等）
    springdoc: true,

    // 语言设置
    i18n: 'zh-CN',

    // URL 配置 - 使用相对路径或配置的前缀
    url: apiBasePath + '/v3/api-docs/swagger-config',
    configUrl: apiBasePath + '/v3/api-docs/swagger-config',

    // 启用配置支持
    configSupport: true,
    securitySupport: false,

    // 禁用 basePath 自动处理（现代框架使用 servers 配置）
    baseSpringFox: false,

    // 框架信息
    framework: FRAMEWORK_CONFIG.framework.name,
    frameworkVersion: FRAMEWORK_CONFIG.framework.version,

    // UI 配置
    customTitle: FRAMEWORK_CONFIG.ui.title,
    customDescription: FRAMEWORK_CONFIG.ui.description,
    disableBasePath: FRAMEWORK_CONFIG.ui.disableBasePath
  };
};

// 导出默认配置
export default FRAMEWORK_CONFIG;
