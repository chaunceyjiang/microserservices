package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"microserservices/files"
	"microserservices/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l:=log.New(os.Stdout,"product-api",log.LstdFlags)
	l.Println("Server Start ...")
	basePath:="/tmp/"
	hh :=handlers.NewProducts(l)
	fh:=handlers.NewFiles(l,files.NewLocal(basePath))
	//sm:=http.NewServeMux()
	// mux 多路复用器
	sm:=mux.NewRouter()
	//sm.Handle("/product",hh)
	// 创建一个路由
	getRouter:=sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",hh.GetProducts)

	postRouter:=sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products",hh.AddProducts)
	// 给路由增加中间件
	postRouter.Use(hh.MiddlewareProductValidation)

	putRouter:=sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}",hh.UpdateProduct)
	// 给路由增加中间件
	putRouter.Use(hh.MiddlewareProductValidation)

	// swagger 文档
	opts := middleware.RedocOpts{SpecURL: "/swagger.json"}
	sh:=middleware.Redoc(opts,nil)
	getRouter.Handle("/docs",sh)
	getRouter.Handle("/swagger.json",http.FileServer(http.Dir("./")))

	// 文件服务应该再拆分为另外一个微服务
	//postFileRouter:=sm.Methods(http.MethodPost).Subrouter()
	//postFileRouter.HandleFunc("/files/{id:[0-9]+}/{filename:[a-z_A-Z]+\\.[a-z]{3}",fh.UploadREST)
	_ = fh.UploadREST
	getFileRouter:=sm.Methods(http.MethodGet).Subrouter()
	getFileRouter.Handle("/files",http.StripPrefix("/files/",http.FileServer(http.Dir("./"))))
	s :=http.Server{
		Addr: ":9000",
		// 传 mux ，可以处理多个handle
		Handler: sm,
		IdleTimeout: 1200000 *time.Second,
		ReadTimeout: 20000*time.Second,
		WriteTimeout: 200000*time.Second,
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
