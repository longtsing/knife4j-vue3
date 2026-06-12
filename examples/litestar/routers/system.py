"""
系统路由
"""
from litestar import get
from models.user import HealthResponse


@get("/api/health", tags=["系统"], summary="健康检查")
async def health_check() -> HealthResponse:
    """健康检查接口"""
    return HealthResponse(status="ok", service="knife4j-vue3-litestar-example")
