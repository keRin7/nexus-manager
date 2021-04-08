package main

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/keRin7/nexus-manager/nexusmanager"
	"github.com/keRin7/nexus-manager/pkg/handler"
	"github.com/keRin7/nexus-manager/pkg/webserver"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//ctx, finish := context.WithCancel(context.Background())

	config := nexusmanager.NewConfig()
	err := env.Parse(config)
	if err != nil {
		logrus.Fatal(err)
	}

	nexus := nexusmanager.New(config)
	//nexus.List()
	fmt.Println()
	//tags := nexus.ListTagsByImage("coolrocket/vnikay-dbupdate")

	//for _, v := range tags {
	//		headers := nexus.GetImageSHA("coolrocket/vnikay", v)
	//	nexus.GetDataV1("coolrocket/vnikay-dbupdate", v)
	//	size := nexus.GetSize("coolrocket/vnikay-dbupdate", v)
	//	fmt.Println(v, size/1024/1024)
	//		fmt.Println(v, headers["Last-Modified"], headers["Date"])
	//}

	//fmt.Println()
	//nexus.ImageManifest("coolrocket/vnikay", "f9a57d6c4a04448ce2e0f6feb5a61945a02fbd29")
	//fmt.Println()
	handlers := handler.NewHandler(nexus)

	srv := new(webserver.Server)
	srv.Run("8080", handlers.InitRoutes())
	//dbcounter := counter.New(config)
	//dbcounter.Init()
	//go dbcounter.TablesLenght(ctx)
	//dbcounter.Start()
	//finish()
}
