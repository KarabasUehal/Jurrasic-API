package models

type Dinosaurus struct {
	ID      int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Species string  `json:"species" gorm:"type:text;not null"`
	Types   string  `json:"types" gorm:"types:text;not null"`
	Height  float64 `json:"height" gorm:"type:double precision;not null"`
	Length  float64 `json:"length" gorm:"type:double precision;not null"`
	Weight  float64 `json:"weight" gorm:"type:double precision;not null"`
	Aquatic bool    `json:"aquatic" gorm:"type:boolean;not null"`
	Flying  bool    `json:"flying" gorm:"type:boolean;not null"`
}
