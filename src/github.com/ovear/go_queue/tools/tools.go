package tools

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/ovear/go_queue/logger"
	"github.com/ovear/go_queue/tasks"
)

var handlerChan chan *tasks.Task
var mainWg *sync.WaitGroup
var isClose bool

func init() {
	//TODO 将 channel 长度写入配置文件
	handlerChan = make(chan *tasks.Task, 1000)
	mainWg = new(sync.WaitGroup)
}

func GetMainWg() *sync.WaitGroup {
	return mainWg
}

func GetHandlerChan() chan *tasks.Task {
	return handlerChan
}

func Fetch(rurl string) (string, int, error) {
	return FetchWithRef(rurl, "")
}

func FetchWithRef(rurl string, referer string) (string, int, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", rurl, nil)
	if err != nil {
		logger.ErrorF("error occur while fetch [%s] [%s]", rurl, err)
		return "", 0, err
	}

	request.Header.Set("User-Agent", "Go_Queue")
	request.Header.Set("Content-Type", "application/x-w-form-urlencoded")
	if referer != "" {
		request.Header.Set("Referer", referer)
	}

	response, err := client.Do(request)
	if err != nil {
		logger.ErrorF("error occur while fetch [%s] [%s]", rurl, err)
		return "", 0, err
	}
	//如果处理失败则回调自身再试一次
	//	if err != nil {
	//		return fetch(rurl)
	//	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	//	realUrl := response.Request.URL

	if err != nil {
		logger.ErrorF("error occur while fetch [%s] [%s]", rurl, err)
		return "", 0, err
	}

	return string(content), response.StatusCode, nil

}

//生成随机字符串
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SetServerClose() {
	isClose = true
	//等待1000毫秒，防止channel写入错误
	time.Sleep(time.Millisecond * 1000)
	close(handlerChan)
}

func IsServerClose() bool {
	return isClose
}
