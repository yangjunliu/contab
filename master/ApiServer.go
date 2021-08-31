package master

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/yangjunliu/crontab/common"
	"net/http"
	"strconv"
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
// POST job={"name":"job1", "command":"echo hello", "cronExpr":"******}
func handleJobSave(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		job     common.Job
		oldJob  common.Job
		err     error
		bytes []byte
	)

	// 2.取表单中的job字段
	if err = json.NewDecoder(r.Body).Decode(&job); err != nil {
		goto ERR
	}

	// 3.保存到etcd
	oldJob, err = G_jobMgr.SaveJob(&job)
	if err != nil {
		goto ERR
	}

	// 4.返回
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	// 5.返回
	if bytes, err = common.BuildResponse(-1, err.Error(), 1); err == nil {
		w.Write(bytes)
	}
}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/job/save", handleJobSave)
	return router
}

// 初始化服务
func InitApiServer() (err error) {
	r := RegisterRouter()
	err = http.ListenAndServe(":"+strconv.Itoa(G_config.ApiPort), r)

	return
}
