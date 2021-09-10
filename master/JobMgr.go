package master

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yangjunliu/crontab/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

type ErrorJob struct {
	Msg string
}

func (e *ErrorJob) Error() string {
	return e.Msg
}

func InitJobMgr() (err error) {
	var (
		client *clientv3.Client
	)
	// 初始化配置
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
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
		kv:     kv,
		lease:  lease,
	}
	return
}

func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob common.Job, err error) {
	// 把任务保存到/cron/jobs/任务名
	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
	)

	jobKey = common.JOB_NAME_PREFIX + job.Name
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

func (jobMgr *JobMgr) DelJob(jobName string) (oldJobs []common.Job, err error) {

	var (
		jobKey  string
		delResp *clientv3.DeleteResponse
		kvPair  *mvccpb.KeyValue
	)

	jobKey = common.JOB_NAME_PREFIX + jobName
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delResp.PrevKvs) > 0 {
		for _, kvPair = range delResp.PrevKvs {
			oldJob := new(common.Job)
			if err = json.Unmarshal(kvPair.Value, oldJob); err == nil {
				oldJobs = append(oldJobs, *oldJob)
			}
		}
	} else {
		err = &ErrorJob{Msg: "任务不存在"}
	}

	return
}
