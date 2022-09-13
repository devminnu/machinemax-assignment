package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/devminnu/assignment/internal/app/register"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	wait := make(chan bool)
	go func() {
		defer func() {
			cancel()
			wait <- true
		}()
		log.Print("waiting for signal")
		<-ctx.Done()
		log.Print("signal received")
	}()
	r := register.New(&http.Client{})
	registeredIds, err := r.Register(ctx)
	if err != nil {
		return
	}
	bs, _ := json.Marshal(registeredIds)
	fmt.Println("registered euids:", string(bs))
	<-wait
}
