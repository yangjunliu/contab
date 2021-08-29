package master

import (
	"context"
	"encoding/json"
	"github.com/yangjunliu/crontab/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error)  {
	var (
		client *clientv3.Client
	)
	// 初始化配置
	config := clientv3.Config{
		Endpoints: G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv: kv,
		lease: lease,
	}
	return
}

func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob common.Job, err error) {
	// 把任务保存到/cron/jobs/任务名
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
	)

	jobKey = "/cron/jobs" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJob); err != nil {
			err = nil
		}
	}
	return
}