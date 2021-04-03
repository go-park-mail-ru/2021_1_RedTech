package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"errors"
)

type dbUserRepository struct {
	db *database.DBManager
}

func NewUserRepository(db *database.DBManager) domain.UserRepository {
	return &dbUserRepository{db: db}
}

func (ur *dbUserRepository) GetById(id uint) (domain.User, error) {
	data, err := ur.db.Query("select id, username, email, avatar from users where id = $1", id)
	if err != nil {
		return domain.User{}, err
	}
	first, ok := data[0].([]interface{})
	if first == nil || !ok {
		return domain.User{}, errors.New("User does not exist")
	}

	user := domain.User{
		ID:       uint(first[0].(int32)),
		Username: first[1].(string),
		Email:    first[2].(string),
		Avatar:   first[3].(string),
	}
	return user, nil
}

func (ur *dbUserRepository) GetByEmail(email string) (domain.User, error) {
	data, err := ur.db.Query("select id, username, email, avatar, password from users where email = $1", email)
	if err != nil {
		return domain.User{}, err
	}
	first, ok := data[0].([]interface{})
	if first == nil || !ok {
		return domain.User{}, errors.New("User does not exist")
	}

	var password [domain.HashLen]byte
	copy(password[:], first[4].([]byte))
	user := domain.User{
		ID:       uint(first[0].(int32)),
		Username: first[1].(string),
		Email:    first[2].(string),
		Avatar:   first[3].(string),
		Password: password,
	}
	return user, nil
}

func (ur *dbUserRepository) Update(user *domain.User) error {
	return ur.db.Exec("update users set username = $1, email = $2, avatar = $3 where id = $4;", user.Username, user.Email, user.Avatar, user.ID)
}

func (ur *dbUserRepository) Store(user *domain.User) (uint, error) {
	data, err := ur.db.Query("insert into users values(default, $1, $2, $3, $4, false) returning id;", user.Username, user.Email, user.Password[:], user.Avatar)
	if err != nil {
		return 0, err
	}
	first, ok := data[0].([]interface{})
	if first == nil || !ok {
		return 0, errors.New("Cannot create user in database")
	}
	return uint(first[0].(int32)), nil
}

func (ur *dbUserRepository) Delete(id uint) error {
	return ur.db.Exec("delete from users where id = $1;", id)
}
