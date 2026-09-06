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
	GetUserByID(userID uint) *UserResponse
	UpdateUser(userID uint, req UpdateUserRequest) (*UserResponse, error)
	DeleteUser(userID uint) error
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

	if claims.TokenTYpe != "refresh" {
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

	if page > int(totalPages) {
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

func (s *userService) GetUserByID(userID uint) *UserResponse {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil
	}

	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func (s *userService) UpdateUser(userID uint, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.repo.FindByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already registered")
		}
		user.Email = *req.Email
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update user")
	}

	response := ToUserResponse(*user)
	return &response, nil
}

func (s *userService) DeleteUser(userID uint) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.repo.Delete(user.ID)
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
