package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"soalhub/pkg/utils"
)

type UserService interface {
	GetAllUsers() ([]UserResponse, error)
	Login(req LoginRequest) (*LoginResponse, error)
	RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, error)
	Register(req RegisterRequest) (*UserResponse, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	userResponses := ToUserResponseList(users)
	return userResponses, nil
}

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
