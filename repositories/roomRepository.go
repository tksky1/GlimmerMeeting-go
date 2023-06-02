package repositories

import (
	"GlimmerMeeting/models"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository() RoomRepository {
	return RoomRepository{MeetingDB}
}

func (r *RoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *RoomRepository) GetById(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) DeleteById(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *RoomRepository) List() (rooms []models.Room, err error) {
	err = r.db.Find(&rooms).Error
	return rooms, err
}
