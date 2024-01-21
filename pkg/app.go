package httpworkers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Tudor036/gohttpworkers/pkg/server"
	"github.com/Tudor036/gohttpworkers/pkg/storage"
)

func Start() {
	storage := storage.NewStorage(storage.DefaultOptions(storage.WithAddr("127.0.0.1:6379"), storage.WithDB(0)))
	sv := server.NewServer(server.DefaultServerOptions(server.WithStorage(storage)))
	run(sv)
}

func run(sv *server.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)

	port := 9110
	go func() {
		addr := fmt.Sprintf(":%d", port)
		if err := http.ListenAndServe(addr, sv.AsHTTPHandler()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started at http://localhost:%d", port)
	<-done
	log.Println("Server stopped...")
}
