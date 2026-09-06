package user

import "gorm.io/gorm"

type UserRepository interface {
	FindAll(param UserQueryParam) ([]User, int64, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// =============== User Management ===============
func (r *userRepository) FindAll(param UserQueryParam) ([]User, int64, error) {
	var users []User
	var totalItems int64

	query := r.db.Model(&User{})

	if param.ID > 0 {
		query = query.Where("id = ?", param.ID)
	}

	if param.Name != "" {
		query = query.Where("name ILIKE ?", "%"+param.Name+"%")
	}

	if param.Email != "" {
		query = query.Where("email ILIKE ?", "%"+param.Email+"%")
	}

	if param.Role != "" {
		query = query.Where("role = ?", param.Role)
	}

	// Get total count
	err := query.Count(&totalItems).Error
	if err != nil {
		return nil, 0, err
	}

	page := param.Page
	if page <= 0 {
		page = 1
	}

	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	err = query.Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
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

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *User) error {
	return r.db.Save(user).Error
}
