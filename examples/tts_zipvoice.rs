/*
使用说明：
1. 确保模型文件位于 `sherpa-onnx-zipvoice-distill-zh-en-emilia` 目录中
2. 运行示例：cargo run --example tts_zipvoice

模型文件说明：
- text_encoder.onnx: 文本编码器模型
- fm_decoder.onnx: 流匹配解码器模型  
- vocos_24khz.onnx: 声码器模型
- pinyin.raw: 拼音字典文件
- espeak-ng-data/: 语音数据目录
*/

use sherpa_rs::tts::{ZipVoiceTts, ZipVoiceTtsConfig};

fn main() {
    let config = ZipVoiceTtsConfig {
        tokens: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\tokens.txt".to_string(),
        text_model: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\text_encoder.onnx".to_string(),
        flow_matching_model: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\fm_decoder.onnx".to_string(),
        vocoder: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\vocos_24khz.onnx".to_string(),
        data_dir: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\espeak-ng-data".to_string(),
        pinyin_dict: "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\pinyin.raw".to_string(),
        feat_scale: 1.0,
        t_shift: 1.0,
        target_rms: 0.5,
        guidance_scale: 1.0,
        ..Default::default()
    };

    let mut tts = ZipVoiceTts::new(config);

    // 参考Go示例的参数设置
    let text = "我是中国人民的儿子，我爱我的祖国。我得祖国是一个伟大的国家，拥有五千年的文明史。";
    let prompt_text = "周日被我射熄火了，所以今天是周一。";
    
    // 读取真实的prompt.wav文件（参考Go示例的处理方式）
    let prompt_path = "c:\\work\\sherpa-rs\\sherpa-onnx-zipvoice-distill-zh-en-emilia\\prompt.wav";
    let (prompt_samples, prompt_sample_rate) = match sherpa_rs::read_audio_file(prompt_path) {
        Ok((samples, sample_rate)) => {
            println!("成功读取prompt.wav文件，采样率: {}Hz，样本数: {}", sample_rate, samples.len());
            (samples, sample_rate)
        }
        Err(e) => {
            println!("读取prompt.wav文件失败: {}", e);
            println!("使用示例数据替代");
            (vec![], 16000)
        }
    };
    
    let speed = 1.0;
    let num_steps = 4; // 参考Go示例设置为4步

    match tts.create_with_zipvoice(
        text,
        prompt_text,
        &prompt_samples,
        prompt_sample_rate,
        speed,
        num_steps,
    ) {
        Ok(audio) => {
            println!("生成音频成功！采样率: {}", audio.sample_rate);
            println!("音频样本数: {}", audio.samples.len());
            
            // 保存音频到文件
            let output_path = "output.wav";
            match sherpa_rs::write_audio_file(output_path, &audio.samples, audio.sample_rate) {
                Ok(_) => println!("音频已保存到: {}", output_path),
                Err(e) => println!("保存音频失败: {}", e),
            }
        }
        Err(e) => {
            println!("生成音频失败: {}", e);
        }
    }

    // ZipVoice特有功能示例（需要提示音频）
    /*
    let prompt_text = "这是一个提示文本";
    let prompt_samples = vec![0.0f32; 1000]; // 示例提示音频数据
    
    let audio = tts.create_with_zipvoice(
        "目标文本",
        prompt_text,
        &prompt_samples,
        22050,
        1.0,
        10,
    ).unwrap();
    
    sherpa_rs::write_audio_file("zipvoice_prompt_audio.wav", &audio.samples, audio.sample_rate).unwrap();
    println!("已创建 zipvoice_prompt_audio.wav");
    */
}