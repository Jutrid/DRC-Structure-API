package models

type Province struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" gorm:"unique;not null"`
	Code          *string `json:"code"`
	PrincipalTown string `json:"principal_town"`
	Surface       string `json:"surface"`
	Population    string `json:"population"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`

	Cities      []City      `json:"-" gorm:"foreignKey:ProvinceID"`
	Territories []Territory `json:"-" gorm:"foreignKey:ProvinceID"`
}

// TableName overrides default
func (Province) TableName() string {
	return "provinces"
}
