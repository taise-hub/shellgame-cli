package usecase

import (
	"github.com/taise-hub/shellgame-cli/domain/model"
	"github.com/taise-hub/shellgame-cli/domain/repository"
	"github.com/taise-hub/shellgame-cli/domain/service"
	"io"
	"log"
	"net"
)

type GameUsecase interface {
	Start(net.Conn) error
	WaitMatch(*model.Player) error
}

type gameInteractor struct {
	consoleRepo  repository.ConsoleRepository
	matchService *service.MatchService
}

func NewGameInteractor(consoleRepo repository.ConsoleRepository, matchService *service.MatchService) GameUsecase {
	return &gameInteractor{
		consoleRepo:  consoleRepo,
		matchService: matchService,
	}
}

// ゲーム開始時に利用する。
// クラアインとから受け取ったWebsocketをコンソールの入出力先であるsocketに接続する。
func (gi *gameInteractor) Start(nconn net.Conn) (err error) {
	cconn, err := gi.consoleRepo.StartShell()
	if err != nil {
		log.Printf("Error in StartShell(): %v\n", err)
		return err
	}
	defer cconn.Close()
	go func() { _, _ = io.Copy(nconn, cconn) }()
	io.Copy(cconn, nconn)
	return
}

// playerをマッチング待ち状態にする。
func (gi *gameInteractor) WaitMatch(player *model.Player) error {
	if err := gi.matchService.Wait(player); err != nil {
		return err
	}
	//TODO: マッチングルーム用のchanelに追加するの処理の実装
	panic("implement me")
}
