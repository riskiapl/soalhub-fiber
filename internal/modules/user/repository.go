package user

import "gorm.io/gorm"

type UserRepository interface {
	FindAll() ([]User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Create(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// =============== Authentication ===============
func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

// =============== User Management ===============
func (r *userRepository) FindAll() ([]User, error) {
	var users []User
	result := r.db.Find(&users).Error

	if result != nil {
		return nil, result
	}

	return users, nil
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	error := r.db.Where("email = ?", email).First(&user).Error

	if error != nil {
		return nil, error
	}

	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*User, error) {
	var user User
	error := r.db.First(&user, id).Error

	if error != nil {
		return nil, error
	}

	return &user, nil
}
