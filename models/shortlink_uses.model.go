package models

import (
	"shorts/database"
	"time"
)

// ShortlinkUse structure
type ShortlinkUse struct {
	ID      uint64    `json:"-" gorm:"primary_key"`
	LinkID  uint64    `json:"-" gorm:"not null"`
	UseTime time.Time `json:"time" gorm:"not null"`
}

// UseCount : returns uses count of each full lunk
func (shortlinkUse ShortlinkUse) UseCount() ([]FullLinkUseCountResponse, error) {
	var linksUses []FullLinkUseCountResponse
	if dbc := database.DB.Table("shortlink_uses").Select("shortlinks.full as full_link, count(1) as uses_count").Group("full_link").
		Joins("left join shortlinks on shortlinks.id = shortlink_uses.link_id").Scan(&linksUses); dbc.Error != nil {
		return nil, dbc.Error
	}

	return linksUses, nil
}
