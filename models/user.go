package models

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Password       string `gorm:"column:password" json:"password"`
	RealName       string `gorm:"column:real_name" json:"real_name"`
	RelatedBooking string `gorm:"column:related_booking" json:"related_booking"`
	RelatedMeeting string `gorm:"column:related_meeting" json:"related_meeting"`
	Username       string `gorm:"column:username" json:"username"`
}

func (u *User) TableName() string {
	return "user_info"
}
