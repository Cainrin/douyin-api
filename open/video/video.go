package video

import (
	"encoding/json"
	"fmt"

	"github.com/amorist/douyin/open/context"
	"github.com/amorist/douyin/util"
)

const (
	// 上传视频
	videoUploadURL string = "https://open.douyin.com/video/upload?access_token=%s&open_id=%s"
	// 初始化分片上传
	videoPartInitURL string = "https://open.douyin.com/video/part/init?access_token=%s&open_id=%s"
	// 分片上传
	videoPartUploadURL string = "https://open.douyin.com/video/part/upload?access_token=%s&open_id=%s"
	// 分片完成上传
	videoPartCompleteURL string = "https://open.douyin.com/video/part/complete?access_token=%s&open_id=%s&upload_id=%s"
	// 创建视频
	videoCreateURL string = "https://open.douyin.com/video/create?access_token=%s&open_id=%s"
	// 删除视频
	videoDeleteURL string = "https://open.douyin.com/video/delete?access_token=%s&open_id=%s"
)

// Video 视频
type Video struct {
	*context.Context
}

// Info 视频信息.
type Info struct {
	util.CommonError

	Video struct {
		VideoID string `json:"video_id"`
		Height  int64  `json:"height"`
		Width   int64  `json:"width"`
	} `json:"video"`
}

type uploadVideoRes struct {
	Message string `json:"message"`
	Data    Info   `json:"data"`
}

// Upload 视频上传.
// refer: https://open.douyin.com/platform/doc/6848798087398295555
func (video *Video) Upload(openid, filename string) (videoInfo Info, err error) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoUploadURL, accessToken, openid)
	field := util.MultipartFormField{
		Fieldname: "video",
		Filename:  filename,
	}

	var response []byte
	response, err = util.PostMultipartForm(field, uri)
	if err != nil {
		return
	}

	var result uploadVideoRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("Upload error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	videoInfo = result.Data
	return
}

// PartInfo .
type PartInfo struct {
	util.CommonError

	UploadID string `json:"upload_id"`
}

type partInfoRes struct {
	Message string   `json:"message"`
	Data    PartInfo `json:"data"`
}

// PartInit 初始化分片上传.
// refer: https://open.douyin.com/platform/doc/6848798087398393859
func (video *Video) PartInit(openid string) (partInfo PartInfo, err error) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoPartInitURL, accessToken, openid)
	var response []byte
	response, err = util.HTTPPost(uri, "")
	if err != nil {
		return
	}

	var result partInfoRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("PartInit error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	partInfo = result.Data
	return
}

type partVideoRes struct {
	Message string           `json:"message"`
	Data    util.CommonError `json:"data"`
}

// PartUpload 视频分片上传.
// refer: https://open.douyin.com/platform/doc/6848798087226460172
func (video *Video) PartUpload(openid, uploadid string, partNumber int64, filename string) (err error) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoPartUploadURL, accessToken, openid)
	field := util.MultipartFormField{
		Fieldname: "video",
		Filename:  filename,
	}

	var response []byte
	response, err = util.PostMultipartForm(field, uri)
	if err != nil {
		return
	}

	var result partVideoRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("Upload error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	return
}

// PartComplete 视频分片完成上传.
// refer: https://open.douyin.com/platform/doc/6848798087398361091
func (video *Video) PartComplete(openid, uploadid string) (videoInfo Info, err error) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoPartCompleteURL, accessToken, openid, uploadid)
	var response []byte
	response, err = util.HTTPPost(uri, "")
	if err != nil {
		return
	}

	var result uploadVideoRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("PartComplete error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	videoInfo = result.Data
	return
}

// CreateVideoReq .
type CreateVideoReq struct {
	VideoID           string   `json:"video_id"`
	CoverTsp          float64  `json:"cover_tsp"`
	GameID            string   `json:"string"`
	PoiID             string   `json:"poi_id"`
	Text              string   `json:"text"`
	MicroAppURL       string   `json:"micro_app_url"`
	MicroAppID        string   `json:"micro_app_id"`
	MicroAppTitle     string   `json:"micro_app_title"`
	AtUsers           []string `json:"at_users"`
	GameContent       string   `json:"game_content"`
	TimelinessKeyword string   `json:"timeliness_keyword"`
	TimelinessLabel   string   `json:"timeliness_label"`
	ArticleID         string   `json:"article_id"`
	ArticleTitle      string   `json:"article_title"`
}

// CreateInfo .
type CreateInfo struct {
	util.CommonError

	ItemID string `json:"item_id"`
}

type createRes struct {
	Message string     `json:"message"`
	Data    CreateInfo `json:"data"`
}

// Create 视频创建.
// refer: https://open.douyin.com/platform/doc/6848798087398328323
func (video *Video) Create(openid string, videoInfo *CreateVideoReq) (info CreateInfo, err error) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoCreateURL, accessToken, openid)
	var response []byte
	response, err = util.PostJSON(uri, videoInfo)
	if err != nil {
		return
	}

	var result createRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("Create error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	info = result.Data
	return
}

type deleteVideoReq struct {
	ItemID string `json:"item_id"`
}

type deleteVideoRes struct {
	Message string           `json:"message"`
	Data    util.CommonError `json:"data"`
}

// Delete 视频删除
// refer: https://open.douyin.com/platform/doc/6848806536383383560#url
func (video *Video) Delete(openid, itemid string) {
	accessToken, err := video.GetAccessToken(openid)
	if err != nil {
		return
	}

	uri := fmt.Sprintf(videoCreateURL, accessToken, openid)

	rep := &deleteVideoReq{
		ItemID: itemid,
	}

	var response []byte
	response, err = util.PostJSON(uri, rep)
	if err != nil {
		return
	}

	var result deleteVideoRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.Data.ErrCode != 0 {
		err = fmt.Errorf("Delete error : errcode=%v , errmsg=%v", result.Data.ErrCode, result.Data.ErrMsg)
		return
	}
	return
}