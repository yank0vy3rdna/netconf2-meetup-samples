package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	http_proxy "github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/adapters/controllers/http_proxy"
	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/repository/sysrepo_repo"
	persist_running "github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/usecases/persist_running"
	"github.com/yank0vy3rdna/netconf2-meetup-samples/sample-app/internal/usecases/proxy"
)

func main() {
	fmt.Println("App starting...")
	repo, err := sysrepo_repo.NewRepo()
	if err != nil {
		panic(err)
	}

	proxyUseCase, err := proxy.NewUseCase(repo)
	if err != nil {
		panic(err)
	}
	persist_running.NewUseCase(repo)

	proxyController := http_proxy.NewProxy(proxyUseCase)

	fmt.Println("App started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("App stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	proxyController.Close(ctx)
	repo.Close()
	fmt.Println("App stopped")
}
