package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"sherpa-onnx-go-custom/sherpa_onnx"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 设置模型路径
	modelDir := "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia"
	
	config := sherpa.OfflineTtsConfig{
		Model: &sherpa.OfflineTtsModelConfig{
			Zipvoice: &sherpa.ZipvoiceModelConfig{
				Tokens:            filepath.Join(modelDir, "tokens.txt"),
				TextModel:         filepath.Join(modelDir, "text_encoder.onnx"),
				FlowMatchingModel: filepath.Join(modelDir, "fm_decoder.onnx"),
				DataDir:           filepath.Join(modelDir, "espeak-ng-data"),
				PinyinDict:        filepath.Join(modelDir, "pinyin.raw"),
				Vocoder:           filepath.Join(modelDir, "vocos_24khz.onnx"),
				FeatScale:         1.0,
				TShift:            1.0,
				TargetRms:         0.5,
				GuidanceScale:     1.0,
			},
			NumThreads: 1,
			Debug:      0,
			Provider:   "cpu",
		},
		MaxNumSentences: 1,
	}

	// 检查模型文件是否存在
	requiredFiles := []string{
		config.Model.Zipvoice.Tokens,
		config.Model.Zipvoice.TextModel,
		config.Model.Zipvoice.FlowMatchingModel,
		config.Model.Zipvoice.Vocoder,
		config.Model.Zipvoice.PinyinDict,
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("模型文件不存在: %s", file)
		}
	}

	// 检查数据目录
	if _, err := os.Stat(config.Model.Zipvoice.DataDir); os.IsNotExist(err) {
		log.Fatalf("数据目录不存在: %s", config.Model.Zipvoice.DataDir)
	}

	filename := "test-zipvoice-custom.wav"
	promptAudio := filepath.Join(modelDir, "prompt.wav")
	promptText := "周日被我射熄火了，所以今天是周一。"
	text := "我是中国人民的儿子，我爱我的祖国。我得祖国是一个伟大的国家，拥有五千年的文明史。"
	zipvoiceNumSteps := 4
	speed := float32(1.0)

	log.Println("输入文本:", text)
	log.Println("输出文件名:", filename)

	log.Println("初始化模型 (可能需要几秒钟)...")
	tts := sherpa.NewOfflineTts(&config)
	defer sherpa.DeleteOfflineTts(tts)
	log.Println("模型创建成功!")

	log.Println("开始生成音频...")
	var audio *sherpa.GeneratedAudio

	if _, err := os.Stat(promptAudio); err == nil {
		log.Println("使用提示音频进行零样本TTS...")
		
		// 读取提示音频
		wave := sherpa.ReadWave(promptAudio)
		defer sherpa.DestroyWave(wave)
		
		if wave == nil {
			log.Fatal("无法读取提示音频文件")
		}
		
		promptSamples := wave.Samples()
		promptSampleRate := wave.SampleRate()
		
		if len(promptSamples) == 0 {
			log.Fatal("提示音频文件为空")
		}
		
		audio = tts.GenerateWithZipvoice(
			text,
			promptText,
			promptSamples,
			promptSampleRate,
			speed,
			int32(zipvoiceNumSteps),
		)
	} else {
		log.Println("使用标准TTS模式...")
		audio = tts.Generate(text, 0, float32(math.Max(float64(speed), 1e-6)))
	}

	log.Println("音频生成完成!")
	
	if ok := audio.Save(filename); !ok {
		log.Fatalf("保存文件失败: %s", filename)
	}
	
	log.Println("音频已保存到:", filename)
}