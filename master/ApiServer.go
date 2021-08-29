package master

import (
	"encoding/json"
	"github.com/yangjunliu/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
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
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		job     common.Job
		oldJob  common.Job
		postJob string
		err     error
	)
	// 1.解析POST表单
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	// 2.取表单中的job字段
	postJob = r.PostForm.Get("job")
	// 3.反序列化etcd
	if err := json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 4.保存到etcd
	oldJob, err = G_jobMgr.SaveJob(&job)
	if err != nil {
		goto ERR
	}
	// 5.返回
	if bytes, err := common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	// 5.返回
	if bytes, err := common.BuildResponse(-1, err.Error(), ""); err == nil {
		w.Write(bytes)
	}
}

// 初始化服务
func InitApiServer() (err error) {

	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动服务端
	go httpServer.Serve(listener)

	return
}
