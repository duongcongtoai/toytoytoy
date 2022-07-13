package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/brpaz/echozap"
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
	ServerHost = "db_test"
)

func getServerAddr() string {
	return fmt.Sprintf("%s:%d", ServerHost, conf.Port)
}

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

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	e := echo.New()
	db := mysql.ConnectDB(conf.Mysql)
	err := mysql.CleanUpTestData(db)
	if err != nil {
		panic(err)
	}
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
	return 0
}
