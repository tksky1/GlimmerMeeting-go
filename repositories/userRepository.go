package repositories

import (
	"GlimmerMeeting/models"
	"gorm.io/gorm"
)

var userRepo *UserRepository

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	if userRepo == nil {
		userRepo = &UserRepository{
			db: MeetingDB,
		}
	}
	return userRepo
}

func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "username = ?", username)
	return &user, result.Error
}

func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *UserRepository) Delete(user *models.User) error {
	result := r.db.Delete(user)
	return result.Error
}

func (r *UserRepository) List() (users []models.User, err error) {
	err = r.db.Find(&users).Error
	return users, err
}
