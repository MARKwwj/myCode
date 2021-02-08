package webmanhuawu

import (
	"TianlangCapturer/src/api"
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	// 导入jpeg解析库
	_ "image/jpeg"
	_ "image/png"
)

// BookInfo BookInfo
type BookInfo struct {
	ID        int            `json:"id"`
	URL       string         `json:"url"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	MD5Name   string         `json:"md5Name"`
	CoverURL  string         `json:"coverURL"`
	Describe  string         `json:"introduce"`
	TotalNum  int            `json:"totalNum"`
	Sections  []*BookSection `json:"chapters"`
	PicSuffix string         `json:"picSuffix"`
}

// GetMD5 GetMD5
func (p *BookInfo) GetMD5() string {
	if p.MD5Name == "" {
		p.MD5Name = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("mhqwe.xyz|%d", p.ID))))
	}
	return p.MD5Name
}

// BookSection BookSection
type BookSection struct {
	URL           string   `json:"url"`
	Name          string   `json:"chapter"`
	Title         string   `json:"subTitle"`
	SortID        int      `json:"sort"`
	MinPage       int      `json:"minPage"`
	MaxPage       int      `json:"maxPage"`
	ImagesURL     []string `json:"imagesURL"`
	TotalFileSize int64    `json:"totalFileSize"`
	mBookInfo     *BookInfo
}

// GraphicNovel GraphicNovel
type GraphicNovel struct {
	mRootURL      string
	mBookListSync sync.WaitGroup
}

// GetBookList GetBookList
func (p *GraphicNovel) GetBookList(pageID int) (result []*BookInfo, err error) {
	api.Info("[GetBookList]PageID:%d", pageID)
	var (
		doc  *goquery.Document
		node *goquery.Selection
	)
	if doc, err = p.doGetDocument(fmt.Sprintf("http://www.mhqwe.xyz/?page.currentPage=%d&queryFilm.orderType=0&filmName=", pageID)); err != nil {
		return
	}

	if node = doc.Find(".book-list"); node == nil {
		return nil, fmt.Errorf("not found '.book-list'")
	}

	if node = node.Find(".book-list-item"); node == nil {
		return nil, fmt.Errorf("not found '.book-list-item'")
	}

	node.Each(func(i int, s *goquery.Selection) {
		if value, ok := s.Attr("onclick"); ok {
			LIndex := strings.Index(value, "'")
			RIndex := strings.LastIndex(value, "'")
			URL := value[LIndex+1 : RIndex]

			LIndex = strings.LastIndex(URL, "/")
			RIndex = strings.LastIndex(URL, ".")
			nID, _ := strconv.ParseInt(URL[LIndex+1:RIndex], 10, 32)

			info := s.Find(".book-list-info")
			if info == nil {
				panic(fmt.Errorf("not found '.book-list-info'"))
			}

			title := info.Find(".book-list-info-title")
			if title == nil {
				panic(fmt.Errorf("not found '.book-list-info-title'"))
			}

			describe := info.Find(".book-list-info-desc")
			if describe == nil {
				panic(fmt.Errorf("not found '.book-list-info-desc'"))
			}

			bookInfo := &BookInfo{
				ID:       int(nID),
				URL:      URL,
				Title:    strings.TrimSpace(title.Text()),
				Describe: strings.TrimSpace(describe.Text()),
			}
			bookInfo.GetMD5()
			result = append(result, bookInfo)

			if _, ok := mHistory.FindRecord(bookInfo.GetMD5()); ok {
				return
			}

			p.mBookListSync.Add(1)
			go func() {
				defer p.mBookListSync.Done()
				if err = p.GetSections(bookInfo); err != nil {
					panic(err)
				}
			}()

		}
	})

	p.mBookListSync.Wait()
	return
}

// GetSections GetSections
func (p *GraphicNovel) GetSections(info *BookInfo) (err error) {
	api.Info("[GetSections]BookID:%d BookMD5:%s URL:%s", info.ID, info.MD5Name, info.URL)
	var (
		doc  *goquery.Document
		node *goquery.Selection
	)
	if doc, err = p.doGetDocument(info.URL); err != nil {
		return
	}

	if node = doc.Find(".book-header").Find(".writer").Find("img"); node == nil {
		return fmt.Errorf("not found '.book-header'")
	}
	info.CoverURL = node.AttrOr("src", "")

	if node = doc.Find(".book-body"); node == nil {
		return fmt.Errorf("not found '.book-body'")
	}

	if node = doc.Find(".list-item"); node == nil {
		return fmt.Errorf("not found '.list-item'")
	}
	node.Each(func(i int, s *goquery.Selection) {
		value, _ := s.Find("a").Attr("onclick")
		valueArr := strings.Split(value, "'")
		if len(valueArr) != 5 {
			panic(fmt.Errorf("parse onclick error: %s", value))
		}
		bookID, linkID := valueArr[1], valueArr[3]

		section := &BookSection{
			mBookInfo: info,
			SortID:    i,
			Name:      fmt.Sprintf("第%d话", i),
			URL:       p.parseURL(bookID, linkID),
		}

		if node := s.Find(".cell-title"); node != nil {
			section.Title = node.Text()
		}

		if strings.TrimSpace(section.Title) == "开始阅读" {
			section.Title = "第1话"
		}

		if i > 0 {
			section.MinPage = info.Sections[len(info.Sections)-1].MaxPage
		}

		if err := p.GetCartoons(section); err != nil {
			panic(err)
		}

		info.Sections = append(info.Sections, section)
	})

	return
}

// GetCartoons GetCartoons
func (p *GraphicNovel) GetCartoons(section *BookSection) (err error) {
	api.Info("[GetCartoons]BookID:%d SectionSortID:%d MinPage:%d URL:%s", section.mBookInfo.ID, section.SortID, section.MinPage, section.URL)
	var (
		doc  *goquery.Document
		node *goquery.Selection
	)
	if doc, err = p.doGetDocument(section.URL); err != nil {
		return
	}

	if node = doc.Find("#imgList"); node == nil {
		return fmt.Errorf("not found '#imgList'")
	}

	if node = node.Find("img"); node == nil {
		return fmt.Errorf("not found 'img'")
	}

	section.ImagesURL = make([]string, len(node.Nodes))
	node.Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("src"); ok {
			section.ImagesURL[i] = url
		}
	})
	section.MaxPage = section.MinPage + len(section.ImagesURL)
	return nil
}

// DownloadBookAllData DownloadBookAllData
func (p *GraphicNovel) DownloadBookAllData(info *BookInfo) bool {
	api.Info("ID:%d Title:%s Sections:%d URL:%s", info.ID, info.Title, len(info.Sections), info.URL)

	//创建主目录
	bookMD5 := info.GetMD5()
	mRootDir := fmt.Sprintf("Data/Manhua/%s", bookMD5)
	os.MkdirAll(mRootDir, os.ModePerm)

	// ChapterInfo ChapterInfo
	type ChapterInfo struct {
		SubTitle    string   `json:"subTitle"`
		Sort        int      `json:"sort"`
		ChapterURLs []string `json:"chapterUrl"`
		ChapterSize int      `json:"chapterSize"`
	}

	// BookJSON BookJSON
	type BookJSON struct {
		Title     string        `json:"title"`
		Introduce string        `json:"introduce"`
		PicURL    string        `json:"picUrl"`
		Type      string        `json:"type"`
		TotalNum  int           `json:"totalNum"`
		MachineID string        `json:"machineId"`
		Chapters  []ChapterInfo `json:"chapters"`
	}

	//API请求JSON数据
	var bookJSON BookJSON
	bookJSON.Type = info.Type
	bookJSON.Title = info.Title
	bookJSON.Introduce = info.Describe
	bookJSON.MachineID = mConfigInfo.MachineID
	bookJSON.PicURL = fmt.Sprintf("cartoonInfo/%s/cover.jpg", bookMD5)
	bookJSON.TotalNum = len(info.Sections)

	//下载封面图
	if err := NewDownloadContext(p.mRootURL+info.CoverURL, filepath.Join(mRootDir, "cover.jpg")).WaitDownload(); err != nil {
		api.Error("Error downloading cover drawing：%s", err.Error())
		return false
	}

	//创建话目录
	info.TotalNum = len(info.Sections)
	for i := 1; i <= info.TotalNum; i++ {
		os.MkdirAll(filepath.Join(mRootDir, fmt.Sprintf("%d", i)), os.ModePerm)
	}

	//下载每话数据
	for i, sections := range info.Sections {
		var chapterInfo ChapterInfo
		chapterInfo.Sort = i + 1
		chapterInfo.SubTitle = sections.Title

		mSaveDir := fmt.Sprintf("%s/%d/", mRootDir, chapterInfo.Sort)
		for index, url := range sections.ImagesURL {
			imgPath := filepath.Join(mSaveDir, fmt.Sprintf("%d.jpg", index))
			context := NewDownloadContext(p.mRootURL+url, imgPath)
			if err := context.WaitDownload(); err != nil {
				api.Error("Error downloading %s", err.Error())
				return false
			}
			sections.TotalFileSize += context.mTotalFileSize

			imgFile, err := os.Open(imgPath)
			if err != nil {
				api.Error("Error open image %s ImgPath:%s", err.Error(), imgPath)
				return false
			}
			config, _, err := image.DecodeConfig(imgFile)
			imgFile.Close()
			if err != nil {
				api.Error("Error decode image %s ImgPath:%s", err.Error(), imgPath)
				return false
			}
			ratio := float32(config.Width) / float32(config.Height)
			chapterURL := fmt.Sprintf("cartoonInfo/%s/%d/%d.jpg_%f", bookMD5, chapterInfo.Sort, index, ratio)
			api.Debug("chapterURL:%s", chapterURL)
			chapterInfo.ChapterURLs = append(chapterInfo.ChapterURLs, chapterURL)
		}

		chapterInfo.ChapterSize = int(sections.TotalFileSize)
		bookJSON.Chapters = append(bookJSON.Chapters, chapterInfo)
	}

	//测试模式下拦截
	if api.IsDebug() {
		return true
	}

	//加密目录下所有文件
	api.XORAllFileData(mRootDir, []byte{
		0x59, 0x3d, 0x94, 0x69, 0x28, 0xed, 0xff, 0xe3, 0xc3, 0x6f, 0xac, 0xc5, 0xec, 0x52, 0x2b, 0xa3,
		0x79, 0x4e, 0xf5, 0x32, 0x43, 0x75, 0x88, 0x03, 0x1a, 0x84, 0x34, 0xcc, 0xb6, 0x53, 0x0d, 0x92,
		0x15, 0xd7, 0x2f, 0x30, 0xbd, 0x60, 0xb5, 0x17, 0x01, 0x9e, 0xdb, 0xb8, 0x56, 0x70, 0x33, 0x54,
		0x2e, 0x3b, 0xbf, 0x6a, 0x8c, 0x04, 0x41, 0xad, 0x6d, 0xf8, 0x58, 0x35, 0x98, 0x99, 0x24, 0x73,
		0x25, 0xf2, 0xb1, 0x5a, 0xb3, 0xc4, 0x8f, 0xd9, 0xef, 0xfb, 0x45, 0xc1, 0x37, 0x2a, 0x93, 0x4c,
		0x86, 0xda, 0x09, 0xae, 0x8d, 0xd8, 0x8a, 0x81, 0x7c, 0x44, 0x6b, 0xea, 0xf9, 0x66, 0x40, 0xa9,
		0xb9, 0x7a, 0x38, 0xe5, 0x29, 0xfc, 0x7f, 0x12, 0x4f, 0xcd, 0xba, 0xc7, 0x6c, 0xd6, 0xa0, 0x10,
		0xe2, 0x5e, 0xf0, 0xd1, 0xfd, 0xbb, 0xa6, 0x63, 0x05, 0xd5, 0x22, 0x9b, 0x9a, 0x00, 0xd3, 0x61,
		0x48, 0xaf, 0xee, 0xa7, 0x46, 0x77, 0x1f, 0x71, 0xde, 0x02, 0x42, 0x9c, 0xa2, 0xc2, 0xcf, 0xce,
		0x9d, 0x64, 0xf6, 0xca, 0xab, 0x14, 0x36, 0x0b, 0xd0, 0xdc, 0x3e, 0x7e, 0x0e, 0x72, 0x5f, 0x20,
		0x49, 0xa5, 0xfe, 0x23, 0xf4, 0x51, 0x95, 0x89, 0x87, 0xb0, 0x5d, 0x0f, 0x2c, 0x39, 0xa4, 0x5c,
		0xf1, 0x13, 0xcb, 0x57, 0x06, 0xf3, 0xeb, 0x97, 0x0c, 0x18, 0xb7, 0x21, 0xbc, 0x90, 0xe0, 0xb2,
		0x96, 0x1b, 0x27, 0xe9, 0x74, 0x19, 0x67, 0x6e, 0xc9, 0x55, 0xc0, 0x9f, 0xe8, 0xbe, 0xd4, 0x7b,
		0x83, 0x16, 0xd2, 0x2d, 0x4a, 0x0a, 0x7d, 0xb4, 0x82, 0x4d, 0xdd, 0x85, 0x3c, 0xe6, 0x50, 0x4b,
		0xc6, 0x80, 0xf7, 0x1d, 0xe7, 0x76, 0x3f, 0xfa, 0xe1, 0x78, 0xa8, 0x68, 0xa1, 0x1c, 0x91, 0xdf,
		0xaa, 0x3a, 0x26, 0x08, 0x8e, 0x62, 0x47, 0x5b, 0x1e, 0x65, 0xc8, 0x07, 0x11, 0x8b, 0x31, 0xe4,
	})

	//SCP上传到资源服务器
	if msg, err := api.SCP(mConfigInfo.SCPAddr, mConfigInfo.SCPPort, mConfigInfo.SCPUser, mConfigInfo.SCPPass, mRootDir, mConfigInfo.SavePath); err == nil {
		api.Info("[SCP]%s Done", bookMD5)
	} else {
		api.Error("[SCP]%s Error:%s\ndetails:\n%s", bookMD5, err.Error(), msg)
		return false
	}

	//HTTP通知资源接口
	jsonData, err := json.Marshal(&bookJSON)
	if err != nil {
		api.Error(err.Error())
		return false
	}
	if api.IsDebug() {
		api.Debug("JsonData:%s", jsonData)
	}
	if resp, err := http.Post(fmt.Sprintf(mConfigInfo.WebAPI+"cartoon/importData"), "application/json", bytes.NewReader(jsonData)); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		api.Info("[WebAPI]%s StateCode:%d Details:\n%s", bookMD5, resp.StatusCode, string(data))
	} else {
		api.Error("[WebAPI]%s Error:%s", bookMD5, err.Error())
		return false
	}

	//清空目录
	api.Info("[Clear]%s", mRootDir)
	os.RemoveAll(mRootDir)
	return true
}

func (p *GraphicNovel) doGetDocument(url string) (*goquery.Document, error) {
	if url[0] == '/' {
		url = p.mRootURL + url
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("res.StatusCode:%d", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (p *GraphicNovel) parseURL(bookID, linkID string) string {
	// aesKey
	block, err := aes.NewCipher([]byte("12cdefgabcdefg12"))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// padding
	srcData := []byte(linkID)
	nBlockSize := block.BlockSize()
	paddingNum := nBlockSize - (len(srcData) % nBlockSize)
	for i := 0; i < paddingNum; i++ {
		srcData = append(srcData, byte(paddingNum))
	}
	// encrypt
	encryptData := make([]byte, len(srcData))
	tmpData := make([]byte, nBlockSize)
	for index := 0; index < len(srcData); index += nBlockSize {
		block.Encrypt(tmpData, srcData[index:index+nBlockSize])
		copy(encryptData, tmpData)
	}
	return fmt.Sprintf("/play?linkId=%s&bookId=%s&key=%s", linkID, bookID, base64.StdEncoding.EncodeToString(encryptData))
}

// DownloadContext DownloadContext
type DownloadContext struct {
	mURL             string
	mSavePath        string
	mMaxRetryCount   int
	mCurreFileSize   int64
	mTotalFileSize   int64
	mDownloadingByte int64
}

// NewDownloadContext NewDownloadContext
func NewDownloadContext(url, filePath string) *DownloadContext {
	context := &DownloadContext{mURL: url, mSavePath: filePath, mMaxRetryCount: 10}
	return context
}

// WaitDownload WaitDownload
func (context *DownloadContext) WaitDownload() (err error) {
	for i := 0; i < context.mMaxRetryCount; i++ {
		api.Info("[Downloader]Try:%d %s << %s", i, context.mSavePath, context.mURL)
		if err = context.tryDownload(); err == nil {
			context.mCurreFileSize = 0
			return nil
		}
		time.Sleep(time.Second)
	}
	return
}

// tryDownload
func (context *DownloadContext) tryDownload() error {
	reqest, err := http.NewRequest("GET", context.mURL, nil)
	if context.mCurreFileSize > 0 {
		reqest.Header.Set("Range", fmt.Sprintf("bytes=%d-", context.mCurreFileSize))
	}

	mClient := &http.Client{Timeout: time.Second * 10}
	resp, err := mClient.Do(reqest)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("StatusCode:%d", resp.StatusCode)
	}

	if context.mTotalFileSize <= 0 {
		// context.mCurreFileSize = 0
		context.mTotalFileSize = int64(resp.ContentLength)
	}

	mSaveFile, err := os.OpenFile(context.mSavePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	mSaveFile.Seek(context.mCurreFileSize, os.SEEK_SET)

	nCount := int64(0)
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			mSaveFile.Write(buffer[:n])
			nCount += int64(n)
			context.mCurreFileSize += int64(n)
			context.mDownloadingByte += int64(n)
		}
		if err != nil {
			if err != io.EOF && api.IsDebug() {
				api.Warn("[Download]progress(%d/%d) fail \n%s", nCount, resp.ContentLength, err.Error())
			}
			break
		}
	}

	mSaveFile.Close()
	resp.Body.Close()
	if context.mCurreFileSize < context.mTotalFileSize {
		return fmt.Errorf("nCurreFileSize(%d) != nTotalFileSize(%d)", context.mCurreFileSize, context.mTotalFileSize)
	}
	return nil
}
