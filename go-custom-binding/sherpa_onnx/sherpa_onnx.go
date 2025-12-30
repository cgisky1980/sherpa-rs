// Package sherpa_onnx provides Go bindings for sherpa-onnx C library
package sherpa_onnx

/*
#cgo windows LDFLAGS: -L${SRCDIR}/../crates/sherpa-rs-sys/sherpa-onnx/build/lib/Release -lsherpa-onnx-c-api -lsherpa-onnx-core -lsherpa-onnx-kaldifst-core -lsherpa-onnx-fst -lsherpa-onnx-fstfar -lsherpa-onnx-fstscript -lsherpa-onnx-kaldifst -lsherpa-onnx-kaldifstcore -lsherpa-onnx-kaldilm -lsherpa-onnx-kws -lsherpa-onnx-punctuation -lsherpa-onnx-speaker -lsherpa-onnx-tts -lsherpa-onnx-vad -lsherpa-onnx-vad-core -lsherpa-onnx-whisper -lsherpa-onnx
#cgo windows CFLAGS: -I${SRCDIR}/../crates/sherpa-rs-sys/sherpa-onnx/sherpa-onnx/c-api

#include <stdlib.h>
#include "c-api.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// OfflineTtsConfig 配置结构体
type OfflineTtsConfig struct {
	Model           *OfflineTtsModelConfig
	RuleFsts        string
	RuleFars        string
	MaxNumSentences int32
}

// OfflineTtsModelConfig 模型配置结构体
type OfflineTtsModelConfig struct {
	Zipvoice *ZipvoiceModelConfig
	NumThreads int32
	Debug      int32
	Provider   string
}

// ZipvoiceModelConfig ZipVoice模型配置
type ZipvoiceModelConfig struct {
	Tokens            string
	TextModel         string
	FlowMatchingModel string
	DataDir           string
	PinyinDict        string
	Vocoder           string
	FeatScale         float32
	TShift            float32
	TargetRms         float32
	GuidanceScale     float32
}

// GeneratedAudio 生成的音频数据
type GeneratedAudio struct {
	ptr *C.SherpaOnnxGeneratedAudio
}

// OfflineTts TTS引擎
type OfflineTts struct {
	ptr *C.SherpaOnnxOfflineTts
}

// NewOfflineTts 创建新的TTS引擎
func NewOfflineTts(config *OfflineTtsConfig) *OfflineTts {
	// 创建C配置结构体
	cConfig := C.SherpaOnnxOfflineTtsConfig{}
	
	// 设置模型配置
	if config.Model != nil {
		cModelConfig := C.SherpaOnnxOfflineTtsModelConfig{}
		
		// 设置ZipVoice配置
		if config.Model.Zipvoice != nil {
			cZipvoiceConfig := C.SherpaOnnxOfflineTtsZipvoiceModelConfig{}
			
			// 设置字符串字段
			if config.Model.Zipvoice.Tokens != "" {
				cZipvoiceConfig.tokens = C.CString(config.Model.Zipvoice.Tokens)
			}
			if config.Model.Zipvoice.TextModel != "" {
				cZipvoiceConfig.text_model = C.CString(config.Model.Zipvoice.TextModel)
			}
			if config.Model.Zipvoice.FlowMatchingModel != "" {
				cZipvoiceConfig.flow_matching_model = C.CString(config.Model.Zipvoice.FlowMatchingModel)
			}
			if config.Model.Zipvoice.DataDir != "" {
				cZipvoiceConfig.data_dir = C.CString(config.Model.Zipvoice.DataDir)
			}
			if config.Model.Zipvoice.PinyinDict != "" {
				cZipvoiceConfig.pinyin_dict = C.CString(config.Model.Zipvoice.PinyinDict)
			}
			if config.Model.Zipvoice.Vocoder != "" {
				cZipvoiceConfig.vocoder = C.CString(config.Model.Zipvoice.Vocoder)
			}
			
			// 设置浮点数字段
			cZipvoiceConfig.feat_scale = C.float(config.Model.Zipvoice.FeatScale)
			cZipvoiceConfig.t_shift = C.float(config.Model.Zipvoice.TShift)
			cZipvoiceConfig.target_rms = C.float(config.Model.Zipvoice.TargetRms)
			cZipvoiceConfig.guidance_scale = C.float(config.Model.Zipvoice.GuidanceScale)
			
			cModelConfig.zipvoice = cZipvoiceConfig
		}
		
		// 设置模型通用配置
		cModelConfig.num_threads = C.int32_t(config.Model.NumThreads)
		cModelConfig.debug = C.int32_t(config.Model.Debug)
		if config.Model.Provider != "" {
			cModelConfig.provider = C.CString(config.Model.Provider)
		}
		
		cConfig.model = cModelConfig
	}
	
	// 设置其他配置
	if config.RuleFsts != "" {
		cConfig.rule_fsts = C.CString(config.RuleFsts)
	}
	if config.RuleFars != "" {
		cConfig.rule_fars = C.CString(config.RuleFars)
	}
	cConfig.max_num_sentences = C.int32_t(config.MaxNumSentences)
	
	// 创建TTS引擎
	ptr := C.SherpaOnnxCreateOfflineTts(&cConfig)
	
	// 清理C字符串
	if config.Model != nil {
		if config.Model.Zipvoice != nil {
			if config.Model.Zipvoice.Tokens != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.tokens))
			}
			if config.Model.Zipvoice.TextModel != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.text_model))
			}
			if config.Model.Zipvoice.FlowMatchingModel != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.flow_matching_model))
			}
			if config.Model.Zipvoice.DataDir != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.data_dir))
			}
			if config.Model.Zipvoice.PinyinDict != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.pinyin_dict))
			}
			if config.Model.Zipvoice.Vocoder != "" {
				C.free(unsafe.Pointer(cConfig.model.zipvoice.vocoder))
			}
		}
		if config.Model.Provider != "" {
			C.free(unsafe.Pointer(cConfig.model.provider))
		}
	}
	if config.RuleFsts != "" {
		C.free(unsafe.Pointer(cConfig.rule_fsts))
	}
	if config.RuleFars != "" {
		C.free(unsafe.Pointer(cConfig.rule_fars))
	}
	
	return &OfflineTts{ptr: ptr}
}

// DeleteOfflineTts 销毁TTS引擎
func DeleteOfflineTts(tts *OfflineTts) {
	C.SherpaOnnxDestroyOfflineTts(tts.ptr)
}

// Generate 生成音频
func (tts *OfflineTts) Generate(text string, sid int32, speed float32) *GeneratedAudio {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	
	cAudio := C.SherpaOnnxOfflineTtsGenerate(tts.ptr, cText, C.int32_t(sid), C.float(speed))
	return &GeneratedAudio{ptr: cAudio}
}

// Samples 获取音频样本数据
func (audio *GeneratedAudio) Samples() []float32 {
	if audio.ptr == nil {
		return nil
	}
	return (*[1 << 30]float32)(unsafe.Pointer(audio.ptr.samples))[:audio.ptr.n:audio.ptr.n]
}

// SampleRate 获取音频采样率
func (audio *GeneratedAudio) SampleRate() int32 {
	if audio.ptr == nil {
		return 0
	}
	return int32(audio.ptr.sample_rate)
}

// Length 获取音频长度（样本数）
func (audio *GeneratedAudio) Length() int32 {
	if audio.ptr == nil {
		return 0
	}
	return int32(audio.ptr.n)
}

// WriteWave 将音频写入WAV文件
func (audio *GeneratedAudio) WriteWave(filename string) int32 {
	if audio.ptr == nil {
		return -1
	}
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	
	return int32(C.SherpaOnnxWriteWave(audio.ptr.samples, audio.ptr.n, audio.ptr.sample_rate, cFilename))
}

// ReadWave 从WAV文件读取音频
func ReadWave(filename string) (*GeneratedAudio, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	
	cWave := C.SherpaOnnxReadWave(cFilename)
	if cWave == nil {
		return nil, fmt.Errorf("failed to read wave file: %s", filename)
	}
	
	// 转换为GeneratedAudio格式
	generatedAudio := &C.SherpaOnnxGeneratedAudio{
		samples: cWave.samples,
		n: cWave.n,
		sample_rate: cWave.sample_rate,
	}
	
	// 销毁原始wave对象
	C.SherpaOnnxDestroyWave(cWave)
	
	return &GeneratedAudio{ptr: generatedAudio}, nil
}

// GenerateWithZipvoice 使用ZipVoice生成音频
func (tts *OfflineTts) GenerateWithZipvoice(text, promptText string, promptSamples []float32, promptSampleRate int32, speed float32, numSteps int32) *GeneratedAudio {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	
	cPromptText := C.CString(promptText)
	defer C.free(unsafe.Pointer(cPromptText))
	
	var cPromptSamples *C.float
	if len(promptSamples) > 0 {
		cPromptSamples = (*C.float)(unsafe.Pointer(&promptSamples[0]))
	}
	
	cAudio := C.SherpaOnnxOfflineTtsGenerateWithZipvoice(
		tts.ptr, 
		cText, 
		cPromptText, 
		cPromptSamples, 
		C.int32_t(len(promptSamples)), 
		C.int32_t(promptSampleRate), 
		C.float(speed), 
		C.int32_t(numSteps))
	
	return &GeneratedAudio{ptr: cAudio}
}

// Save 保存音频到文件
func (audio *GeneratedAudio) Save(filename string) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	
	result := C.SherpaOnnxWriteWave(audio.ptr.samples, audio.ptr.n, audio.ptr.sample_rate, cFilename)
	return result == 1
}

// ReadWave 读取WAV文件
func ReadWave(filename string) *Wave {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	
	cWave := C.SherpaOnnxReadWave(cFilename)
	return &Wave{ptr: cWave}
}

// Wave WAV文件结构体
type Wave struct {
	ptr *C.SherpaOnnxWave
}

// Samples 获取音频样本
func (w *Wave) Samples() []float32 {
	if w.ptr == nil || w.ptr.samples == nil {
		return nil
	}
	
	length := int(w.ptr.n)
	samples := make([]float32, length)
	
	cSamples := (*[1 << 30]C.float)(unsafe.Pointer(w.ptr.samples))[:length:length]
	for i := 0; i < length; i++ {
		samples[i] = float32(cSamples[i])
	}
	
	return samples
}

// SampleRate 获取采样率
func (w *Wave) SampleRate() int32 {
	if w.ptr == nil {
		return 0
	}
	return int32(w.ptr.sample_rate)
}

// DestroyWave 销毁WAV对象
func DestroyWave(wave *Wave) {
	C.SherpaOnnxDestroyWave(wave.ptr)
}