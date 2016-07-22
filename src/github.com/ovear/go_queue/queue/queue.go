package queue

import (
	"sync"

	"github.com/ovear/go_queue/config"
	"github.com/ovear/go_queue/logger"
	"github.com/ovear/go_queue/tasks"
	"github.com/ovear/go_queue/tools"
)

var threadNum int
var ch chan *tasks.Task
var mainWg *sync.WaitGroup

func init() {
	logger.Info("queue thread init..")
	threadNum = config.GetThreadNum()
	ch = tools.GetHandlerChan()
	mainWg = tools.GetMainWg()
	logger.InfoF("queue handler thread num [%d]", threadNum)
	//这里使用一个新的线程处理，防止因为 channel 长度不足导致的死锁
	go initUnhandledTasks()
}

func initUnhandledTasks() {
	tArr := tasks.GetUnhandledTasks()
	tArr = append(tArr, tasks.GetHandingTasks()...)
	logger.InfoF("load unhandled tasks from db [%d]", len(tArr))
	for _, t := range tArr {
		ch <- t
	}
}

func StartQueue() {
	logger.Info("queue thread start..")
	for i := 0; i < threadNum; i++ {
		mainWg.Add(1)
		go func(tnum int) {
			logger.InfoF("Handle thread [%d] start", tnum)
			for t := range ch {
				//设置任务为正在处理中
				t.Status = 2
				succ, _ := tasks.UpdateTaskToDb(t)
				if !succ {
					continue
				}
				logger.InfoF("[Thread-%d] handing [%s]", tnum, t.Url)
				content, statusCode, err := tools.Fetch(t.Url)
				if err != nil {
					//标记为处理失败
					t.Status = 4
					tasks.UpdateTaskToDb(t)
					tasks.NewTaskLogToDb(t.Uuid, t.Url, statusCode, err.Error())
					continue
				}
				logger.DebugF("[Thread-%d] fetch [%s] completed [%d]", tnum, t.Url, statusCode)
				//更新为处理完毕
				t.Status = 3
				t.Result = content
				succ, _ = tasks.UpdateTaskToDb(t)
				if !succ {
					continue
				}
				tasks.NewTaskLogToDb(t.Uuid, t.Url, statusCode, t.Result)
				logger.InfoF("[Thread-%d] handiing [%s] completed", tnum, t.Uuid)

			}
			logger.InfoF("[Thread-%d]Receive exit signal...", tnum)
			mainWg.Done()
		}(i)
	}
}
