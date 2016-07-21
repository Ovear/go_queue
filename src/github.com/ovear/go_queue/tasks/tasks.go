package tasks

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ovear/go_queue/logger"
	"github.com/ovear/go_queue/mysql"
)

var db *sql.DB

func init() {}

func getDb() *sql.DB {
	if db == nil {
		db = mysql.GetDb()
	}
	return db
}

func Test() {
	logger.Debug("tasks test start...")
	mdb := getDb()

	rdata := mdb.QueryRow("select * from tasks limit 1")
	mtask := new(Task)
	rdata.Scan(&mtask.Uuid, &mtask.Url, &mtask.Status,
		&mtask.Result, &mtask.Expected, &mtask.LastExecute, &mtask.CreatedAt)
	logger.Debug(mtask)
	logger.Debug("uuid -", mtask.Uuid)
	logger.Debug("url -", mtask.Url)
	//	mtask.Uuid = mtask.Uuid + "1"
	mtask.Status = 2
	//	UpdateTaskToDb(mtask)
	//InsertTaskToDb(mtask)
	NewTaskLogToDb(mtask.Uuid, mtask.Url, 200, "result hahah")

}

func NewTaskToDb(uuid, url string) (succ bool, err error) {
	return InsertTaskToDb(NewTask(uuid, url))
}

func NewTask(uuid, url string) (obj *Task) {
	obj = new(Task)
	obj.Uuid = uuid
	obj.Url = url
	//默认为1，未执行
	obj.Status = 1
	obj.Expected = ""
	obj.LastExecute = getCurTimeAsDbFormat()
	obj.CreatedAt = getCurTimeAsDbFormat()
	return
}

func InsertTaskToDb(t *Task) (succ bool, err error) {
	if t == nil {
		return
	}
	mdb := getDb()
	pstmt, err := mdb.Prepare("insert into tasks " +
		"(uuid, url, status, result, expected, lastExecute, createdAt) " +
		"values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.ErrorF("[task] insert task[%#v] failed [%#v]", t, err)
		pstmt.Close()
		return
	}
	result, err := pstmt.Exec(t.Uuid, t.Url, t.Status, t.Result,
		t.Expected, t.LastExecute, t.CreatedAt)
	pstmt.Close()
	if err != nil {
		logger.ErrorF("[task] insert task[%#v] failed [%#v]", t, err)
		return
	}
	if affect, _ := result.RowsAffected(); affect < 1 {
		logger.ErrorF("[task] insert task[%#v] failed [%#v] unexpected affect[%d]", t, err, affect)
		return
	}

	succ = true
	return
}

func UpdateTaskToDb(t *Task) (succ bool, err error) {
	if t == nil {
		return
	}
	mdb := getDb()
	pstmt, err := mdb.Prepare("update tasks " +
		"set url = ?, status = ?, result = ?, lastExecute = now() where uuid = ?")
	if err != nil {
		logger.ErrorF("[task] update task[%#v] failed [%#v]", t, err)
		pstmt.Close()
		return
	}
	result, err := pstmt.Exec(t.Url, t.Status, t.Result, t.Uuid)
	pstmt.Close()
	if err != nil {
		logger.ErrorF("[task] update task[%#v] failed [%#v]", t, err)
		return
	}
	if affect, _ := result.RowsAffected(); affect < 1 {
		logger.ErrorF("[task] update task[%#v] failed [%#v] unexpected affect[%d]", t, err, affect)
		return
	}

	succ = true
	return
}

//获得特定状态的Tasks, status=0则为所有
func GetTasksFromDb(status int) (tLists []*Task) {
	tLists = make([]*Task, 0)
	db := mysql.GetDb()
	sql := "select * from tasks"
	if status != 0 {
		sql = fmt.Sprintf("select * from tasks where status = %d", status)
	}

	rows, err := db.Query(sql)
	if err != nil {
		logger.Error("GetTasksFromDb failed ,", err)
		return
	}
	for rows.Next() {
		t := new(Task)
		rows.Scan(&t.Uuid, &t.Url, &t.Status, &t.Result,
			&t.Expected, &t.LastExecute, &t.CreatedAt)
		tLists = append(tLists, t)
	}
	return
}

func GetUnhandledTasks() (tLists []*Task) {
	return GetTasksFromDb(1)
}

func GetHandingTasks() (tLists []*Task) {
	return GetTasksFromDb(2)
}

type Task struct {
	Uuid     string
	Url      string
	Status   int
	Result   string
	Expected string
	//上次执行时间 2016-07-21 01:37:45
	LastExecute string
	//创建时间 2016-07-21 01:37:45
	CreatedAt string
}

func NewTaskLogToDb(task_uuid, url string, code int, result string) (succ bool, err error) {
	return InsertTaskLogToDb(NewTaskLog(task_uuid, url, code, result))
}

func InsertTaskLogToDb(tl *TaskLog) (succ bool, err error) {
	if tl == nil {
		return
	}
	mdb := getDb()
	pstmt, err := mdb.Prepare("insert into tasks_logs " +
		"(task_uuid, url, code, result, createdAt) " +
		"values (?, ?, ?, ?, ?)")
	if err != nil {
		logger.ErrorF("[task] insert task_log[%#v] failed [%#v]", tl, err)
		pstmt.Close()
		return
	}
	result, err := pstmt.Exec(tl.Task_uuid, tl.Url, tl.Code, tl.Result, tl.CreatedAt)
	pstmt.Close()
	if err != nil {
		logger.ErrorF("[task] insert task[%#v] failed [%#v]", tl, err)
		return
	}
	if affect, _ := result.RowsAffected(); affect < 1 {
		logger.ErrorF("[task] insert task[%#v] failed [%#v] unexpected affect[%d]", tl, err, affect)
		return
	}
	tl_id, _ := result.LastInsertId()
	tl.Id = int(tl_id)

	succ = true
	return
}

func NewTaskLog(task_uuid, url string, code int, result string) (tl *TaskLog) {
	tl = new(TaskLog)
	tl.Task_uuid = task_uuid
	tl.Url = url
	tl.Code = code
	tl.Result = result
	tl.CreatedAt = getCurTimeAsDbFormat()
	return
}

type TaskLog struct {
	Id        int
	Task_uuid string
	Url       string
	//http code
	Code      int
	Result    string
	CreatedAt string
}

func getCurTimeAsDbFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
