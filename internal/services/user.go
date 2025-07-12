package services

import (
	"errors"
	"gorm.io/gorm"
	"my-go-backend/pkg/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) GetAllUsers(page, limit int) (*models.PaginatedResponse, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data:       userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}
