package models

type City struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	ProvinceID uint    `json:"province_id" gorm:"index;not null"`
	Name       string  `json:"name" gorm:"not null"`
	Code       *string `json:"code"`

	Province *Province `json:"province,omitempty" gorm:"foreignKey:ProvinceID"`
}

func (City) TableName() string {
	return "cities"
}
