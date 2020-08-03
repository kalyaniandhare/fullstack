package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kalyaniandhare/fullstack/api/models"
)

var users = []models.LogConfig{
	models.LogConfig{
		LogLevel: "DEBUG",
		Interval:  500,
		FilePath: "/home/app.log",
	},
	models.LogConfig{
		LogLevel: "DEBUG",
		Interval:  300,
		FilePath: "/home/heroku/app.log",
	},
	models.LogConfig{
		LogLevel: "DEBUG",
		Interval:  300,
		FilePath: "/home/heroku/app.log",
	},
}

var posts = []models.Post{
	models.Post{
		AlertLogLevel:   "ERROR",
		DateTime: "2020/31/07",
		AlertMessage: "First log",
	},
	models.Post{
		AlertLogLevel:   "DEBUG",
		DateTime: "2020/01/08",
		AlertMessage: "Second on 1st aug log",
	},
	models.Post{
		AlertLogLevel:   "DEBUG",
		DateTime: "2020/01/08",
		AlertMessage: "Second on 1st aug log",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.LogConfig{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.LogConfig{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	//err = db.Debug().Model(&models.Post{}).AddForeignKey("logconfig_id", "users(id)", "cascade", "cascade").Error
	//if err != nil {
	//	log.Fatalf("attaching foreign key error: %v", err)
	//}

	for i, _ := range posts {
		err = db.Debug().Model(&models.LogConfig{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].LogConfigID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}