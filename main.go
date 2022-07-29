package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"winprockill/internal/config"
	"winprockill/internal/handler"
	"winprockill/internal/service"

	_ "embed"
)

//go:embed index.html
var tmpl []byte

//go:embed nssm.exe
var nssm []byte

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	{
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}

	args := os.Args

	cmdService := service.NewWinCommand(args[0], c.ProcessNamePattern, nssm)

	if len(args) == 2 {
		if args[1] != "install" {
			log.Fatalln("invalid command ", args[0])
		}

		err = cmdService.InstallAsService()
		if err != nil {
			log.Fatalln("install as service: ", err)
		}

		log.Println("installed as a service")
		return
	}

	mux := http.DefaultServeMux
	{
		h := handler.New(cmdService, tmpl)

		mux.HandleFunc("/", h.Ui)
		mux.HandleFunc("/list", h.Processes)
		mux.HandleFunc("/control", h.Control)
	}

	server := &http.Server{
		Addr:    ":" + c.HTTPPort,
		Handler: mux,
	}

	log.Println("start server at: ", c.HTTPPort)
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("server start: ", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	s := <-signals
	log.Println("program shutdown, get signal from os: ", s)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Println("server shutdown: ", err)
	}
}
