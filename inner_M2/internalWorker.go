package inner_M2

import (
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/rabbit_MQ"
	"log"
	"runtime"
)

func WorkerUP() {

	workerUP()
}

func workerUP() {

	cores := runtime.NumCPU()
	query := make(chan rabbit_MQ.Work, cores)

	go rabbit_MQ.GetFromM1(query)

	for i := 0; i < cores; i++ {
		go worker(query)
	}
}

func worker(query chan rabbit_MQ.Work) {

	for w := range query {
		w := w
		modifyData(&w)
		if err := rabbit_MQ.PostM2(w.Body, w.ContType, w.ID); err != nil {
			log.Printf("can't post data to MQ_M2: %s", err)
		}
	}
}

func modifyData(w *rabbit_MQ.Work) {

	mod := []byte("modified data: ")
	w.Body = append(mod, w.Body...)
}
