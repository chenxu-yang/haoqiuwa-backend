package collect

import (
	"fmt"
	"time"
	"wxcloudrun-golang/internal/pkg/model"
)

type Service struct {
	CollectDao   *model.Collect
	UserEventDao *model.UserEvent
	SurveyDao    *model.Survey
}

func NewService() *Service {
	return &Service{
		CollectDao:   &model.Collect{},
		UserEventDao: &model.UserEvent{},
		SurveyDao:    &model.Survey{},
	}
}

func (s *Service) ToggleCollectVideo(openID string, fileID string, picURL string, videoType int32) (*model.Collect,
	error) {
	// 查询是否已经收藏过
	collects, err := s.CollectDao.Gets(&model.Collect{OpenID: openID, FileID: fileID})
	fmt.Println(collects)
	if err != nil {
		return nil, err
	}
	if len(collects) > 0 {
		collect, err := s.CollectDao.Update(&model.Collect{ID: collects[0].ID, Status: collects[0].Status * (-1)})
		if err != nil {
			return nil, err
		}
		return collect, nil
	}
	// 创建收藏
	if videoType == 0 {
		videoType = 1
	}
	collect, err := s.CollectDao.Create(&model.Collect{
		OpenID:      openID,
		FileID:      fileID,
		PicURL:      picURL,
		VideoType:   videoType,
		Status:      1,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return collect, nil
}

func (s *Service) GetCollectByUser(userOpenID string, videoType int32) ([]model.Collect, error) {
	if videoType == 0 {
		videoType = 1
	}
	collects, err := s.CollectDao.Gets(&model.Collect{OpenID: userOpenID, Status: 1, VideoType: videoType})
	if err != nil {
		return nil, err
	}
	return collects, nil
}

func (s *Service) CollectUserEvent(openID string, fileID string, eventType int32, fromPage string,
	videoType int32) (string,
	error) {
	data, err := s.UserEventDao.Create(&model.UserEvent{
		OpenID:    openID,
		FileID:    fileID,
		EventType: eventType,
		FromPage:  fromPage,
		VideoType: videoType,
	})
	if err != nil {
		return "", err
	}
	return data.FileID, nil
}

func (s *Service) CreateSurvey(openID string, content string) (*model.Survey, error) {
	data, err := s.SurveyDao.Create(&model.Survey{OpenID: openID, Content: content, CreatedTime: time.Now(), UpdatedTime: time.Now()})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetUserDownload(openID string) (int32, error) {
	data, err := s.UserEventDao.Gets(&model.UserEvent{OpenID: openID, EventType: 2})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int32(len(data)), nil
}

func (s *Service) GetUserDownloads(openID string, videoType int32) ([]model.UserEvent, error) {
	data, err := s.UserEventDao.Gets(&model.UserEvent{OpenID: openID, EventType: 2, VideoType: videoType})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}
