package repository

import (
	"Redioteka/internal/pkg/database"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"fmt"
)

const (
	insertSub = `insert into subscriptions values(default, $1, $2, $3) on conflict (user_id) do update set expires = $2, actual = $3;`
	deleteSub = `update subscriptions set actual = false where user_id = $1;`
)

type subscriptionRepository struct {
	db *database.DBManager
}

func NewSubscriptionRepository(db *database.DBManager) *subscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (sr *subscriptionRepository) Create(sub *domain.Subscription) error {
	err := sr.db.Exec(insertSub, sub.UserID, sub.Expiraton.Unix(), sub.Actual)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot add sub for user id: ", sub.UserID))
		return err
	}
	return nil
}

func (sr *subscriptionRepository) Delete(id uint) error {
	err := sr.db.Exec(deleteSub, id)
	if err != nil {
		log.Log.Warn(fmt.Sprint("Cannot delete sub for user id: ", id))
		return err
	}
	return nil
}
