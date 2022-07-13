package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brpaz/echozap"
	"github.com/duongcongtoai/toytoytoy/cmd/migration"
	"github.com/duongcongtoai/toytoytoy/internal/infras/mysql"
	"github.com/duongcongtoai/toytoytoy/internal/services"
	"github.com/duongcongtoai/toytoytoy/internal/transport/http"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	conf struct {
		Mysql mysql.Config
		Port  int
	}
	configPath = flag.String("config", "/configs/config.yaml", "path to config file")
)

func init() {
	flag.Parse()
	viper.SetConfigFile(*configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	e := echo.New()
	migration.Up(conf.Mysql.DSN)

	db := mysql.ConnectDB(conf.Mysql)
	zapLogger, _ := zap.NewProduction()

	e.Use(echozap.ZapLogger(zapLogger))

	wagerSvc := services.NewWagerSvc(db, &mysql.WagerRepo{})
	purchaseSvc := services.NewPurchaseSvc(db, &mysql.WagerRepo{}, &mysql.PurchaseRepo{})

	http.BindAPI(e, wagerSvc, purchaseSvc)

	go func() {
		e.Start(fmt.Sprintf(":%d", conf.Port))
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	e.Shutdown(context.Background())
}
