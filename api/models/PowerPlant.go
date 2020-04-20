package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type PowerPlant struct {
	ID                   uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name                 string    `gorm:"size:255;not null;unique" json:"name"`
	ShortName            string    `gorm:"size:10;not null;unique" json:"short_name"`
	EtsoCode             string    `gorm:"size:255;not null;" json:"etso_code"`
	InstalledPower       float64   `gorm:"not null" json:"installed_power"`
	Order                uint32    `gorm:"not null" json:"order"`
	CreatorUser          User      `json:"creator_user"`
	CreationTime         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"creation_time"`
	CreatorUserId        uint32    `gorm:"not null" json:"creator_user_id"`
	ModifierUser         User      `json:"modifier_user"`
	LastModificationTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_modification_time"`
	LastModifierUserId   uint32    `gorm:json:"last_modifier_user_id"`
}

func (p *PowerPlant) PrepareForCreate() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.ShortName = html.EscapeString(strings.TrimSpace(p.ShortName))
	p.EtsoCode = html.EscapeString(strings.TrimSpace(p.EtsoCode))
	p.InstalledPower = p.InstalledPower
	p.Order = p.Order
	p.CreatorUser = User{}
	p.CreationTime = time.Now()
	p.LastModificationTime = time.Now()
}

func (p *PowerPlant) PrepareForUpdate() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.ShortName = html.EscapeString(strings.TrimSpace(p.ShortName))
	p.EtsoCode = html.EscapeString(strings.TrimSpace(p.EtsoCode))
	p.InstalledPower = p.InstalledPower
	p.Order = p.Order
	p.ModifierUser = User{}
	p.LastModificationTime = time.Now()
}

func (p *PowerPlant) Validate() error {

	if p.Name == "" {
		return errors.New("Required Name")
	}
	if p.ShortName == "" {
		return errors.New("Required ShortName")
	}
	if p.EtsoCode == "" {
		return errors.New("Required EtsoCode")
	}
	if p.InstalledPower < 1 {
		return errors.New("Required InstalledPower")
	}
	if p.Order < 1 {
		return errors.New("Required Order")
	}
	return nil
}

func (p *PowerPlant) SavePowerPlant(db *gorm.DB) (*PowerPlant, error) {
	var err error
	err = db.Debug().Model(&PowerPlant{}).Create(&p).Error
	if err != nil {
		return &PowerPlant{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorUserId).Take(&p.CreatorUser).Error
		if err != nil {
			return &PowerPlant{}, err
		}
	}
	return p, nil
}

func (p *PowerPlant) FindAllPowerPlants(db *gorm.DB) (*[]PowerPlant, error) {
	var err error
	powerPlants := []PowerPlant{}
	err = db.Debug().Model(&PowerPlant{}).Limit(100).Find(&powerPlants).Error
	if err != nil {
		return &[]PowerPlant{}, err
	}
	if len(powerPlants) > 0 {
		for i, _ := range powerPlants {
			err := db.Debug().Model(&User{}).Where("id = ?", powerPlants[i].CreatorUserId).Take(&powerPlants[i].CreatorUser).Error
			if err != nil {
				return &[]PowerPlant{}, err
			}

			if powerPlants[i].LastModifierUserId > 0 {
				err2 := db.Debug().Model(&User{}).Where("id = ?", powerPlants[i].LastModifierUserId).Take(&powerPlants[i].ModifierUser).Error
				if err2 != nil {
					return &[]PowerPlant{}, err2
				}
			}
		}
	}
	return &powerPlants, nil
}

func (p *PowerPlant) FindPowerPlantByID(db *gorm.DB, pid uint64) (*PowerPlant, error) {
	var err error
	err = db.Debug().Model(&PowerPlant{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &PowerPlant{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorUserId).Take(&p.CreatorUser).Error
		if err != nil {
			return &PowerPlant{}, err
		}
	}
	return p, nil
}

func (p *PowerPlant) UpdateAPowerPlant(db *gorm.DB) (*PowerPlant, error) {

	var err error

	err = db.Debug().Model(&PowerPlant{}).Where("id = ?", p.ID).Updates(PowerPlant{Name: p.Name, ShortName: p.ShortName, EtsoCode: p.EtsoCode, InstalledPower: p.InstalledPower, Order: p.Order, LastModifierUserId: p.LastModifierUserId, LastModificationTime: time.Now()}).Error
	if err != nil {
		return &PowerPlant{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.LastModifierUserId).Take(&p.ModifierUser).Error
		if err != nil {
			return &PowerPlant{}, err
		}
	}
	return p, nil
}

func (p *PowerPlant) DeleteAPowerPlant(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&PowerPlant{}).Where("id = ? and author_id = ?", pid, uid).Take(&PowerPlant{}).Delete(&PowerPlant{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("PowerPlant not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
