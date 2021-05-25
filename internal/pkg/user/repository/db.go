package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/cast"
	"Redioteka/internal/pkg/utils/log"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	querySelectID         = "select id, username, email, avatar from users where id = $1;"
	querySelectEmail      = "select id, username, email, avatar, password from users where email = $1;"
	queryInsert           = "insert into users values(default, $1, $2, $3, $4, false) returning id;"
	queryDelete           = "delete from users where id = $1;"
	querySelectFavourites = `select distinct m.id, m.title, m.avatar, m.rating, m.is_free 
							from movies as m join user_favs as uf on m.id = uf.movie_id
							join users as u on u.id = uf.user_id 
							where u.id = $1;`
	queryUpdate    = "update users set username = $1, email = $2, avatar = $3 where id = $4"
	querySelectSub = "select expires, actual from subscriptions where user_id = $1;"
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
		log.Log.Warn(fmt.Sprint("Cannot get user from db with id: ", id))
		return domain.User{}, err
	}
	if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("User with id: %d - not found in db", id))
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
		log.Log.Warn(fmt.Sprint("Cannot get user from db with email: ", email))
		return domain.User{}, err
	}
	if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("User with email: %s - not found in db", email))
		return domain.User{}, errors.New("User does not exist")
	}

	first := data[0]
	user := domain.User{
		ID:       cast.ToUint(first[0]),
		Username: cast.ToString(first[1]),
		Email:    cast.ToString(first[2]),
		Avatar:   cast.ToString(first[3]),
		Password: first[4],
	}
	return user, nil
}

func buildUpdateQuery(update *domain.User) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updateQuery := psql.Update("users").Where(sq.Eq{"id": update.ID})
	if update.Username != "" {
		updateQuery = updateQuery.Set("username", update.Username)
	}
	if update.Email != "" {
		updateQuery = updateQuery.Set("email", update.Email)
	}
	if update.Avatar != "" {
		updateQuery = updateQuery.Set("avatar", update.Avatar)
	}
	return updateQuery.ToSql()
}

func (ur *dbUserRepository) Update(user *domain.User) error {
	updateQuery, params, err := buildUpdateQuery(user)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Can't construct user %v update query: %v", user.ID, err))
		return err
	}
	err = ur.db.Exec(updateQuery, params...)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot update user in db with id: ", user.ID))
	}
	return err
}

func (ur *dbUserRepository) Store(user *domain.User) (uint, error) {
	data, err := ur.db.Query(queryInsert, user.Username, user.Email, user.Password[:], user.Avatar)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot insert user in db with username: %s email: %s", user.Username, user.Email))
		return 0, err
	}
	if len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("No id was returned by inserting user with username: %s email: %s", user.Username, user.Email))
		return 0, errors.New("Cannot create user in database")
	}
	return cast.ToUint(data[0][0]), nil
}

func (ur *dbUserRepository) Delete(id uint) error {
	err := ur.db.Exec(queryDelete, id)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot delete user in db with id: ", id))
	}
	return err
}

func (ur *dbUserRepository) GetFavouritesByID(id uint) ([]domain.Movie, error) {
	data, err := ur.db.Query(querySelectFavourites, id)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot get favourites of user with id: %d", id))
		return nil, user.NotFoundError
	}

	result := make([]domain.Movie, 0)
	for _, movie := range data {
		result = append(result, domain.Movie{
			ID:     cast.ToUint(movie[0]),
			Title:  cast.ToString(movie[1]),
			Avatar: cast.ToString(movie[2]),
			Rating: cast.ToFloat(movie[3]),
			IsFree: cast.ToBool(movie[4]),
		})
	}
	return result, nil
}

func (ur *dbUserRepository) CheckSub(id uint) time.Time {
	data, err := ur.db.Query(querySelectSub, id)
	if err != nil || len(data) == 0 {
		log.Log.Warn(fmt.Sprintf("Cannot get sub of user with id: %d", id))
		return time.Now()
	}
	if !cast.ToBool(data[0][1]) {
		log.Log.Info("Sub is not actual")
		return time.Now()
	}
	sec := cast.ToInt(data[0][0])
	return time.Unix(int64(sec), 0)
}
