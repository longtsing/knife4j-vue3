from datetime import datetime
from typing import Optional
from fastapi import APIRouter, Query

router = APIRouter(tags=["测试"])


@router.get("/test", summary="测试接口")
async def test_endpoint(
    name: Optional[str] = Query(default="World", description="名称"),
    count: int = Query(default=1, description="数量")
):
    """测试接口，用于验证 API 文档功能"""
    return {
        "message": f"Hello {name}",
        "count": count,
        "timestamp": datetime.now().isoformat()
    }
