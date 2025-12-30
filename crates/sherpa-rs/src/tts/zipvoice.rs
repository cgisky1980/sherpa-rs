use std::{mem, ptr::null};

use crate::{utils::cstring_from_str, OnnxConfig};
use eyre::{bail, Result};
use sherpa_rs_sys;

use super::{CommonTtsConfig, TtsAudio};

pub struct ZipVoiceTts {
    tts: *const sherpa_rs_sys::SherpaOnnxOfflineTts,
}

#[derive(Default)]
pub struct ZipVoiceTtsConfig {
    pub tokens: String,
    pub text_model: String,
    pub flow_matching_model: String,
    pub vocoder: String,
    pub data_dir: String,
    pub pinyin_dict: String,
    pub feat_scale: f32,
    pub t_shift: f32,
    pub target_rms: f32,
    pub guidance_scale: f32,

    pub onnx_config: OnnxConfig,
    pub tts_config: CommonTtsConfig,
}

impl ZipVoiceTts {
    pub fn new(config: ZipVoiceTtsConfig) -> Self {
        let tts = unsafe {
            let tokens = cstring_from_str(&config.tokens);
            let text_model = cstring_from_str(&config.text_model);
            let flow_matching_model = cstring_from_str(&config.flow_matching_model);
            let vocoder = cstring_from_str(&config.vocoder);
            let data_dir = cstring_from_str(&config.data_dir);
            let pinyin_dict = cstring_from_str(&config.pinyin_dict);

            let provider = cstring_from_str(&config.onnx_config.provider);

            let tts_config = config.tts_config.to_raw();

            let model_config = sherpa_rs_sys::SherpaOnnxOfflineTtsModelConfig {
                num_threads: config.onnx_config.num_threads,
                vits: mem::zeroed::<_>(),
                debug: config.onnx_config.debug.into(),
                provider: provider.as_ptr(),
                matcha: mem::zeroed::<_>(),
                kokoro: mem::zeroed::<_>(),
                kitten: mem::zeroed(),
                zipvoice: sherpa_rs_sys::SherpaOnnxOfflineTtsZipvoiceModelConfig {
                    tokens: tokens.as_ptr(),
                    text_model: text_model.as_ptr(),
                    flow_matching_model: flow_matching_model.as_ptr(),
                    vocoder: vocoder.as_ptr(),
                    data_dir: data_dir.as_ptr(),
                    pinyin_dict: pinyin_dict.as_ptr(),
                    feat_scale: config.feat_scale,
                    t_shift: config.t_shift,
                    target_rms: config.target_rms,
                    guidance_scale: config.guidance_scale,
                },
            };
            let config = sherpa_rs_sys::SherpaOnnxOfflineTtsConfig {
                max_num_sentences: config.tts_config.max_num_sentences,
                model: model_config,
                rule_fars: tts_config.rule_fars.map(|v| v.as_ptr()).unwrap_or(null()),
                rule_fsts: tts_config.rule_fsts.map(|v| v.as_ptr()).unwrap_or(null()),
                silence_scale: config.tts_config.silence_scale,
            };
            sherpa_rs_sys::SherpaOnnxCreateOfflineTts(&config)
        };

        Self { tts }
    }

    pub fn create(&mut self, text: &str, sid: i32, speed: f32) -> Result<TtsAudio> {
        unsafe { super::create(self.tts, text, sid, speed) }
    }

    pub fn create_with_zipvoice(
        &mut self,
        text: &str,
        prompt_text: &str,
        prompt_samples: &[f32],
        prompt_sample_rate: u32,
        speed: f32,
        num_steps: i32,
    ) -> Result<TtsAudio> {
        unsafe {
            let text = cstring_from_str(text);
            let prompt_text = cstring_from_str(prompt_text);

            let audio_ptr = sherpa_rs_sys::SherpaOnnxOfflineTtsGenerateWithZipvoice(
                self.tts,
                text.as_ptr(),
                prompt_text.as_ptr(),
                prompt_samples.as_ptr(),
                prompt_samples.len() as i32,
                prompt_sample_rate as i32,
                speed,
                num_steps,
            );

            if audio_ptr.is_null() {
                bail!("audio is null");
            }
            let audio = audio_ptr.read();

            if audio.n.is_negative() {
                bail!("no samples found");
            }
            if audio.samples.is_null() {
                bail!("audio samples are null");
            }
            let samples: &[f32] = std::slice::from_raw_parts(audio.samples, audio.n as usize);
            let samples = samples.to_vec();
            let sample_rate = audio.sample_rate;
            let duration = (samples.len() as i32) / sample_rate;

            // Free
            sherpa_rs_sys::SherpaOnnxDestroyOfflineTtsGeneratedAudio(audio_ptr);

            Ok(TtsAudio {
                samples,
                sample_rate: sample_rate as u32,
                duration,
            })
        }
    }
}

unsafe impl Send for ZipVoiceTts {}
unsafe impl Sync for ZipVoiceTts {}

impl Drop for ZipVoiceTts {
    fn drop(&mut self) {
        unsafe {
            sherpa_rs_sys::SherpaOnnxDestroyOfflineTts(self.tts);
        }
    }
}
