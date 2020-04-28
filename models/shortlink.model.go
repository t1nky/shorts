package models

import (
	h "shorts/helper"

	"github.com/jinzhu/gorm"
)

// Shortlink structure
type Shortlink struct {
	ID      uint64 `json:"id" gorm:"primary_key"`
	Short   string `json:"short" gorm:"unique;not null"`
	Full    string `json:"full" gorm:"not null"`
	OwnerID uint64 `json:"ownerId" gorm:"not null"`

	Uses []ShortlinkUse `gorm:"ForeignKey:LinkID" json:"uses"`
}

// AfterCreate for updating `short` field
func (s *Shortlink) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(s).Update("short", h.MakeShortlinkFromID(s.ID))
	return
}
