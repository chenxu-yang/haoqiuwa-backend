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
	Status int32    `json:"status"`
	Hour   int32    `json:"hour"`
	Videos []*Video `json:"videos"`
}

type Video struct {
	StartMinute int32  `json:"start_minute"`
	EndMinute   int32  `json:"end_minute"`
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
	distinctHours := make(map[int]bool)
	for _, link := range allLinks {
		links := strings.Split(link, "/")
		hour := strings.Split(links[len(links)-1], "-")[0]
		hourInt, _ := strconv.Atoi(hour[1:])
		distinctHours[hourInt] = true
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
	// if last hour is 5min before time now, set status to 1, if now is 10-03, last hour is 10-00, set status to 1
	if len(hours) > 0 && time.Now().Hour() == hours[len(hours)-1]+1 && time.Now().Minute() < 6 {
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
	eventDetail := &EventDetail{Videos: []*Video{}}
	for _, link := range allLinks {
		links := strings.Split(link, "/")
		minuteString := strings.Split(strings.Split(links[len(links)-1], "-")[1], ".")[0]
		minute, _ := strconv.Atoi(minuteString)
		eventDetail.Videos = append(eventDetail.Videos, &Video{
			StartMinute: int32(minute),
			EndMinute:   int32(minute + 5),
			Url:         link,
			PicUrl: "cloud://prod-2gicsblt193f5dc8.7072-prod-2gicsblt193f5dc8-1318337180/highlight/court10/testcover." +
				"png",
		})
	}
	eventDetail.Hour = int32(hour)
	if time.Now().Hour() == hour+1 && time.Now().Minute() < 6 {
		eventDetail.Status = 1
	}
	return eventDetail, nil
}
