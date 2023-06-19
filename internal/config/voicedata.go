package config

type VoiceData struct {
	VoiceConf    *VoiceConf
	VoiceSpec    *VoiceSpec
	OnnxPath     string
	OnnxJsonPath string
}

func (v *VoiceData) Downloaded() bool {
	return v.VoiceSpec != nil
}

func (v *VoiceData) Downloadable() bool {
	return v.VoiceConf.DownloadUrl != ""
}
