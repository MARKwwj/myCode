package apptianlang

import (
	"TianlangCapturer/src/api"
	"fmt"
	"time"
)

// ReqCategory ReqCategory
type ReqCategory struct {
	ID   string `json:"category_id"`
	Time string `json:"time"`
}

// PullCategory PullCategory
func (d *Disguiser) PullCategory(cateID int32) (err error) {
	api.Info("[PullCategory]CateID:%d", cateID)
	var req ReqCategory
	req.ID = fmt.Sprintf("%d", cateID)
	req.Time = fmt.Sprintf("%d", time.Now().Unix())

	var result Result
	if result, err = d.doServer("PUT", "https://leiluohaoxiong.vip/api/v2/cate", &req); err != nil {
		return
	}

	// cateList := result.GetValue("cateList").([]interface{})
	// for _, value := range cateList {
	// 	category := value.(map[string]interface{})
	// 	id := int32(category["id"].(float64))
	// 	name := category["title"].(string)
	// 	d.SetCategoryName(id, name, cateID)
	// 	d.mVideoCategoryList = append(d.mVideoCategoryList, id)
	// }

	cateList := result.GetValue("cateList").AsSlice()
	for _, value := range cateList {
		id := value.Find("id").AsInt32()
		name := value.Find("title").AsString()
		d.SetCategoryName(id, name, cateID)
		d.mVideoCategoryList = append(d.mVideoCategoryList, id)
	}

	return
}
