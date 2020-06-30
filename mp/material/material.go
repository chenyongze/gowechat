package material

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/chenyongze/gowechat/mp/base"
	"github.com/chenyongze/gowechat/util"
	"github.com/chenyongze/gowechat/wxcontext"
)

const (
	addNewsURL          = "https://api.weixin.qq.com/cgi-bin/material/add_news"
	addMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/add_material"
	delMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/del_material"
	getMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/get_material"
	batchgetMaterialURL = "https://api.weixin.qq.com/cgi-bin/material/batchget_material"
)

//Material 素材管理
type Material struct {
	base.MpBase
}

//NewMaterial init
func NewMaterial(context *wxcontext.Context) *Material {
	material := new(Material)
	material.Context = context
	return material
}

//Article 永久图文素材
type Article struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
}

//reqArticles 永久性图文素材请求信息
type reqArticles struct {
	Articles []*Article `json:"articles"`
}

//resArticles 永久性图文素材返回结果
type resArticles struct {
	util.CommonError

	MediaID string `json:"media_id"`
}

//AddNews 新增永久图文素材
func (material *Material) AddNews(articles []*Article) (mediaID string, err error) {
	req := &reqArticles{articles}

	responseBytes, err := material.HTTPPostJSONWithAccessToken(addNewsURL, req)
	if err != nil {
		return
	}

	var res resArticles
	err = json.Unmarshal(responseBytes, res)
	if err != nil {
		return
	}
	mediaID = res.MediaID
	return
}

//PostMediaType PostMediaType
type resMedia struct {
	MediaID string `json:"media_id"`
}

//ArticleDetail 永久图文素材
type ArticleDetail struct {
	Title              string `json:"title"`
	ThumbMediaID       string `json:"thumb_media_id"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	ShowCoverPic       int    `json:"show_cover_pic"`
	Content            string `json:"content"`
	ContentSourceURL   string `json:"content_source_url"`
	URL                string `json:"url"`
	ThumbURL           string `json:"thumb_url"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}

//ReqNews ReqNews
type ReqNews struct {
	CreateTime int         `json:"create_time"`
	UpdateTime int         `json:"update_time"`
	NewsItem   []*ArticleDetail `json:"news_item"`
}

// GetNews 获取素材
func (material *Material) GetNews(mediaID string) (res *ReqNews, err error) {
	req := &resMedia{
		MediaID: mediaID,
	}
	responseBytes, err := material.HTTPPostJSONWithAccessToken(getMaterialURL, req)
	fmt.Println("responseBytes:", responseBytes)
	fmt.Println("err1:", err)
	if err != nil {
		return
	}

	res = new(ReqNews)
	err = json.Unmarshal(responseBytes, res)
	fmt.Println("res:", res)
	fmt.Println("err2:", err)
	if err != nil {
		return
	}
	// mediaID = res.MediaID
	return
}

//ReqNewsList ReqNewsList
type ReqNewsList struct {
	Item       []*ReqNewsListItem `json:"item"`
	TotalCount int                `json:"total_count"`
	ItemCount  int                `json:"item_count"`
}

// ReqNewsListItem ReqNewsListItem
type ReqNewsListItem struct {
	MediaID    string   `json:"media_id"`
	UpdateTime int      `json:"update_time"`
	Content    *ReqNews `json:"content"`
}

//resNewsList resNewsList
// type=news offset=0 count=1
type resNewsList struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}

// GetNewsList GetNewsList
func (material *Material) GetNewsList(type_cn string, offset, count int) (res *ReqNewsList, err error) {
	req := &resNewsList{
		Type:   type_cn,
		Offset: offset,
		Count:  count,
	}
	responseBytes, err := material.HTTPPostJSONWithAccessToken(batchgetMaterialURL, req)
	fmt.Println("responseBytes:", responseBytes)
	fmt.Println("err1:", err)
	if err != nil {
		return
	}

	res = new(ReqNewsList)
	err = json.Unmarshal(responseBytes, res)
	fmt.Println("res:", res)
	fmt.Println("err2:", err)
	if err != nil {
		return
	}
	// mediaID = res.MediaID
	return
}

//resAddMaterial 永久性素材上传返回的结果
type resAddMaterial struct {
	util.CommonError

	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}

//AddMaterial 上传永久性素材（处理视频需要单独上传）
func (material *Material) AddMaterial(mediaType MediaType, filename string) (mediaID string, url string, err error) {
	if mediaType == MediaTypeVideo {
		err = errors.New("永久视频素材上传使用 AddVideo 方法")
	}
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&type=%s", addMaterialURL, accessToken, mediaType)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	var resMaterial resAddMaterial
	err = json.Unmarshal(response, &resMaterial)
	if err != nil {
		return
	}
	if resMaterial.ErrCode != 0 {
		err = fmt.Errorf("AddMaterial error : errcode=%v , errmsg=%v", resMaterial.ErrCode, resMaterial.ErrMsg)
		return
	}
	mediaID = resMaterial.MediaID
	url = resMaterial.URL
	return
}

type reqVideo struct {
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

//AddVideo 永久视频素材文件上传
func (material *Material) AddVideo(filename, title, introduction string) (mediaID string, url string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&type=video", addMaterialURL, accessToken)

	videoDesc := &reqVideo{
		Title:        title,
		Introduction: introduction,
	}
	var fieldValue []byte
	fieldValue, err = json.Marshal(videoDesc)
	if err != nil {
		return
	}

	fields := []util.MultipartFormField{
		{
			IsFile:    true,
			Fieldname: "video",
			Filename:  filename,
		},
		{
			IsFile:    true,
			Fieldname: "description",
			Value:     fieldValue,
		},
	}

	var response []byte
	response, err = util.PostMultipartForm(fields, uri)
	if err != nil {
		return
	}

	var resMaterial resAddMaterial
	err = json.Unmarshal(response, &resMaterial)
	if err != nil {
		return
	}
	if resMaterial.ErrCode != 0 {
		err = fmt.Errorf("AddMaterial error : errcode=%v , errmsg=%v", resMaterial.ErrCode, resMaterial.ErrMsg)
		return
	}
	mediaID = resMaterial.MediaID
	url = resMaterial.URL
	return
}

type reqDeleteMaterial struct {
	MediaID string `json:"media_id"`
}

//DeleteMaterial 删除永久素材
func (material *Material) DeleteMaterial(mediaID string) error {
	_, err := material.HTTPPostJSONWithAccessToken(delMaterialURL, reqDeleteMaterial{mediaID})
	if err != nil {
		return err
	}
	return nil
}
