# 从源代码构建sherpa-onnx C库的PowerShell脚本

# 设置环境变量
$env:SHERPA_BUILD_SHARED_LIBS = "1"

# 切换到sherpa-onnx目录
cd "c:\work\sherpa-rs\crates\sherpa-rs-sys\sherpa-onnx"

# 创建构建目录
if (Test-Path "build") {
    Remove-Item -Recurse -Force "build"
}
New-Item -ItemType Directory -Path "build" | Out-Null
cd "build"

# 配置CMake
Write-Host "配置CMake..."
cmake .. -DCMAKE_BUILD_TYPE=Release ^
         -DSHERPA_ONNX_ENABLE_C_API=ON ^
         -DSHERPA_ONNX_ENABLE_TTS=ON ^
         -DSHERPA_ONNX_ENABLE_BINARY=ON ^
         -DBUILD_SHARED_LIBS=ON ^
         -DSHERPA_ONNX_BUILD_C_API_EXAMPLES=ON

if ($LASTEXITCODE -ne 0) {
    Write-Error "CMake配置失败"
    exit 1
}

# 构建库
Write-Host "构建sherpa-onnx库..."
cmake --build . --config Release --parallel

if ($LASTEXITCODE -ne 0) {
    Write-Error "构建失败"
    exit 1
}

# 设置库路径环境变量
$sherpa_lib_path = "c:\work\sherpa-rs\crates\sherpa-rs-sys\sherpa-onnx\build\lib\Release"
$env:SHERPA_LIB_PATH = $sherpa_lib_path

Write-Host "构建完成！库文件位于: $sherpa_lib_path"
Write-Host "请设置环境变量: set SHERPA_LIB_PATH=$sherpa_lib_path"