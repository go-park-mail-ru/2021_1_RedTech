package domain

type Stream struct {
	Video  string `json:"video_path,omitempty"`
	Season int    `json:"season,omitempty"`
	Series int    `json:"series,omitempty"`
}

//go:generate mockgen -destination=../stream/repository/mock/mock_repo.go -package=mock Redioteka/internal/pkg/domain StreamRepository
type StreamRepository interface {
	GetStream(id uint) ([]Stream, error)
}

//go:generate mockgen -destination=../stream/usecase/mock/mock_usecase.go -package=mock Redioteka/internal/pkg/domain StreamUsecase
type StreamUsecase interface {
	GetStream(id uint) ([]Stream, error)
}
