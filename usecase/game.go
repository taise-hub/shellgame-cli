package usecase

import (
	"github.com/taise-hub/shellgame-cli/domain/repository"
)

type GameUsecase interface {
	Start() error
	Matching() error
}

type gameInteractor struct {
	shellRepo repository.ConsoleRepository
}

func NewGameInteractor(shellRepo repository.ConsoleRepository) GameUsecase {
	return &gameInteractor{
		shellRepo: shellRepo,
	}
}

// マッチング完了後、バトルを開始するためのメソッド
func (gc *gameInteractor) Start() error {
	panic("implement me")
}

// ゲームのマッチングをするためのメソッド
// 詳細まで考えれていないので、最終的にはもう少し細かく分かれる
func (gu *gameInteractor) Matching() error {
	panic("implement me")
}