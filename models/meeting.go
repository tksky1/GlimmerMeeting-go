package models

import "encoding/json"

type Timepiece struct {
	BeginHour   int `json:"beginhour"`
	BeginMinute int `json:"beginminute"`
	EndHour     int `json:"endhour"`
	EndMinute   int `json:"endminute"`
}

type Meeting struct {
	ID           int       `json:"meetingid"`
	Attendees    []int     `json:"attendees"`
	BookerID     int       `json:"bookerid"`
	Duration     Timepiece `json:"duration"`
	RoomID       int       `json:"roomid"`
	RoomLocation string    `json:"roomlocation"`
	Theme        string    `json:"theme"`
	Day          string    `json:"day"`
}

// ToMeetingRecord 从Meeting转换到MeetingRecord
func (m *Meeting) ToMeetingRecord() *MeetingRecord {
	attendees, _ := json.Marshal(m.Attendees)
	return &MeetingRecord{
		BeginHour:    m.Duration.BeginHour,
		BeginMinute:  m.Duration.BeginMinute,
		EndHour:      m.Duration.EndHour,
		EndMinute:    m.Duration.EndMinute,
		Attendees:    string(attendees),
		BookerID:     m.BookerID,
		RoomID:       m.RoomID,
		Theme:        m.Theme,
		RoomLocation: m.RoomLocation,
		ID:           uint(m.ID),
		Day:          m.Day,
	}
}

type MeetingRecord struct {
	ID           uint   `gorm:"primaryKey;column:meeting_id"`
	BeginHour    int    `gorm:"not null"`
	BeginMinute  int    `gorm:"not null"`
	EndHour      int    `gorm:"not null"`
	EndMinute    int    `gorm:"not null"`
	Attendees    string `gorm:"not null"`
	BookerID     int    `gorm:"not null"`
	RoomID       int    `gorm:"not null"`
	Theme        string `gorm:"not null"`
	RoomLocation string `gorm:"not null"`
	Day          string `gorm:"not null"`
}

func (mr *MeetingRecord) ToMeeting() *Meeting {
	var attendees []int
	json.Unmarshal([]byte(mr.Attendees), &attendees)
	return &Meeting{
		ID:        int(mr.ID),
		Attendees: attendees,
		BookerID:  mr.BookerID,
		Duration: Timepiece{BeginHour: mr.BeginHour, BeginMinute: mr.BeginMinute,
			EndHour: mr.EndHour, EndMinute: mr.EndMinute},
		RoomID:       mr.RoomID,
		RoomLocation: mr.RoomLocation,
		Theme:        mr.Theme,
		Day:          mr.Day,
	}
}

func (*MeetingRecord) TableName() string {
	return "meeting"
}

type MeetingRecordList []*MeetingRecord

func (l MeetingRecordList) Len() int {
	return len(l)
}

func (l MeetingRecordList) Less(i, j int) bool {
	return l[i].BeginHour < l[j].BeginHour || (l[i].BeginHour == l[j].BeginHour && l[i].BeginMinute < l[j].BeginMinute)
}

func (l MeetingRecordList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func TimePieceConflict(i, j Timepiece) bool {
	start1 := i.BeginHour*60 + i.BeginMinute
	end1 := i.EndHour*60 + i.EndMinute
	start2 := j.BeginHour*60 + j.BeginMinute
	end2 := j.EndHour*60 + j.EndMinute

	if start1 < end2 && start2 < end1 {
		return true
	}
	return false
}
