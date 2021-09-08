package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// 定时任务
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

type JobName struct {
	Name string `json:"name"`
}

// 列表
type JobList struct {
	Jobs []Job
}

// HTTP接口
type Response struct {
	ErrNo int         `json:"err_no"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errNo int, msg string, data interface{}) (resp []byte, err error) {
	response := Response{
		ErrNo: errNo,
		Msg:   msg,
		Data:  data,
	}

	resp, err = json.Marshal(response)
	return
}

// api error
type ErrorApi struct {
	Code int
	Msg  string
}

func (ae *ErrorApi) Error() string {
	return "code:" + strconv.Itoa(ae.Code) + " message:" + ae.Msg
}

func ErrorResponse(w http.ResponseWriter, err ErrorApi) {
	errBytes, e := BuildResponse(err.Code, err.Msg, "")
	if e != nil {
		errBytes, _ = BuildResponse(10000, "coding error", "")
	}

	w.Write(errBytes)
}
