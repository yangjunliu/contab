package master

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/yangjunliu/crontab/common"
)

// http的接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例对象
	G_apiServer *ApiServer
)

// 任务保存到etcd中
// POST job={"name":"job1", "command":"echo hello", "cronExpr":"* * * * * *}
func handleJobSave(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		job    common.Job
		oldJob common.Job
		err    error
		bytes  []byte
	)

	// 取表单中的job字段
	if err = json.NewDecoder(r.Body).Decode(&job); err != nil {
		common.ErrorResponse(w, common.ErrorApi{Code: 10001, Msg: err.Error()})
		return
	}

	// 保存到etcd
	oldJob, err = G_jobMgr.SaveJob(&job)
	if err != nil {
		common.ErrorResponse(w, common.ErrorApi{Code: 10002, Msg: err.Error()})
		return
	}

	// 保存成功
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
}

// 任务从etcd中删除
func handleJobDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		job     common.Job
		oldJobs []common.Job
		err     error
		bytes   []byte
	)

	// 获取表单
	if err = json.NewDecoder(r.Body).Decode(&job); err != nil {
		common.ErrorResponse(w, common.ErrorApi{Code: 10003, Msg: err.Error()})
	}

	// 删除job
	oldJobs, err = G_jobMgr.DelJob(job.Name)
	if err != nil {
		common.ErrorResponse(w, common.ErrorApi{Code: 10004, Msg: err.Error()})
	}

	if bytes, err = common.BuildResponse(0, "success", oldJobs); err == nil {
		w.Write(bytes)
	}
}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/job/save", handleJobSave)
	router.POST("/job/delete", handleJobDelete)
	return router
}

// 初始化服务
func InitApiServer() (err error) {
	r := RegisterRouter()
	err = http.ListenAndServe(":"+strconv.Itoa(G_config.ApiPort), r)

	return
}
