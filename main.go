package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-dress/config"
	"go-dress/route"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var conf string
	flag.StringVar(&conf, "c", "", "config file path")
	config.Init(conf)
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	var host string
	flag.StringVar(&host, "h", cfg.Server.Host, "")
	var port int
	flag.IntVar(&port, "p", cfg.Server.Port, "")
	var mode string
	flag.StringVar(&mode, "m", "release", "")
	flag.Parse()
	gin.SetMode(mode)
	router := gin.Default()
	route.Init(router)
	go func() {
		err = router.Run(host + ":" + strconv.Itoa(port))
		if err != nil {
			fmt.Println(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = ctx.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server exiting")
}
