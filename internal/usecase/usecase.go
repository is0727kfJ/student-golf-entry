package usecase

import (
	"errors"

	"github.com/is0727kfJ/student-golf-entry/internal/models"
	"github.com/is0727kfJ/student-golf-entry/internal/repository"
)

type ITournamentUsecase interface {
	GetTournaments() ([]models.Tournament, error)
	ApplyEntry(req models.EntryRequest) error
}

type tournamentUsecase struct {
	repo repository.ITournamentRepository // 倉庫番（インターフェース）を保持する
}

func NewTournamentUsecase(repo repository.ITournamentRepository) ITournamentUsecase {
	return &tournamentUsecase{repo: repo}
}

func (u *tournamentUsecase) GetTournaments() ([]models.Tournament, error) {
	return u.repo.GetAll()
}

func (u *tournamentUsecase) ApplyEntry(req models.EntryRequest) error {
	// ビジネスルールの検証（入力値が空でないか等）
	if req.TournamentID == "" || req.UserID == "" {
		return errors.New("invalid_input")
	}
	// 倉庫番にトランザクション処理を依頼
	return u.repo.CreateEntryTx(req.TournamentID, req.UserID)
}
