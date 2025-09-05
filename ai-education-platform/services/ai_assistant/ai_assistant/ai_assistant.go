package main

import (
	"flag"
	"fmt"

	"ai-education-platform/services/ai_assistant/ai_assistant/internal/config"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/handler"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ai-assistant-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册路由
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting AI assistant service at %s:%d...\n", c.Host, c.Port)
	server.Start()
}