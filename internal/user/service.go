package user

import (
	"github.com/emday4prez/fs-go/internal/domain"
	"github.com/google/uuid"
)

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) Register(username, password string) (*domain.User, error) {

	id := uuid.NewString()

	newUser, err := domain.NewUser(id, username, password)
	if err != nil {
		return nil, err
	}

	err = s.storage.Save(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *Service) Login(username, password string) (*domain.User, error) {
	user, err := s.storage.FindByUsername(username)
	if err != nil {

		return nil, err
	}

	if !user.CheckPassword(password) {
		// If the password doesn't match, return a generic error.
		return nil, ErrUserNotFound
	}

	return user, nil
}
