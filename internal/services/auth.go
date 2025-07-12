package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"my-go-backend/pkg/models"
	"time"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Register(req *models.CreateUserRequest) (*models.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
