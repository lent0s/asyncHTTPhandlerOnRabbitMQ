package incoming_M1

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type control struct {
	interrupt bool
	clientsON sync.WaitGroup
}

func ServerUP(host, port string, timeout int, exit chan struct{}) {

	serverUP(host, port, timeout, exit)
}

func serverUP(host, port string, tOut int, exit chan struct{}) {

	c := &control{}
	mux := http.NewServeMux()
	mux.HandleFunc("/cmd", c.worker)

	srv := &http.Server{
		Addr:              host + ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * time.Duration(tOut+5),
		ErrorLog:          nil,
		BaseContext:       nil,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("server running on: %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server on %s fatal: %s", srv.Addr, err)
			close(exit)
			return
		}
	}()

	<-exit
	c.interrupt = true
	end := make(chan struct{})

	go func() {
		time.Sleep(time.Second * time.Duration(tOut+1))
		end <- struct{}{}
	}()

	go func() {
		c.clientsON.Wait()
		end <- struct{}{}
	}()

	<-end
	serverDN(srv, &wg)
}

func serverDN(srv *http.Server, wg *sync.WaitGroup) {

	log.Printf("server stopping on: %s\n", srv.Addr)
	go func() {
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("error during server shutdown on %s: %v\n", srv.Addr, err)
		}
	}()
	wg.Wait()
	log.Printf("server stopped on: %s\n", srv.Addr)
}
