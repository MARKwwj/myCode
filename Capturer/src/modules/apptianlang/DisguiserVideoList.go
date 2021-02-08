package apptianlang

import (
	"TianlangCapturer/src/api"
	"TianlangCapturer/src/model"
	"fmt"
	"strings"
	"time"
)

// ReqVideoList ReqVideoList
type ReqVideoList struct {
	Time       string `json:"time"`
	PageNum    string `json:"page"`
	LimitNum   string `json:"limit"`
	CategoryID string `json:"category_id"`
}

// PullVideoList PullVideoList
func (d *Disguiser) PullVideoList(categoryID int32, pageNum int32) (videoInfoList []*model.VideoInfo, err error) {
	api.Info("[PullVideoList]CateID:%d PageNum:%d", categoryID, pageNum)
	var req ReqVideoList
	req.Time = fmt.Sprintf("%d", time.Now().Unix())
	req.PageNum = fmt.Sprintf("%d", pageNum)
	req.LimitNum = "15"
	req.CategoryID = fmt.Sprintf("%d", categoryID)

	var result Result
	if result, err = d.doServer("PUT", "https://leiluohaoxiong.vip/api/v2/video/list", &req); err != nil {
		return nil, err
	}

	// videoList := result.GetValue("videoList").([]interface{})
	// for _, value := range videoList {
	// 	videoInfo := &model.VideoInfo{
	// 		MajorCategory: d.GetCategorySuperiorName(categoryID),
	// 		MinorCategory: d.GetCategoryName(categoryID),
	// 	}
	// 	model.UnmarshalJSON(value, videoInfo)
	// 	videoInfoList = append(videoInfoList, videoInfo)
	// }

	videoList := result.GetValue("videoList").AsRawSlice()
	for _, value := range videoList {
		videoInfo := &model.VideoInfo{
			MajorCategory: d.GetCategorySuperiorName(categoryID),
			MinorCategory: d.GetCategoryName(categoryID),
		}
		videoInfo.Title = strings.TrimSpace(videoInfo.Title)
		model.UnmarshalJSON(value, videoInfo)
		videoInfoList = append(videoInfoList, videoInfo)
	}

	return
}

// PullAllVideoList PullAllVideoList
func (d *Disguiser) PullAllVideoList(categoryID int32) (err error) {
	api.Info("[PullAllVideoList]CateID:%d", categoryID)
	var pageNum int32
	var videoInfoList []*model.VideoInfo
	for {
		pageNum++
		if videoInfoList, err = d.PullVideoList(categoryID, pageNum); err != nil {
			break
		}
		for _, videoInfo := range videoInfoList {
			AddDownloadTask(videoInfo)
		}
		if len(videoInfoList) < 15 {
			break
		}
		for len(mDownloadTask) >= 15 {
			api.Info("Wait for the download task, sleep(30s)")
			time.Sleep(time.Second * 30)
		}
	}
	return
}
