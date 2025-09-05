package main

import (
	"flag"
	"fmt"

	"ai-education-platform/services/course/course/internal/config"
	"ai-education-platform/services/course/course/internal/handler"
	"ai-education-platform/services/course/course/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/course-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册路由
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting course service at %s:%d...\n", c.Host, c.Port)
	server.Start()
}