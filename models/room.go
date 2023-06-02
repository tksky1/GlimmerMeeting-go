package models

type Room struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Info           string `gorm:"size:1000" json:"info"`
	Location       string `gorm:"size:255" json:"location"`
	Name           string `gorm:"size:255" json:"name"`
	RelatedBooking string `gorm:"size:1000" json:"related_booking"`
	StatID         int    `json:"stat_id"`
}

func (r *Room) TableName() string {
	return "room_info"
}
