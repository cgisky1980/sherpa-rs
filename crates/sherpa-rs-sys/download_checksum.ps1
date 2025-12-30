$baseUrl = "https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.12.19"
$files = @(
    "sherpa-onnx-v1.12.19-win-x64-static.tar.bz2",
    "sherpa-onnx-v1.12.19-win-x64-shared.tar.bz2",
    "sherpa-onnx-v1.12.19-linux-x64-static.tar.bz2",
    "sherpa-onnx-v1.12.19-linux-x64-shared.tar.bz2",
    "sherpa-onnx-v1.12.19-linux-aarch64-static.tar.bz2",
    "sherpa-onnx-v1.12.19-linux-aarch64-shared-cpu.tar.bz2",
    "sherpa-onnx-v1.12.19-osx-universal2-static.tar.bz2",
    "sherpa-onnx-v1.12.19-osx-universal2-shared.tar.bz2",
    "sherpa-onnx-v1.12.19-android.tar.bz2",
    "sherpa-onnx-v1.12.19-ios.tar.bz2"
)

$outputDir = "C:\work\sherpa-rs\crates\sherpa-rs-sys\downloads"
New-Item -ItemType Directory -Force -Path $outputDir

foreach ($file in $files) {
    $url = "$baseUrl/$file"
    $outputPath = Join-Path $outputDir $file
    Write-Host "Downloading $file..."
    
    try {
        Invoke-WebRequest -Uri $url -OutFile $outputPath -TimeoutSec 300
        Write-Host "  Downloaded to $outputPath"
        
        $hash = Get-FileHash -Algorithm SHA256 -Path $outputPath
        Write-Host "  SHA256: $($hash.Hash)"
    }
    catch {
        Write-Host "  Error downloading $file: $_"
    }
}

Write-Host ""
Write-Host "Done. Files are in $outputDir"
