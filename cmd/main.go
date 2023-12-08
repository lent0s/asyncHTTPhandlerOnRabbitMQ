package cmd

import (
	"fmt"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/config"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/incoming_M1"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/inner_M2"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/logger"
	"github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/rabbit_MQ"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Local() {

	local(*config.ReadFlag())
}

func local(cfg config.Config) {

	logger.StartLog(cfg.LogPathFolder, cfg.LogFileName, cfg.LogMaxFileSizeKB)

	exit := make(chan struct{})
	closeMQ := make(chan struct{})
	defer fmt.Print("\r\t\b")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go gracefulShutdown(exit, &wg)

	go func() {
		defer wg.Done()
		rabbit_MQ.ServerUP(cfg.RConnect, cfg.RTimeout, closeMQ, &wg)
	}()

	wg.Wait()
	wg.Add(3)
	go inner_M2.WorkerUP()

	go func() {
		defer wg.Done()
		timeout := cfg.RTimeout
		incoming_M1.ServerUP(cfg.ServerHost, cfg.ServerPort, timeout, exit)
		closeMQ <- struct{}{}
	}()

	wg.Wait()
}

func gracefulShutdown(exit chan struct{}, wg *sync.WaitGroup) {

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal,
		syscall.SIGTERM,
		syscall.SIGINT)

	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		defer w.Done()
		go visualizationWorkConsole(exit, &w)
		<-exitSignal
		log.Println("completion signal received")
		close(exit)
	}()
	w.Wait()
	wg.Done()

	visualizationExitConsole()
}

func visualizationExitConsole() {

	dur := time.Duration(100)

	fmt.Print("exit")
	for {
		for j := 0; j < 3; j++ {
			time.Sleep(dur * time.Millisecond)
			fmt.Print(".")
		}
		for j := 0; j < 3; j++ {
			time.Sleep(dur * time.Millisecond)
			fmt.Print("\b \b")
		}
	}
}

func visualizationWorkConsole(exit chan struct{}, wg *sync.WaitGroup) {

	defer wg.Done()

	dur := time.Duration(100)
	fmt.Print("working ")

	for {
		select {
		case <-exit:
			fmt.Print("\r         \r")
			return
		default:
			fmt.Print("/")
			time.Sleep(dur * time.Millisecond)
			fmt.Print("\bꟷ")
			time.Sleep(dur * time.Millisecond)
			fmt.Print("\b\\")
			time.Sleep(dur * time.Millisecond)
			fmt.Print("\b|")
			time.Sleep(dur * time.Millisecond)
			fmt.Print("\b")
		}
	}
}
