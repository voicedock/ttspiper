package config

type VoiceData struct {
	Voice     *Voice
	VoiceSpec *VoiceSpec
}

// Voice configuration
type Voice struct {
	Lang        string `json:"lang"`
	Speaker     string `json:"speaker"`
	DownloadUrl string `json:"download_url"`
	License     string `json:"license"`
}

// VoiceSpec voice specification
type VoiceSpec struct {
	Audio     VoiceSpecAudio     `json:"audio"`
	Espeak    VoiceSpecEspeak    `json:"espeak"`
	Inference VoiceSpecInference `json:"inference"`
}

type VoiceSpecAudio struct {
	SampleRate int `json:"sample_rate"`
}

type VoiceSpecEspeak struct {
	// Espeak voice language
	Voice string `json:"voice"`
}

type VoiceSpecInference struct {
	NoiseScale  float64 `json:"noise_scale"`
	LengthScale int     `json:"length_scale"`
	NoiseW      float64 `json:"noise_w"`
}
