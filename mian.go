package main

import (
	"context"
	"log"
	"microserservices/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l:=log.New(os.Stdout,"product-api",log.LstdFlags)
	l.Println("Server Start ...")
	hh :=handlers.NewProducts(l)
	sm:=http.NewServeMux()
	// mux 多路复用器
	sm.Handle("/",hh)

	s :=http.Server{
		Addr: ":9000",
		// 传 mux ，可以处理多个handle
		Handler: sm,
		IdleTimeout: 120 *time.Second,
		ReadTimeout: 2*time.Second,
		WriteTimeout: 2*time.Second,
	}
	go func() {
		if err:=s.ListenAndServe();err!=nil{
			l.Fatalln(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt,os.Kill)
	// 没有收到信号，就会阻塞在此处
	sig:=<-sigChan
	l.Println("接收到信号，优雅的退出..",sig)
	ctx,_:=context.WithTimeout(context.Background(),20*time.Second)
	s.Shutdown(ctx)
}
