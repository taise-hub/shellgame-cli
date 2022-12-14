package main

import (
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"github.com/taise-hub/shellgame-cli/server/infrastructure"
	"github.com/taise-hub/shellgame-cli/server/interfaces"
	"github.com/taise-hub/shellgame-cli/server/usecase"
	"log"
	"net/http"
)

func main() {
	containerHandler, err := infrastructure.NewContainerHandler()
	if err != nil {
		log.Fatal(err)
		return
	}
	consoleRepo := interfaces.NewContainerRepository(containerHandler)
	gameUsecase := usecase.NewGameInteractor(consoleRepo)
	gameController := interfaces.NewGameController(gameUsecase)

	go model.GetMatchingRoom().Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/profiles", gameController.Profile)
	mux.HandleFunc("/players", gameController.Match)
	mux.HandleFunc("/waitmatch", gameController.WaitMatch)
	mux.HandleFunc("/shell", gameController.Start)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":80", mux)
}
