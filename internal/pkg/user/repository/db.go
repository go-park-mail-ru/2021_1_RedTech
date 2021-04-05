package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/cast"
	"errors"
)

const (
	querySelectID    = "select id, username, email, avatar from users where id = $1;"
	querySelectEmail = "select id, username, email, avatar, password from users where email = $1;"
	queryUpdate      = "update users set username = $1, email = $2, avatar = $3 where id = $4;"
	queryInsert      = "insert into users values(default, $1, $2, $3, $4, false) returning id;"
	queryDelete      = "delete from users where id = $1;"
)

type dbUserRepository struct {
	db *database.DBManager
}

func NewUserRepository(db *database.DBManager) domain.UserRepository {
	return &dbUserRepository{db: db}
}

func (ur *dbUserRepository) GetById(id uint) (domain.User, error) {
	data, err := ur.db.Query(querySelectID, id)
	if err != nil {
		return domain.User{}, err
	}
	if data == nil {
		return domain.User{}, errors.New("User does not exist")
	}

	first := data[0]
	user := domain.User{
		ID:       cast.ToUint(first[0]),
		Username: cast.ToString(first[1]),
		Email:    cast.ToString(first[2]),
		Avatar:   cast.ToString(first[3]),
	}
	return user, nil
}

func (ur *dbUserRepository) GetByEmail(email string) (domain.User, error) {
	data, err := ur.db.Query(querySelectEmail, email)
	if err != nil {
		return domain.User{}, err
	}
	if data == nil {
		return domain.User{}, errors.New("User does not exist")
	}

	first := data[0]
	var password [domain.HashLen]byte
	copy(password[:], first[4])
	user := domain.User{
		ID:       cast.ToUint(first[0]),
		Username: cast.ToString(first[1]),
		Email:    cast.ToString(first[2]),
		Avatar:   cast.ToString(first[3]),
		Password: password,
	}
	return user, nil
}

func (ur *dbUserRepository) Update(user *domain.User) error {
	return ur.db.Exec(queryUpdate, user.Username, user.Email, user.Avatar, user.ID)
}

func (ur *dbUserRepository) Store(user *domain.User) (uint, error) {
	data, err := ur.db.Query(queryInsert, user.Username, user.Email, user.Password[:], user.Avatar)
	if err != nil {
		return 0, err
	}
	if data == nil {
		return 0, errors.New("Cannot create user in database")
	}
	return cast.ToUint(data[0][0]), nil
}

func (ur *dbUserRepository) Delete(id uint) error {
	return ur.db.Exec(queryDelete, id)
}
