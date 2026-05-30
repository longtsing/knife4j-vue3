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

echo [✓] 找到前端编译产物: %DIST_DIR%
echo.

:: 复制到 Java Spring Boot 示例
echo [1/3] 复制到 java-springboot\src\main\resources\static\
if not exist "%~dp0java-springboot\src\main\resources\static" mkdir "%~dp0java-springboot\src\main\resources\static"
xcopy /E /I /Y "%DIST_DIR%\*" "%~dp0java-springboot\src\main\resources\static\" >nul
echo       [✓] 完成
echo.

:: 复制到 FastAPI 示例
echo [2/3] 复制到 fastapi\static\
if not exist "%~dp0fastapi\static" mkdir "%~dp0fastapi\static"
xcopy /E /I /Y "%DIST_DIR%\*" "%~dp0fastapi\static\" >nul
echo       [✓] 完成
echo.

:: 复制到 LiteStar 示例
echo [3/3] 复制到 litestar\static\
if not exist "%~dp0litestar\static" mkdir "%~dp0litestar\static"
xcopy /E /I /Y "%DIST_DIR%\*" "%~dp0litestar\static\" >nul
echo       [✓] 完成
echo.

echo ========================================
echo   配置完成！
echo ========================================
echo.
echo 启动方式：
echo.
echo   Java Spring Boot:
echo     cd java-springboot
echo     mvn spring-boot:run
echo     访问 http://localhost:8080/doc.html
echo.
echo   FastAPI:
echo     cd fastapi
echo     pip install -r requirements.txt
echo     uvicorn main:app --host 0.0.0.0 --port 8000 --reload
echo     访问 http://localhost:8000/doc.html
echo.
echo   LiteStar:
echo     cd litestar
echo     pip install -r requirements.txt
echo     uvicorn main:app --host 0.0.0.0 --port 8000 --reload
echo     访问 http://localhost:8000/doc.html
echo.
pause
