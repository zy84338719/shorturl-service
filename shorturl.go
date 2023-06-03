package main

import (
	"flag"
	"fmt"
	"shorturl/internal/logic/shorturl"

	"shorturl/internal/config"
	"shorturl/internal/handler"
	"shorturl/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shorturl-api-dev.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	var cc config.Conf
	conf.MustLoad(*configFile, &cc)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c, cc)

	go shorturl.TimerUpdateCount(ctx.DB, ctx.Cache)
	go shorturl.TimerCheckStatus(ctx.DB, ctx.Cache)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
