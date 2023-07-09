package event

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"wxcloudrun-golang/internal/pkg/model"
	"wxcloudrun-golang/internal/pkg/tcos"
)

type Service struct {
	EventDao   *model.Event
	VideoDao   *model.Video
	CourtDao   *model.Court
	CollectDao *model.Collect
}

func NewService() *Service {
	return &Service{
		EventDao: &model.Event{},
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

func (s *Service) GetEvents(courtID string) ([]Event, error) {
	// get today's date like 20210101
	today := time.Now().Format("20060102")
	results := make([]Event, 0)
	// get cos links
	allLinks, err := tcos.GetCosFileList(fmt.Sprintf("highlight/court%s/%s/v", courtID, today))
	if err != nil {
		log.Printf("get cos file list error: %v", err)
		return nil, err
	}
	// get hours by links, links are like 4042-prod/highlight/court1/20210101/10-32.mp4, 10-32 means hour and minute
	distinctHours := make(map[int]int)
	for _, link := range allLinks {
		links := strings.Split(link, "/")
		hour := strings.Split(links[len(links)-1], "-")[0]
		hourInt, _ := strconv.Atoi(hour[1:])
		distinctHours[hourInt] += 1
	}
	// get hour by order
	hours := make([]int, 0)
	for hour := range distinctHours {
		hours = append(hours, hour)
	}
	if len(hours) == 0 {
		return results, nil
	}
	// sort hours
	sort.Slice(hours, func(i, j int) bool { return hours[i] > hours[j] })
	// get events by hours
	for _, hour := range hours {
		results = append(results, Event{StartTime: int32(hour), EndTime: int32(hour + 1), CourtName: courtID, Status: 0})
	}
	print("time now:", time.Now().Hour())
	if time.Now().Hour() == hours[0] {
		results[0].Status = 1
	} else if time.Now().Hour() == hours[0]+1 && time.Now().Minute() < 10 {
		results[0].Status = 1
	}
	return results, nil
}

func (s *Service) GetVideos(courtID string, hour int, openID string) (*EventDetail, error) {
	today := time.Now().Format("20060102")
	allLinks, err := tcos.GetCosFileList(fmt.Sprintf("highlight/court%s/%s/v%d", courtID, today, hour))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	picLinks, err := tcos.GetCosFileList(fmt.Sprintf("highlight/court%s/%s/p%d", courtID, today, hour))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// order by minute
	sort.Slice(allLinks, func(i, j int) bool {
		ssi := strings.Split(allLinks[i], "/")
		ssj := strings.Split(allLinks[j], "/")
		return strings.Compare(ssi[len(ssi)-1], ssj[len(ssj)-1]) < 0
	})
	sort.Slice(picLinks, func(i, j int) bool {
		ssi := strings.Split(picLinks[i], "/")
		ssj := strings.Split(picLinks[j], "/")
		return strings.Compare(ssi[len(ssi)-1], ssj[len(ssj)-1]) < 0
	})
	eventDetail := &EventDetail{VideoSeries: []*VideoSeries{}}
	firstHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "00"), EndTime: fmt.Sprintf("%d:%s", hour,
		"15")}
	secondHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "15"), EndTime: fmt.Sprintf("%d:%s", hour,
		"30")}
	thirdVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "30"), EndTime: fmt.Sprintf("%d:%s", hour,
		"45")}
	fourthVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "45"), EndTime: fmt.Sprintf("%d:%s", hour+1,
		"00")}
	for index := range allLinks {
		isCollected := false
		collects, err := s.CollectDao.Gets(&model.Collect{OpenID: openID, Status: 1, FileID: allLinks[index]})
		if err != nil {
			return nil, err
		}
		if len(collects) > 0 {
			isCollected = true
		}
		links := strings.Split(allLinks[index], "/")
		minuteString := strings.Split(strings.Split(links[len(links)-1], "-")[1], ".")[0]
		minute, _ := strconv.Atoi(minuteString)
		if minute <= 15 {
			firstHalfVideo.Videos = append(firstHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else if minute > 15 && minute <= 30 {
			secondHalfVideo.Videos = append(secondHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else if minute > 30 && minute <= 45 {
			thirdVideo.Videos = append(thirdVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else {
			fourthVideo.Videos = append(fourthVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		}
	}
	if len(firstHalfVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 25 && len(secondHalfVideo.Videos) == 0 {
			firstHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, firstHalfVideo)
	}
	if len(secondHalfVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 40 && len(thirdVideo.Videos) == 0 {
			secondHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, secondHalfVideo)
	}
	if len(thirdVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 55 && len(fourthVideo.Videos) == 0 {
			thirdVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, thirdVideo)
	}
	if len(fourthVideo.Videos) > 0 {
		if time.Now().Hour() == hour || (time.Now().Hour() == hour+1 && time.Now().Minute() < 10) {
			fourthVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo)
	}
	return eventDetail, nil
}

func (s *Service) GetRecord(courtID string, hour int, openID string) (*EventDetail, error) {
	today := time.Now().Format("20060102")
	allLinks, err := tcos.GetCosFileList(fmt.Sprintf("record/court%s/%s/v%d", courtID, today, hour))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	picLinks, err := tcos.GetCosFileList(fmt.Sprintf("record/court%s/%s/p%d", courtID, today, hour))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// order by minute
	sort.Slice(allLinks, func(i, j int) bool {
		ssi := strings.Split(allLinks[i], "/")
		ssj := strings.Split(allLinks[j], "/")
		return strings.Compare(ssi[len(ssi)-1], ssj[len(ssj)-1]) < 0
	})
	sort.Slice(picLinks, func(i, j int) bool {
		ssi := strings.Split(picLinks[i], "/")
		ssj := strings.Split(picLinks[j], "/")
		return strings.Compare(ssi[len(ssi)-1], ssj[len(ssj)-1]) < 0
	})
	eventDetail := &EventDetail{VideoSeries: []*VideoSeries{}}
	firstHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "00"), EndTime: fmt.Sprintf("%d:%s", hour,
		"15")}
	secondHalfVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "15"), EndTime: fmt.Sprintf("%d:%s", hour,
		"30")}
	thirdVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "30"), EndTime: fmt.Sprintf("%d:%s", hour,
		"45")}
	fourthVideo := &VideoSeries{StartTime: fmt.Sprintf("%d:%s", hour, "45"), EndTime: fmt.Sprintf("%d:%s", hour+1,
		"00")}
	for index := range allLinks {
		isCollected := false
		collects, err := s.CollectDao.Gets(&model.Collect{OpenID: openID, Status: 1, FileID: allLinks[index]})
		if err != nil {
			return nil, err
		}
		if len(collects) > 0 {
			isCollected = true
		}
		links := strings.Split(allLinks[index], "/")
		minuteString := strings.Split(strings.Split(links[len(links)-1], "-")[1], ".")[0]
		minute, _ := strconv.Atoi(minuteString)
		if minute <= 15 {
			firstHalfVideo.Videos = append(firstHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else if minute > 15 && minute <= 30 {
			secondHalfVideo.Videos = append(secondHalfVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else if minute > 30 && minute <= 45 {
			thirdVideo.Videos = append(thirdVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else {
			fourthVideo.Videos = append(fourthVideo.Videos, &Video{
				IsCollected: isCollected,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		}
	}
	if len(firstHalfVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 25 && len(secondHalfVideo.Videos) == 0 {
			firstHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, firstHalfVideo)
	}
	if len(secondHalfVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 40 && len(thirdVideo.Videos) == 0 {
			secondHalfVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, secondHalfVideo)
	}
	if len(thirdVideo.Videos) > 0 {
		if time.Now().Hour() == hour && time.Now().Minute() < 55 && len(fourthVideo.Videos) == 0 {
			thirdVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, thirdVideo)
	}
	if len(fourthVideo.Videos) > 0 {
		if time.Now().Hour() == hour || (time.Now().Hour() == hour+1 && time.Now().Minute() < 10) {
			fourthVideo.Status = 1
		}
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, fourthVideo)
	}
	return eventDetail, nil
}
