package main

import (
	"context"
	"fmt"
	"github.com/go-text-parse/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	shutdownReq chan bool
	reqCount uint32
}

var port = os.Getenv("httpPort")

func StartServer() *Server {

	//if port == "" {
	//	port = ":80"
	//}

	fmt.Print(port)

	srv := &Server{
		Server:      http.Server{
			Addr:              ":9090",
			ReadTimeout:       10*time.Second,
			WriteTimeout:      10*time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.GetuploadForm).Methods("GET")
	router.HandleFunc("/upload", controllers.GetStatistic).Methods("POST")

	srv.Handler = router

	return srv
}

func main() {

	server := StartServer()

	done := make(chan bool)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		log.Println("Server starting")
		done <- true
	}()

	server.WaitShutdown()

	 <- done
	 log.Printf("DONE")

}

func (s *Server) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))

	if !atomic.CompareAndSwapUint32(&s.reqCount,0, 1){
		log.Printf("Shutdown through API call in progress..")
		return
	}

	go func() {
		s.shutdownReq <- true
	}()
}

func (s *Server) WaitShutdown()  {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <- irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stoping http server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}
