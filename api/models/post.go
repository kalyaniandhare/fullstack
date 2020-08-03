package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	AlertLogLevel     string    `gorm:"size:255;not null;" json:"alert_log_level"`
	DateTime   string    `gorm:"size:255;not null;" json:"datetime"`
	AlertMessage   string      `gorm:"size:255;not null;" json:"alermessage"`
	LogConfig    LogConfig      `json:"logconfig"`
	LogConfigID uint64         `gorm:"not null" json:"logconfig"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.AlertLogLevel = html.EscapeString(strings.TrimSpace(p.AlertLogLevel))
	p.AlertMessage = html.EscapeString(strings.TrimSpace(p.AlertMessage))
	p.LogConfig = LogConfig{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {

	if p.AlertLogLevel == "" {
		return errors.New("Required Title")
	}
	if p.AlertMessage == "" {
		return errors.New("Required Content")
	}
	return nil
}

func (p *Post) SavePostNEW(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	return p, nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
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

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&LogConfig{}).Where("id = ?", posts[i].LogConfigID).Take(&posts[i].LogConfig).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
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
