package xss

import (
	"Redioteka/internal/pkg/domain"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeUser(user *domain.User) {
	sanitizer := bluemonday.UGCPolicy()
	user.Username = sanitizer.Sanitize(user.Username)
	user.Email = sanitizer.Sanitize(user.Email)
}
