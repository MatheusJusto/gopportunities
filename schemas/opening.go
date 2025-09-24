package schemas

import (
	"gorm.io/gorm"
)

type Opening struct {
	gorm.Model
	Role         string
	Compony      string
	Localization string
	Remote       bool
	Link         string
	Salary       int64
}
