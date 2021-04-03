package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"errors"
	"fmt"
)

type dbUserRepository struct {
	db *database.DBManager
}

func NewUserRepository() domain.UserRepository {
	return &dbUserRepository{db: database.Manager}
}

func (ur *dbUserRepository) GetById(id uint) (domain.User, error) {
	data, err := ur.db.Query("select id, username, email, avatar from users where id = $1", id)
	if err != nil {
		return domain.User{}, err
	}
	first := data[0]
	if first == nil {
		return domain.User{}, errors.New("User does not exist")
	}
	fmt.Println(first)
	user := first.(domain.User)
	return user, nil
}

func (ur *dbUserRepository) GetByEmail(email string) (domain.User, error) {
	data, err := ur.db.Query("select id, username, email, avatar from users where email = $1", email)
	if err != nil {
		return domain.User{}, err
	}
	first := data[0]
	if first == nil {
		return domain.User{}, errors.New("User does not exist")
	}
	fmt.Println(first)
	user := first.(domain.User)
	return user, nil
}

func (ur *dbUserRepository) Update(user *domain.User) error {
	return ur.db.Exec("update users set username = $1, email = $2, avatar = $3 where id = $4;", user.Username, user.Email, user.Avatar, user.ID)
}

func (ur *dbUserRepository) Store(user *domain.User) (uint, error) {
	err := ur.db.Exec("insert into users values(default, $1, $2, $3, $4, false);", user.Username, user.Email, user.Password, user.Avatar)
	if err != nil {
		return 0, err
	}
	newUser, err := ur.GetByEmail(user.Email)
	return newUser.ID, nil
}

func (ur *dbUserRepository) Delete(id uint) error {
	return ur.db.Exec("delete from users where id = $1;", id)
}
