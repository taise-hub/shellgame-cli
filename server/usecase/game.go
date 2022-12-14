package usecase

import (
	"context"
	"github.com/taise-hub/shellgame-cli/common"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"github.com/taise-hub/shellgame-cli/server/domain/repository"
	"io"
	"log"
	"net"
	"time"
)

type GameInteractor struct {
	consoleRepo  repository.ConsoleRepository
}

func NewGameInteractor(consoleRepo repository.ConsoleRepository) *GameInteractor {
	return &GameInteractor{
		consoleRepo:  consoleRepo,
	}
}

// ゲーム開始時に利用する。
// クラアインとから受け取ったコネクションをコンソールの入出力先である別のコネクションに接続する。
func (gi *GameInteractor) Start(nconn net.Conn) (err error) {
	cconn, err := gi.consoleRepo.StartShell()
	if err != nil {
		log.Printf("Error in StartShell(): %v\n", err)
		return err
	}
	defer cconn.Close()

	go func() { io.Copy(cconn, nconn) }()
	io.Copy(nconn, cconn)
	return
}

func (gi *GameInteractor) ExtractMatchingProfiles(exceptID string) []*common.Profile {
	mroom := model.GetMatchingRoom()
	players := mroom.GetMatchingPlayers()
	var profiles []*common.Profile
	for _, v := range players {
		if v.Profile.ID == exceptID {
			continue
		}
		profiles = append(profiles, v.Profile)
	}
	return profiles
}

// playerをマッチング待ち状態にする。
func (gi *GameInteractor) WaitMatch(player *model.MatchingPlayer) {
	mroom := model.GetMatchingRoom()
	mroom.GetRegisterChan() <- player
	go player.ReadPump()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		player.WritePump(ctx)
		cancel()
	}()
}
