package usecase

import (
	"Redioteka/internal/pkg/actor/repository/mock"
	"Redioteka/internal/pkg/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestActorUsecase_GetById(t *testing.T) {
	retActor := domain.Actor{
		ID:        1,
		FirstName: "Johny",
		LastName:  "Depp",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRepoMock := mock.NewMockActorRepository(ctrl)
	uc := NewActorUsecase(actorRepoMock)
	actorRepoMock.EXPECT().GetById(retActor.ID).Times(1).Return(retActor, nil)
	currentActor, currentError := uc.GetById(1)
	require.NoError(t, currentError)
	require.Equal(t, currentActor, retActor)
}

func TestActorUsecase_Search(t *testing.T) {
	retActors := []domain.Actor{
		{
			ID:        1,
			FirstName: "Johny",
			LastName:  "First",
		},
		{
			ID:        2,
			FirstName: "Johny",
			LastName:  "Second",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRepoMock := mock.NewMockActorRepository(ctrl)
	uc := NewActorUsecase(actorRepoMock)
	query := "Johny"
	actorRepoMock.EXPECT().Search(query).Times(1).Return(retActors, nil)
	currentActor, currentError := uc.Search(query)
	require.NoError(t, currentError)
	require.Equal(t, currentActor, retActors)
}
