package model

import (
	"log"
)

var (
	matchingRoom *MatchingRoom = &MatchingRoom{
		Players:    make(map[uint32]*MatchingPlayer),
		message:    make(chan *MatchingMessage),
		register:   make(chan *MatchingPlayer),
		unregister: make(chan *MatchingPlayer),
	}
)

// shellgame-cliサーバ上で一つだけ存在。
// 対戦待ち状態の管理を行う。
type MatchingRoom struct {
	Players    map[uint32]*MatchingPlayer // 誰がMatchigRoomにいるのか把握するために利用。
	message    chan *MatchingMessage
	register   chan *MatchingPlayer
	unregister chan *MatchingPlayer
}

func GetMatchingRoom() *MatchingRoom {
	return matchingRoom
}

func (mr *MatchingRoom) GetRegisterChan() chan<- *MatchingPlayer {
	return mr.register
}

func (mr *MatchingRoom) GetUnregisterChan() chan<- *MatchingPlayer {
	return mr.unregister
}

func (mr *MatchingRoom) GetMatchingPlayers() []*MatchingPlayer {
	var players []*MatchingPlayer
	for _, player := range mr.Players {
		players = append(players, player)
	}
	return players
}

func (mr *MatchingRoom) Run() {
	for {
		select {
		case player := <-mr.register:
			log.Printf("[+] %s entered the room.\n", player.GetName())
			mr.Players[player.GetID()] = player
		case player := <-mr.unregister:
			log.Printf("[+] %s exited the room.\n", player.GetName())
			if _, ok := mr.Players[player.GetID()]; ok {
				close(player.matchingChan)
				delete(mr.Players, player.GetID())
			}
		case msg := <-mr.message:
			switch msg.Data {
			case OFFER:
				log.Printf("[+] OFFER: %s to %s\n", mr.Players[msg.Source].GetName(), mr.Players[msg.Dest].GetName())
				mr.HandleOffer(msg)
				mr.Players[msg.Dest].matchingChan <- msg
			case CANCEL_OFFER:
				log.Printf("[+] CANCEL OFFER: %s to %s\n", mr.Players[msg.Source].GetName(), mr.Players[msg.Dest].GetName())
				mr.HandleCancelOffer(msg)
				mr.Players[msg.Dest].matchingChan <- msg
			case ACCEPT:
				log.Printf("[+] ACCEPT OFFER: %s to %s\n", mr.Players[msg.Source].GetName(), mr.Players[msg.Dest].GetName())
				mr.HandleAccept(msg)
				mr.Players[msg.Source].matchingChan <- msg
				mr.Players[msg.Dest].matchingChan <- msg
			case DENY:
				log.Printf("[+] DENY OFFER: %s to %s\n", mr.Players[msg.Source].GetName(), mr.Players[msg.Dest].GetName())
				mr.HandleDeny(msg)
				mr.Players[msg.Dest].matchingChan <- msg
			default:
				err := &MatchingMessage{Data: ERROR}
				mr.Players[msg.Source].matchingChan <- err
			}
		}
	}
}

// 申請のハンドリング
// 申請者と承諾者のステータスが共にWAITINGであること確認し、両者のステータスをNEGOTIATINGにする。
func (mr *MatchingRoom) HandleOffer(msg *MatchingMessage) {
	// 受信者がRoomにいることを確認
	_, ok := mr.Players[msg.Dest]
	if !ok {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	if mr.Players[msg.Source].Status == NEGOTIATING || mr.Players[msg.Dest].Status == NEGOTIATING {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	mr.Players[msg.Source].Status = NEGOTIATING
	mr.Players[msg.Dest].Status = NEGOTIATING
}

// 申請キャンセルのハンドリング
// 申請者と承諾者のステータスが共にNEGITIATINGであること確認し、両者のステータスをWAITINGにする。
func (mr *MatchingRoom) HandleCancelOffer(msg *MatchingMessage) {
	_, ok := mr.Players[msg.Dest]
	if !ok {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	// FIXME: 交渉中でないPlayerを宛先にしてMessageを送信することで交渉解除することが可能。
	if mr.Players[msg.Source].Status == WAITING || mr.Players[msg.Dest].Status == WAITING {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	mr.Players[msg.Source].Status = WAITING
	mr.Players[msg.Dest].Status = WAITING
}

// 申請に対する承諾のハンドリング
// 申請者と承諾者のステータスが共にNEGITIATINGであること確認する。
func (mr *MatchingRoom) HandleAccept(msg *MatchingMessage) {
	_, ok := mr.Players[msg.Dest]
	if !ok {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	if mr.Players[msg.Source].Status == WAITING || mr.Players[msg.Dest].Status == WAITING {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}
}

// 申請に対する不承諾のハンドリング
// 申請者と承諾者のステータスが共にNEGITIATINGであること確認する。
func (mr *MatchingRoom) HandleDeny(msg *MatchingMessage) {
	//  以下の条件を満たす場合、申請者に対してマッチングが成立しなかったことを通達し、申請者と承認者のMatchingStateをWAITINGにする。
	_, ok := mr.Players[msg.Dest]
	if !ok {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	// FIXME: 交渉中でないPlayerを宛先にしてMessageを送信することで交渉解除することが可能。
	if mr.Players[msg.Source].Status == WAITING || mr.Players[msg.Dest].Status == WAITING {
		err := &MatchingMessage{Data: ERROR}
		mr.Players[msg.Source].matchingChan <- err
		return
	}

	mr.Players[msg.Source].Status = WAITING
	mr.Players[msg.Dest].Status = WAITING
}