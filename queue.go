package main

import (
	"fmt"

	"github.com/ovear/go_queue/logger"
	"github.com/ovear/go_queue/queue"
	_ "github.com/ovear/go_queue/tasks"
	"github.com/ovear/go_queue/tools"
	"github.com/ovear/go_queue/web"
)

func main() {
	fmt.Println("Hello, Queue")

	logger.Info("System init start")
	web.Start()

	//mysql.InitDb()

	//	tasks.Test()
	queue.StartQueue()
	tools.GetMainWg().Wait()
}
