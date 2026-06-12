"""
数据模型模块
"""
from pydantic import BaseModel, Field
from typing import List, Optional


class UserCreate(BaseModel):
    """创建用户请求体"""
    name: str = Field(description="用户名", examples=["张三"])
    email: str = Field(description="邮箱地址", examples=["zhangsan@example.com"])
    role: str = Field(description="用户角色", default="user", examples=["admin"])


class User(UserCreate):
    """用户响应模型"""
    id: int = Field(description="用户ID", examples=[1])


class HealthResponse(BaseModel):
    """健康检查响应"""
    status: str = Field(description="服务状态", examples=["ok"])
    service: str = Field(description="服务名称", examples=["knife4j-vue3-litestar-example"])
