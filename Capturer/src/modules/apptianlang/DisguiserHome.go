package apptianlang

import "TianlangCapturer/src/api"

// PullHome PullHome
func (d *Disguiser) PullHome() (err error) {
	api.Info("[PullHome]")
	var result Result
	if result, err = d.doServer("POST", "https://leiluohaoxiong.vip/api/v2/home", nil); err != nil {
		return
	}

	// categoryList := result.GetValue("cate").([]interface{})
	// for _, value := range categoryList {
	// 	category := value.(map[string]interface{})
	// 	id := int32(category["id"].(float64))
	// 	name := category["title"].(string)
	// 	d.SetCategoryName(id, name, 0)
	// 	if err = d.PullCategory(id); err != nil {
	// 		return
	// 	}
	// }

	categoryList := result.GetValue("cate").AsSlice()
	for _, value := range categoryList {
		id := value.Find("id").AsInt32()
		name := value.Find("title").AsString()
		d.SetCategoryName(id, name, 0)
		if err = d.PullCategory(id); err != nil {
			return
		}
	}

	return
}
