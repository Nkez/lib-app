package main

import (
	library_app "github.com/Nkez/lib-app.git"
	"github.com/Nkez/lib-app.git/pkg/handler"
	"log"
)

func main() {
	srv := new(library_app.Server)
	handlers := new(handler.Handler)
	if err := srv.Run("8080", handlers.InitRouter()); err != nil {
		log.Fatal("Error while running server")
	}
}
