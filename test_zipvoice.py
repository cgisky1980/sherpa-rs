#!/usr/bin/env python3
"""
简单的ZipVoice测试脚本
"""

import os
import sys

# 添加sherpa-onnx Python模块路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'crates', 'sherpa-rs-sys', 'sherpa-onnx', 'python'))

try:
    import sherpa_onnx
    print("成功导入sherpa_onnx模块")
    
    # 检查ZipVoice模型配置
    config = sherpa_onnx.OfflineTtsZipvoiceModelConfig()
    print("ZipVoice配置类可用")
    
    # 检查模型文件是否存在
    model_dir = "sherpa-onnx-zipvoice-distill-zh-en-emilia"
    required_files = [
        "text_encoder.onnx",
        "fm_decoder.onnx", 
        "vocos_24khz.onnx",
        "tokens.txt",
        "pinyin.raw"
    ]
    
    print("检查模型文件:")
    for file in required_files:
        file_path = os.path.join(model_dir, file)
        if os.path.exists(file_path):
            print(f"  ✓ {file}")
        else:
            print(f"  ✗ {file} - 文件不存在")
    
    # 检查espeak-ng-data目录
    data_dir = os.path.join(model_dir, "espeak-ng-data")
    if os.path.exists(data_dir):
        print("  ✓ espeak-ng-data目录")
    else:
        print("  ✗ espeak-ng-data目录不存在")
        
    print("\nZipVoice模型文件检查完成！")
    
except ImportError as e:
    print(f"导入sherpa_onnx失败: {e}")
    print("请确保sherpa-onnx Python模块已正确构建")
    
except Exception as e:
    print(f"测试过程中出现错误: {e}")