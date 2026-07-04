"""
系统路由
"""
from litestar import get, post, Request
from litestar.datastructures import UploadFile
from litestar.params import Body
from models.user import HealthResponse
from typing import Dict, Any


@get("/health", tags=["系统"], summary="健康检查")
async def health_check() -> HealthResponse:
    """健康检查接口"""
    return HealthResponse(status="ok", service="knife4j-vue3-litestar-example")


@get("/echo/get", tags=["系统"], summary="GET请求回显")
async def echo_get(request: Request) -> Dict[str, Any]:
    """
    GET请求回显接口
    
    显示收到的请求头和查询参数
    """
    return {
        "headers": dict(request.headers),
        "query_params": dict(request.query_params),
    }


@post("/echo/post", tags=["系统"], summary="POST请求回显")
async def echo_post(
    request: Request,
    data: Dict[str, Any] = Body(description="POST请求体")
) -> Dict[str, Any]:
    """
    POST请求回显接口
    
    显示收到的请求头和POST请求内容
    """
    return {
        "headers": dict(request.headers),
        "body": data
    }
