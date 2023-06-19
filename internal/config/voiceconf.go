package config

// VoiceConf configuration
type VoiceConf struct {
	Lang        string `json:"lang"`
	Speaker     string `json:"speaker"`
	DownloadUrl string `json:"download_url"`
	License     string `json:"license"`
}
