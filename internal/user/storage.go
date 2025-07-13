package user

import (
	"errors"
	"sync"

	"github.com/emday4prez/fs-go/internal/domain"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("username is already taken")
)

type Storage interface {
	Save(user *domain.User) error
	FindByUsername(username string) (*domain.User, error)
}

type InMemoryStorage struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users: make(map[string]*domain.User),
	}
}

// adds a new user to the in memory map
func (s *InMemoryStorage) Save(user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Username]; exists {
		return ErrUsernameTaken
	}

	s.users[user.Username] = user
	return nil
}

func (s *InMemoryStorage) FindByUsername(username string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}
