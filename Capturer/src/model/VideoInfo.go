package model

import (
	"crypto/md5"
	"fmt"
)

// VideoInfo VideoInfo
type VideoInfo struct {
	ID            int32  `json:"id"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	Thumb         string `json:"thumb"`
	CateID        string `json:"category_id"`
	M3U8Path      string `json:"path"`
	MajorCategory string `json:"-"`
	MinorCategory string `json:"-"`
	RetryCount    int32  `json:"-"`
}

// GetMD5 GetMD5
func (p *VideoInfo) GetMD5() string {
	value := fmt.Sprintf("tianlang.app|%d", p.ID)
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}
