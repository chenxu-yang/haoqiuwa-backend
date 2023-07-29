package model

import (
	"time"
	"wxcloudrun-golang/internal/pkg/db"
)

type Video struct {
	ID          int32     `gorm:"primary_key" json:"id"`
	FilePath    string    `gorm:"column:file_path" json:"file_path"`
	Date        int32     `gorm:"column:date" json:"date"`
	Time        string    `gorm:"column:time" json:"time"`
	Type        int32     `gorm:"column:type" json:"type"`
	Court       int32     `gorm:"column:court" json:"court"`
	Hour        int32     `gorm:"column:hour" json:"hour"`
	FileName    string    `gorm:"column:file_name" json:"file_name"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
}

// GORM table name for Video struct
func (obj *Video) TableName() string {
	return "t_video"
}
func (obj *Video) Create(video *Video) (*Video, error) {
	err := db.Get().Create(video).Error
	return video, err
}

func (obj *Video) Get(video *Video) (*Video, error) {
	result := new(Video)
	err := db.Get().Table(obj.TableName()).Where(video).First(result).Error
	return result, err
}

func (obj *Video) Gets(video *Video) ([]Video, error) {
	results := make([]Video, 0)
	err := db.Get().Table(obj.TableName()).Where(video).Find(&results).Error
	return results, err
}

func (obj *Video) Update(video *Video) (*Video, error) {
	err := db.Get().Table(obj.TableName()).Where("id = ?", video.ID).Updates(video).Error
	return video, err
}

func (obj *Video) Delete(video *Video) error {
	return db.Get().Delete(video, "id = ?", video.ID).Error
}
