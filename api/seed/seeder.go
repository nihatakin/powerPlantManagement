package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/nihatakin/powerPlantManagement/api/models"
)

var users = []models.User{
	models.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: "password",
	},
	models.User{
		Name:     "Nihat AKIN",
		Email:    "akin.nht@gmail.com",
		Password: "password",
	},
	models.User{
		Name:     "Leo AKIN",
		Email:    "leo.akin@gmail.com",
		Password: "password",
	},
}

var powerPlants = []models.PowerPlant{
	models.PowerPlant{
		Name:           "Bandirma",
		ShortName:      "BND1",
		EtsoCode:       "40W000000024604G",
		InstalledPower: 930.80,
		Order:          1,
	},
	models.PowerPlant{
		Name:           "Bandirma2",
		ShortName:      "BND2",
		EtsoCode:       "40W0000031943675",
		InstalledPower: 607.20,
		Order:          2,
	},
	models.PowerPlant{
		Name:           "Kandil",
		ShortName:      "KND",
		EtsoCode:       "40W0000004120593",
		InstalledPower: 207.92,
		Order:          3,
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.PowerPlant{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.PowerPlant{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i, _ := range powerPlants {
		powerPlants[i].CreatorUserId = users[0].ID
		err = db.Debug().Model(&models.PowerPlant{}).Create(&powerPlants[i]).Error
		if err != nil {
			log.Fatalf("cannot seed power plants table: %v", err)
		}
	}
}
