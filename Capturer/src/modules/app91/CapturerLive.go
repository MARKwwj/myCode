package app91

import (
	"TianlangCapturer/src/api"
)

// Capture91LiveModule Capture91LiveModule
type Capture91LiveModule struct {
	mAESKey  string
	mSignKey string
	mLiveKey string
}

// Name Name
func (p *Capture91LiveModule) Name() string { return "91.live" }

// Run Run
func (p *Capture91LiveModule) Run() error {
	return nil
}

func init() {
	mCapture91Live := new(Capture91LiveModule)
	mCapture91Live.mAESKey = "e79465cfbb39cjdusimcuekd3b066a6e"
	mCapture91Live.mSignKey = "132f1537f85sjdpcm59f7e318b9epa51"
	mCapture91Live.mLiveKey = "ljhlksslgkjfhlksuo8472rju6p2od03"
	api.RegisterCapturer(mCapture91Live)
}
