package event

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"wxcloudrun-golang/internal/pkg/model"
)

var cosLink = "cloud://prod-2gicsblt193f5dc8.7072-prod-2gicsblt193f5dc8-1318337180/"

type Service struct {
	VideoDao   *model.Video
	CourtDao   *model.Court
	CollectDao *model.Collect
}

func NewService() *Service {
	return &Service{
		VideoDao:   &model.Video{},
		CourtDao:   &model.Court{},
		CollectDao: &model.Collect{},
	}
}

type Event struct {
	StartTime int32  `json:"start_time"`
	EndTime   int32  `json:"end_time"`
	CourtName string `json:"court_name"`
	Status    int32  `json:"status"`
}

type EventDetail struct {
	VideoSeries []*VideoSeries `json:"video_series"`
}

type VideoSeries struct {
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Status    int32    `json:"status"`
	Videos    []*Video `json:"videos"`
}

type Video struct {
	IsCollected bool   `json:"is_collected"`
	Url         string `json:"url"`
	PicUrl      string `json:"pic_url"`
}

func (s *Service) GetEvents(courtID string, date int32) ([]Event, error) {
	// get today's date like 20210101
	results := make([]Event, 0)
	// get cos links
	hours, err := s.VideoDao.GetDistinctHours(date)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(hours) == 0 {
		return results, nil
	}
	// get events by hours
	for _, hour := range hours {
		results = append(results, Event{StartTime: hour, EndTime: hour + 1, CourtName: courtID, Status: 0})
	}
	if time.Now().Hour() == int(hours[0]) {
		results[0].Status = 1
	} else if time.Now().Hour() == int(hours[0])+1 && time.Now().Minute() < 10 {
		results[0].Status = 1
	}
	return results, nil
}

func (s *Service) GetVideos(date int32, courtID int32, hour int32, openID string) (*EventDetail, error) {
	eventDetail := &EventDetail{VideoSeries: []*VideoSeries{}}
	videos, err := s.VideoDao.GetVideos(date, courtID, hour, 1)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pictures, err := s.VideoDao.GetPictures(date, courtID, hour, 1)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firstHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "00"), EndTime: fmt.Sprintf("%d:%s", hour,
		"15"), Videos: []*Video{}}
	secondHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "15"), EndTime: fmt.Sprintf("%d:%s", hour,
		"30"), Videos: []*Video{}}
	thirdVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "30"), EndTime: fmt.Sprintf("%d:%s", hour,
		"45"), Videos: []*Video{}}
	fourthVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "45"), EndTime: fmt.Sprintf("%d:%s", hour+1,
		"00"), Videos: []*Video{}}
	for index := range videos {
		isCollected := false
		collects, err := s.CollectDao.Gets(&model.Collect{OpenID: openID, Status: 1, FileID: videos[index].FilePath})
		if err != nil {
			return nil, err
		}
		if len(collects) > 0 {
			isCollected = true
		}
		links := strings.Split(videos[index].FileName, "-")
		minuteString := strings.Split(links[1], ".")[0]
		minute, _ := strconv.Atoi(minuteString)
		if minute <= 15 {
			firstHalfVideo.Videos = append(firstHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else if minute > 15 && minute <= 30 {
			secondHalfVideo.Videos = append(secondHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else if minute > 30 && minute <= 45 {
			thirdVideo.Videos = append(thirdVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else {
			fourthVideo.Videos = append(fourthVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		}
	}
	// if date is not today, return
	if time.Now().Format("20060102") != strconv.Itoa(int(date)) {
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo, thirdVideo, secondHalfVideo, firstHalfVideo)
		return eventDetail, nil
	}
	if len(fourthVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) || (time.Now().Hour() == int(hour)+1 && time.Now().Minute() < 10) {
			fourthVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo)
	}
	if len(thirdVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 55 && len(fourthVideo.Videos) == 0 {
			thirdVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, thirdVideo)
	}
	if len(secondHalfVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 40 && len(thirdVideo.Videos) == 0 {
			secondHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, secondHalfVideo)
	}
	if len(firstHalfVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 25 && len(secondHalfVideo.Videos) == 0 {
			firstHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, firstHalfVideo)
	}
	return eventDetail, nil
}

func (s *Service) GetRecord(date int32, courtID int32, hour int32, openID string) (*EventDetail, error) {
	videos, err := s.VideoDao.GetVideos(date, courtID, hour, 2)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pictures, err := s.VideoDao.GetPictures(date, courtID, hour, 2)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	eventDetail := &EventDetail{VideoSeries: []*VideoSeries{}}
	firstHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "00"), EndTime: fmt.Sprintf("%d:%s", hour,
		"15"), Videos: []*Video{}}
	secondHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "15"), EndTime: fmt.Sprintf("%d:%s", hour,
		"30"), Videos: []*Video{}}
	thirdVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "30"), EndTime: fmt.Sprintf("%d:%s", hour,
		"45"), Videos: []*Video{}}
	fourthVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "45"), EndTime: fmt.Sprintf("%d:%s", hour+1,
		"00"), Videos: []*Video{}}
	for index := range videos {
		isCollected := false
		collects, err := s.CollectDao.Gets(&model.Collect{OpenID: openID, Status: 1, FileID: videos[index].FilePath})
		if err != nil {
			return nil, err
		}
		if len(collects) > 0 {
			isCollected = true
		}
		links := strings.Split(videos[index].FileName, "-")
		minuteString := strings.Split(links[1], ".")[0]
		minute, _ := strconv.Atoi(minuteString)
		if minute <= 15 {
			firstHalfVideo.Videos = append(firstHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else if minute > 15 && minute <= 30 {
			secondHalfVideo.Videos = append(secondHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else if minute > 30 && minute <= 45 {
			thirdVideo.Videos = append(thirdVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		} else {
			fourthVideo.Videos = append(fourthVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         videos[index].FilePath,
				PicUrl:      pictures[index].FilePath,
			})
		}
	}
	// if date is not today, return
	if time.Now().Format("20060102") != strconv.Itoa(int(date)) {
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo, thirdVideo, firstHalfVideo, secondHalfVideo)

		return eventDetail, nil
	}
	if len(fourthVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) || (time.Now().Hour() == int(hour)+1 && time.Now().Minute() < 10) {
			fourthVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo)
	}
	if len(thirdVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 55 && len(fourthVideo.Videos) == 0 {
			thirdVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, thirdVideo)
	}
	if len(secondHalfVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 40 && len(thirdVideo.Videos) == 0 {
			secondHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, secondHalfVideo)
	}
	if len(firstHalfVideo.Videos) > 0 {
		if time.Now().Hour() == int(hour) && time.Now().Minute() < 25 && len(secondHalfVideo.Videos) == 0 {
			firstHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, firstHalfVideo)
	}
	return eventDetail, nil
}

func (s *Service) StoreVideo(video *model.Video) (string, error) {
	// get file path
	var typeString string
	if video.Type == 1 {
		typeString = "highlight"
	} else {
		typeString = "record"
	}
	filePath := fmt.Sprintf("%s%s/court%d/%d/%s", cosLink, typeString, video.Court, video.Date, video.FileName)
	record, err := s.VideoDao.Create(&model.Video{
		FilePath:    filePath,
		Date:        video.Date,
		Hour:        video.Hour,
		FileName:    video.FileName,
		Type:        video.Type,
		Court:       video.Court,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	return record.FilePath, nil
}
