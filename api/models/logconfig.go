package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)
type LogConfig struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	LogLevel     string    `gorm:"size:255;not null" json:"log_level"`
	Interval   int    `gorm:"size:255;not null;" json:"interval"`
	FilePath   string      `gorm:"size:255;not null" json:"file_path"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *LogConfig) Prepare() {
	p.ID = 0
	p.LogLevel = html.EscapeString(strings.TrimSpace(p.LogLevel))
	p.FilePath = html.EscapeString(strings.TrimSpace(p.FilePath))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *LogConfig) Validate() error {

	if p.LogLevel == "" {
		return errors.New("Required Title")
	}
	if p.FilePath == "" {
		return errors.New("Required Content")
	}
	if p.Interval < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *LogConfig) SaveLogConfig(db *gorm.DB) (*LogConfig, error) {
	var err error
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &LogConfig{}, err
	}

	return p, nil
}


func (p *LogConfig) FindAllLogs(db *gorm.DB) (*[]LogConfig, error) {
	var err error
	posts := []LogConfig{}
	err = db.Debug().Model(&LogConfig{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]LogConfig{}, err
	}
	return &posts, nil
}
func (p *Post) FindLogByID(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&LogConfig{}).Where("id = ?", p.LogConfigID).Take(&p.LogConfig).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}