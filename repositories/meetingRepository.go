package repositories

import (
	"GlimmerMeeting/models"
	"gorm.io/gorm"
)

type MeetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository() MeetingRepository {
	return MeetingRepository{MeetingDB}
}

func (r *MeetingRepository) Create(meeting *models.MeetingRecord) error {
	return r.db.Create(meeting).Error
}

func (r *MeetingRepository) GetById(id uint) (*models.MeetingRecord, error) {
	var meeting models.MeetingRecord
	err := r.db.First(&meeting, id).Error
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (r *MeetingRepository) GetByBookerID(bookerID int) (list []*models.MeetingRecord, err error) {
	result := MeetingDB.Where("booker_id = ?", bookerID).Find(&list)
	err = result.Error
	return list, err
}

func (r *MeetingRepository) GetByDayAndRoomID(day string, roomID int) (list []*models.MeetingRecord, err error) {
	result := MeetingDB.Where("day = ? AND room_id = ?", day, roomID).Find(&list)
	err = result.Error
	return list, err
}

func (r *MeetingRepository) Update(meeting *models.MeetingRecord) error {
	return r.db.Save(meeting).Error
}

func (r *MeetingRepository) DeleteById(id uint) error {
	return r.db.Delete(&models.MeetingRecord{}, id).Error
}

func (r *MeetingRepository) List() (meetings []models.MeetingRecord, err error) {
	err = r.db.Find(&meetings).Error
	return meetings, err
}
