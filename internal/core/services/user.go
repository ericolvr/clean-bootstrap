package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ericolvr/sec-backend/internal/core/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	if err := user.ValidateUser(); err != nil {
		return err
	}
	if err := s.validateUniqueFields(ctx, user); err != nil {
		return err
	}

	// Generate random password
	var plainPassword string
	if user.Password == "" {
		randomPassword, err := domain.GenerateRandomPassword()
		if err != nil {
			return errors.New("failed to generate random password: " + err.Error())
		}
		plainPassword = randomPassword
		fmt.Printf("🔐 Random password generated for %s: %s\n", user.Name, plainPassword)
	} else {
		plainPassword = user.Password
		fmt.Printf("🔐 Provided password for %s: %s\n", user.Name, plainPassword)
	}

	// make password hash before save on database
	hashedPassword, err := domain.HashPassword(plainPassword)
	if err != nil {
		return errors.New("erro ao fazer hash da senha: " + err.Error())
	}
	user.Password = hashedPassword
	user.Status = true

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.List(ctx, limit, offset)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("ID is required")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("ID must be a valid number")
	}

	return s.userRepo.GetByID(ctx, idInt)
}

func (s *UserService) GetByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	if mobile == "" {
		return nil, errors.New("celular é obrigatório")
	}

	return s.userRepo.GetByMobile(ctx, mobile)
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	if user.ID <= 0 {
		return errors.New("ID is required for update")
	}

	if err := user.ValidateUserUpdate(); err != nil {
		return err
	}

	// Verificar se o usuário existe
	exists, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	if err := s.validateUniqueFields(ctx, exists); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("ID must be a valid number")
	}

	user, err := s.userRepo.GetByID(ctx, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.userRepo.Delete(ctx, idInt); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) validateUniqueFields(ctx context.Context, user *domain.User) error {
	existing, err := s.userRepo.GetByMobile(ctx, user.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// If mobile exists but belongs to the same user being updated, it's OK
	if existing != nil && existing.ID != user.ID {
		return fmt.Errorf("already exists a user with this mobile")
	}
	return nil
}
