package service

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"wxcloudrun-golang/internal/app/collect"
	"wxcloudrun-golang/internal/app/court"
	"wxcloudrun-golang/internal/app/event"
	"wxcloudrun-golang/internal/app/recommend"
	"wxcloudrun-golang/internal/app/user"
	"wxcloudrun-golang/internal/pkg/model"
	"wxcloudrun-golang/internal/pkg/resp"

	"github.com/gin-gonic/gin"
)

type Service struct {
	UserService      *user.Service
	CourtService     *court.Service
	EventService     *event.Service
	CollectService   *collect.Service
	RecommendService *recommend.Service
}

func NewService() *Service {
	return &Service{
		UserService:      user.NewService(),
		CourtService:     court.NewService(),
		EventService:     event.NewService(),
		CollectService:   collect.NewService(),
		RecommendService: recommend.NewService(),
	}
}

// WeChatLogin /wechat/applet_login?code=xxx [get]  路由
// 微信小程序登录
func (s *Service) WeChatLogin(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	var phoneReq PhoneReq
	body, _ := ioutil.ReadAll(c.Request.Body)
	_ = json.Unmarshal(body, &phoneReq)
	// 根据code获取 openID 和 session_key
	wxLoginResp, err := s.UserService.WXLogin(openID, phoneReq.CloudID)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(wxLoginResp, err))
}

type courtReq struct {
	Court int32 `json:"court"`
}

// StoreCourt
func (s *Service) StoreCourt(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	var courtReq courtReq
	_ = json.Unmarshal(body, &courtReq)
	err := s.UserService.StoreCourt(openID, courtReq.Court)
	c.JSON(200, resp.ToStruct(nil, err))
}

// 主页面相关

type PhoneReq struct {
	CloudID string `json:"cloud_id"`
}

// GetUserPhone
func (s *Service) GetUserPhone(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}

}

// ToggleCollectVideo 收藏视频
func (s *Service) ToggleCollectVideo(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	newCollect := &model.Collect{}
	err := json.Unmarshal(body, newCollect)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	collectRecord, err := s.CollectService.ToggleCollectVideo(openID, newCollect.FileID, newCollect.PicURL)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(collectRecord, err))
}

// GetCounts 获取场地
func (s *Service) GetCounts(c *gin.Context) {
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	counts, err := s.CourtService.GetCourts(latitude, longitude)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(counts, err))
}

func (s *Service) GetCountInfo(c *gin.Context) {
	countID := c.Param("id")
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	countIDInt, _ := strconv.Atoi(countID)
	countInfo, err := s.CourtService.GetCountInfo(int32(countIDInt), latitude, longitude)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(countInfo, err))
}

// GetEvents 获取用户所属事件的视频
func (s *Service) GetEvents(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	courtID := c.Query("court")
	results, err := s.EventService.GetEvents(courtID)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(results, err))
}

// GetVideos 获取事件
func (s *Service) GetVideos(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	courtID := c.Query("court")
	hour := c.Query("hour")
	hourInt, _ := strconv.Atoi(hour)
	event, err := s.EventService.GetVideos(courtID, hourInt, openID)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(event, err))
}

// GetRecords 获取录像
func (s *Service) GetRecords(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	courtID := c.Query("court")
	hour := c.Query("hour")
	hourInt, _ := strconv.Atoi(hour)
	data, err := s.EventService.GetRecord(courtID, hourInt, openID)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(data, err))
}

// GetCollectVideos 获取用户收藏的视频
func (s *Service) GetCollectVideos(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	collects, err := s.CollectService.GetCollectByUser(openID)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(collects, err))
}

// GetRecommendVideos 获取推荐视频
func (s *Service) GetRecommendVideos(c *gin.Context) {
	videos, err := s.RecommendService.GetRecommend()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(videos, err))
}

// JudgeLocation 判断用户是否在场地内
func (s *Service) JudgeLocation(c *gin.Context) {
	countID := c.Param("id")
	countIDInt, _ := strconv.Atoi(countID)
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	result, err := s.CourtService.JudgeLocation(int32(countIDInt), latitude, longitude)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(result, err))
}

// CollectUserEvent 下载视频记录
func (s *Service) CollectUserEvent(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	userEvent := &model.UserEvent{}
	err := json.Unmarshal(body, userEvent)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	data, err := s.CollectService.CollectUserEvent(openID, userEvent.FileID, userEvent.EventType, userEvent.FromPage,
		userEvent.VideoType)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(data, err))
}

// CollectSurvey 下载问卷记录
func (s *Service) CollectSurvey(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	data, err := s.CollectService.CreateSurvey(openID, string(body))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(data, err))
}

// GetUserDownload
func (s *Service) GetUserDownload(c *gin.Context) {
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		c.JSON(400, "请先登录")
		return
	}
	data, err := s.CollectService.GetUserDownload(openID)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp.ToStruct(data, err))
}
