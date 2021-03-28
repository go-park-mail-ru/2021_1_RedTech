package repository

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"errors"
	"fmt"
	"sync"
)

type mapUserRepository struct {
	sync.Mutex
	users map[uint]*domain.User
}

func NewMapUserRepository() domain.UserRepository {
	return &mapUserRepository{
		users: make(map[uint]*domain.User),
	}
}

func (m *mapUserRepository) GetById(id uint) (domain.User, error) {
	m.Lock()
	user, inMap := m.users[id]
	m.Unlock()
	if !inMap {
		return domain.User{}, errors.New("not in map")
	}
	return *user, nil
}

func (m *mapUserRepository) GetByEmail(email string) (domain.User, error) {
	m.Lock()
	defer m.Unlock()
	for _, user := range m.users {
		if user.Email == email {
			return *user, nil
		}
	}
	return domain.User{}, errors.New("not in map")
}

func updateUser(user *domain.User, update *domain.User) {
	if update.Username != "" {
		user.Username = update.Username
	}
	if update.Email != "" {
		user.Email = update.Email
	}
	if update.Avatar != "" {
		user.Avatar = update.Avatar
	}
}

func (m *mapUserRepository) Update(update *domain.User) error {
	old, err := m.GetById(update.ID)
	if err != nil {
		return user.NotFoundError
	}
	updateUser(&old, update)

	err = m.Delete(update.ID)
	if err != nil {
		return fmt.Errorf("old user deleting error %s", err)
	}
	_, err = m.Store(&old)
	if err != nil {
		return fmt.Errorf("update user storing error: %s", err)
	}
	return nil
}

func (m *mapUserRepository) Store(user *domain.User) (uint, error) {
	m.Lock()
	defer m.Unlock()
	_, inMap := m.users[user.ID]
	if inMap {
		return 0, errors.New("user already in map")
	}
	// if uninitialized, create new
	if user.ID == 0 {
		user.ID = uint(len(m.users)) + 1
	}
	m.users[user.ID] = user
	return user.ID, nil
}

func (m *mapUserRepository) Delete(id uint) error {
	m.Lock()
	defer m.Unlock()
	_, inMap := m.users[id]
	if !inMap {
		return errors.New("user not in map")
	}
	delete(m.users, id)
	return nil
}
