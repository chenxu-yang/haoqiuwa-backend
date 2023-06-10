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
	EventDao *model.Event
	VideoDao *model.Video
	CourtDao *model.Court
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
	Status      int32          `json:"status"`
	Hour        int32          `json:"hour"`
	VideoSeries []*VideoSeries `json:"video_series"`
}

type VideoSeries struct {
	StartMinute int32    `json:"start_minute"`
	EndMinute   int32    `json:"end_minute"`
	Videos      []*Video `json:"videos"`
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
		log.Println(err)
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
	// sort hours
	sort.Slice(hours, func(i, j int) bool { return hours[i] > hours[j] })
	// get events by hours
	for _, hour := range hours {
		results = append(results, Event{StartTime: int32(hour), EndTime: int32(hour + 1), CourtName: courtID, Status: 0})
	}
	if distinctHours[hours[len(hours)-1]] < 6 {
		results[len(results)-1].Status = 1
	}
	return results, nil
}

func (s *Service) GetEventInfo(courtID string, hour int) (*EventDetail, error) {
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
	firstHalfVideo := &VideoSeries{}
	secondHalfVideo := &VideoSeries{}
	for index := range allLinks {
		//links := strings.Split(link, "/")
		//minuteString := strings.Split(strings.Split(links[len(links)-1], "-")[1], ".")[0]
		//minute, _ := strconv.Atoi(minuteString)
		if index <= 2 {
			firstHalfVideo.StartMinute = 0
			firstHalfVideo.EndMinute = 30
			firstHalfVideo.Videos = append(firstHalfVideo.Videos, &Video{
				IsCollected: false,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		} else {
			secondHalfVideo.StartMinute = 30
			secondHalfVideo.EndMinute = 60
			secondHalfVideo.Videos = append(secondHalfVideo.Videos, &Video{
				IsCollected: false,
				Url:         allLinks[index],
				PicUrl:      picLinks[index],
			})
		}
	}
	eventDetail.Hour = int32(hour)
	if len(allLinks) < 6 {
		eventDetail.Status = 1
	}
	if len(firstHalfVideo.Videos) > 0 {
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, firstHalfVideo)
	}
	if len(secondHalfVideo.Videos) > 0 {
		eventDetail.VideoSeries = append(eventDetail.VideoSeries, secondHalfVideo)
	}
	return eventDetail, nil
}
