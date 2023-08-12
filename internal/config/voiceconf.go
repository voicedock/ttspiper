package config

// VoiceConf configuration
type VoiceConf struct {
	Lang                string `json:"lang"`
	Speaker             string `json:"speaker"`
	DownloadOnnxUrl     string `json:"download_onnx_url"`
	DownloadOnnxJsonUrl string `json:"download_onnx_json_url"`
	License             string `json:"license"`
}
