package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"soalhub/pkg/utils"
)

type UserService interface {
	Login(req LoginRequest) (*LoginResponse, error)
	RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, error)
	Register(req RegisterRequest) (*UserResponse, error)
	GetAllUsers(param UserQueryParam) (*UserListResponse, error)
	GetUserByID(targetID uint, userID uint, role string) (*UserResponse, error)
	UpdateUser(targetID uint, userID uint, role string, req UpdateUserRequest) (*UserResponse, error)
	DeleteUser(targetID uint, userID uint, role string) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// ================ Authentication ================
func (s *userService) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// 3. Generate JWT Token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to create access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to create refresh token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         ToUserResponse(*user),
	}, nil
}

func (s *userService) RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, errors.New("failed to create new access token")
	}

	// Generate new refresh token
	newRefreshToken, err := utils.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, errors.New("failed to create new refresh token")
	}

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userService) Register(req RegisterRequest) (*UserResponse, error) {
	// Check if user with the same email already exists
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	response := ToUserResponse(*user)
	return &response, nil
}

// ================ User Management ================
func (s *userService) GetAllUsers(param UserQueryParam) (*UserListResponse, error) {
	users, totalItems, err := s.repo.FindAll(param)
	if err != nil {
		return nil, err
	}

	page := param.Page
	if page <= 0 {
		page = 1
	}

	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	if totalPages > 0 && page > int(totalPages) {
		return nil, errors.New("page number exceeds total pages")
	}

	return &UserListResponse{
		Users: ToUserResponseList(users),
		Meta: MetaPagination{
			CurrentPage: page,
			TotalPages:  int(totalPages),
			Limit:       limit,
			TotalItems:  int(totalItems),
		},
	}, nil
}

func (s *userService) GetUserByID(targetID uint, userID uint, role string) (*UserResponse, error) {
	if role != "admin" && targetID != userID {
		return nil, errors.New("unauthorized access")
	}

	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := ToUserResponse(*user)
	return &response, nil
}

func (s *userService) UpdateUser(targetID uint, userID uint, role string, req UpdateUserRequest) (*UserResponse, error) {
	if role != "admin" && targetID != userID {
		return nil, errors.New("unauthorized access")
	}

	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if role != "admin" && req.Role != user.Role {
		return nil, errors.New("unauthorized to change role")
	}

	if req.Email != user.Email {
		existingUser, _ := s.repo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != targetID {
			return nil, errors.New("email already registered")
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	err = s.repo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update user")
	}

	response := ToUserResponse(*user)
	return &response, nil
}

func (s *userService) DeleteUser(targetID uint, userID uint, role string) error {
	if role != "admin" && targetID != userID {
		return errors.New("unauthorized access")
	}

	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.repo.Delete(user.ID)
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
