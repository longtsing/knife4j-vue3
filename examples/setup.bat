@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   Knife4j Vue3 示例项目一键配置
echo ========================================
echo.

:: 检查前端编译产物
set "DIST_DIR=%~dp0..\dist"
if not exist "%DIST_DIR%\doc.html" (
    echo [!] 未找到前端编译产物: %DIST_DIR%
    echo.
    echo 请先编译前端项目：
    echo   cd %~dp0..
    echo   pnpm install
    echo   pnpm build
    echo.
    pause
    exit /b 1
)

echo [√] 找到前端编译产物: %DIST_DIR%
echo.

:: 复制到 Java Spring Boot 示例（Maven 插件会自动处理，但复制一份方便 IDE 直接运行）
echo [1/5] 准备 java-springboot（Maven 插件启动时自动复制）
echo       [√] 跳过（Maven resources 插件会自动复制到 target/classes/static/api/）
echo.

:: 复制到 FastAPI 示例（代码直接引用 ../../dist，无需复制）
echo [2/5] 准备 fastapi（代码自动读取 ../../dist）
echo       [√] 跳过（FastAPI 直接引用 dist 目录）
echo.

:: 复制到 LiteStar 示例（代码直接引用 ../../dist，无需复制）
echo [3/5] 准备 litestar（代码自动读取 ../../dist）
echo       [√] 跳过（LiteStar 直接引用 dist 目录）
echo.

:: 复制到 Go Gin 示例
echo [4/5] 复制到 go-gin\static\
if not exist "%~dp0go-gin\static" mkdir "%~dp0go-gin\static"
xcopy /E /I /Y "%DIST_DIR%\*" "%~dp0go-gin\static\" >nul
echo       [√] 完成
echo.

:: 复制到 Go 标准库示例
echo [5/5] 复制到 go-stdlib\static\
if not exist "%~dp0go-stdlib\static" mkdir "%~dp0go-stdlib\static"
xcopy /E /I /Y "%DIST_DIR%\*" "%~dp0go-stdlib\static\" >nul
echo       [√] 完成
echo.

echo ========================================
echo   配置完成！
echo ========================================
echo.
echo ┌─────────────────────────────────────────────────────────────┐
echo │  启动方式                                                   │
echo ├─────────────────────────────────────────────────────────────┤
echo │                                                             │
echo │  Java Spring Boot (端口 8080):                              │
echo │    cd java-springboot                                       │
echo │    mvn clean spring-boot:run                                │
echo │    访问 http://localhost:8080/api/doc.html                  │
echo │                                                             │
echo │  Python FastAPI (端口 8000):                                │
echo │    cd fastapi                                               │
echo │    pip install -r requirements.txt                          │
echo │    python main.py                                           │
echo │    访问 http://localhost:8000/doc.html                      │
echo │    （认证: hxgis/hxgis12345 或 hbxqx/hbxqx168）            │
echo │                                                             │
echo │  Python LiteStar (端口 8000):                               │
echo │    cd litestar                                              │
echo │    pip install -r requirements.txt                          │
echo │    uvicorn main:app --root-path /api --reload               │
echo │    访问 http://localhost:8000/api/doc.html                  │
echo │                                                             │
echo │  Go Gin (端口 8080):                                        │
echo │    cd go-gin                                                │
echo │    go mod tidy                                              │
echo │    go run main.go                                           │
echo │    访问 http://localhost:8080/doc.html                      │
echo │                                                             │
echo │  Go 标准库 (端口 8080, 零依赖):                             │
echo │    cd go-stdlib                                             │
echo │    go run main.go                                           │
echo │    访问 http://localhost:8080/doc.html                      │
echo │                                                             │
echo └─────────────────────────────────────────────────────────────┘
echo.
pause
