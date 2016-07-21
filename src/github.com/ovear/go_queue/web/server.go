package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ovear/go_queue/config"
	"github.com/ovear/go_queue/logger"
	"github.com/ovear/go_queue/tasks"
	"github.com/ovear/go_queue/tools"
)

const SERVER_INFO = "Go Queue"

var extra = map[string]string{
	"Server": SERVER_INFO,
}
var port int
var key string

func init() {
	port = config.GetHttpPort()
	key = config.GetAuthKey()
}

func HeaderSetter(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for k, v := range extra {
			rw.Header().Set(k, v)
		}
		fn(rw, req)
	})
}

func Start() {
	go func() {
		addr := fmt.Sprintf(":%d", port)
		logger.Info("Web server init start...")
		logger.InfoF("Web server init, listening at [%s]", addr)

		http.HandleFunc("/newTask", HeaderSetter(newTask))
		http.HandleFunc("/shutdown", HeaderSetter(shutdownHandler))

		http.HandleFunc("/testAdd", HeaderSetter(testAdd))

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logger.Fatal("Web server init failed, addr ", addr)
		}
	}()
}

func testAdd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprint(w, getJsonMsg(false, e))
		}
	}()
	if tools.IsServerClose() {
		panic("server closing")
	}
	//解析form
	r.ParseForm()
	form := make(map[string]string)
	for k, v := range r.Form {
		form[k] = v[0]
	}

	if form["pwd"] != key {
		fmt.Fprintf(w, "error: unknown")
		return
	}

	uuid := tools.GetRandomString(32)
	url := "http://192.168.44.1"
	if uuid == "" || url == "" {
		fmt.Fprint(w, getJsonMsg(false, "uuid or url is empty"))
	}

	t := tasks.NewTask(uuid, url)
	_, err := tasks.InsertTaskToDb(t)
	checkWebErr(err)

	tools.GetHandlerChan() <- t
	fmt.Fprint(w, getJsonMsg(true, "notify success"))
}

func newTask(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprint(w, getJsonMsg(false, e))
		}
	}()
	if tools.IsServerClose() {
		panic("server closing")
	}
	//解析form
	r.ParseForm()
	form := make(map[string]string)
	for k, v := range r.Form {
		form[k] = v[0]
	}

	if form["pwd"] != key {
		fmt.Fprintf(w, "error: unknown")
		return
	}

	uuid := form["uuid"]
	url := form["url"]
	if uuid == "" || url == "" {
		fmt.Fprint(w, getJsonMsg(false, "uuid or url is empty"))
	}

	t := tasks.NewTask(uuid, url)
	_, err := tasks.InsertTaskToDb(t)
	checkWebErr(err)

	tools.GetHandlerChan() <- t
	fmt.Fprint(w, getJsonMsg(true, "notify success"))
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	//解析form
	r.ParseForm()
	if v, exist := r.Form["pwd"]; !exist || v[0] != key {
		fmt.Fprintf(w, "error: unknown")
		return
	}
	fmt.Fprintf(w, "shutting down")
	logger.InfoF("System shutting down by %s", r.RemoteAddr)
	go tools.SetServerClose()

}

func getJsonMsg(succ bool, msg interface{}) (text string) {
	m := make(map[string]interface{})
	m["success"] = succ
	m["msg"] = fmt.Sprint(msg)
	tmp, _ := json.Marshal(m)
	text = string(tmp)
	return
}

func checkWebErr(err error) {
	if err != nil {
		logger.ErrorF("handing web request failed [%s]", err)
		panic(err)
	}
}
