/**
 * LiteStar + Knife4j Vue3 集成配置
 * 为 LiteStar 框架优化的 Knife4j 配置
 */

// LiteStar 特定的配置常量
export const LITESTAR_CONFIG = {
  // API 端点配置
  endpoints: {
    swaggerConfig: '/v3/api-docs/swagger-config',
    swaggerConfigFallback: '/api/v3/api-docs/swagger-config',
    swaggerResources: '/swagger-resources',
    uiConfig: '/swagger-resources/configuration/ui',
    openApiSchema: '/schema/openapi.json',
    apiDocs: '/v3/api-docs'
  },
  
  // 框架信息
  framework: {
    name: 'LiteStar',
    version: '2.17.0',
    description: 'LiteStar ASGI Framework with OpenAPI 3.0 Support'
  },
  
  // UI 定制配置
  ui: {
    title: 'LiteStar API Documentation',
    description: '基于 LiteStar 框架的 API 文档，使用 Knife4j Vue3 界面',
    enableLiteStarFeatures: true,
    showFrameworkInfo: true,
    customFooter: 'Powered by LiteStar + Knife4j Vue3',
    // 禁用自动添加 basePath 功能
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

// LiteStar 特定的初始化选项
export const getLiteStarKnife4jOptions = () => {
  return {
    // 禁用 SpringDoc 模式以避免自动添加 basePath
    springdoc: false,
    
    // 语言设置
    i18n: 'zh-CN',
    
    // URL 配置 - 使用完整路径
    url: LITESTAR_CONFIG.endpoints.swaggerConfig,
    configUrl: LITESTAR_CONFIG.endpoints.uiConfig,
    
    // 启用配置支持
    configSupport: true,
    securitySupport: false,
    
    // 禁用 basePath 自动处理
    baseSpringFox: false,
    
    // LiteStar 特定配置
    framework: LITESTAR_CONFIG.framework.name,
    frameworkVersion: LITESTAR_CONFIG.framework.version,
    
    // UI 配置
    enableLiteStarFeatures: true,
    customTitle: LITESTAR_CONFIG.ui.title,
    customDescription: LITESTAR_CONFIG.ui.description,
    disableBasePath: true
  };
};

// 导出默认配置
export default LITESTAR_CONFIG;
