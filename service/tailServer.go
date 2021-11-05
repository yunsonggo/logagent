package service

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic string `json:"topic"`
}

type TextMsg struct {
	Msg string
	Topic string
}

type TailObj struct {
	Tail *tail.Tail
	Conf CollectConf
	Status int
	ExitChan chan int
}

type TailObjMgr struct {
	TailObjs []*TailObj
	MsgChan chan *TextMsg
	Lock sync.Mutex
}

var TMgr *TailObjMgr

func GetOneLine() (msg *TextMsg) {
	msg = <-TMgr.MsgChan
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {

	TMgr = &TailObjMgr{
		MsgChan: make(chan *TextMsg, chanSize),
	}

	if len(conf) == 0 {
		logs.Error("invalid config for log collect, conf:%v", conf)
		return
	}

	for _, v := range conf {
		CreateNewTask(v)
	}

	return
}

func UpdateConfig(confs []CollectConf) (err error) {
	TMgr.Lock.Lock()
	defer TMgr.Lock.Unlock()

	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range TMgr.TailObjs {
			if oneConf.LogPath == obj.Conf.LogPath {
				isRunning = true
				break
			}
		}

		if isRunning {
			continue
		}

		CreateNewTask(oneConf)
	}

	var tailObjs []*TailObj
	for _, obj := range TMgr.TailObjs {
		obj.Status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.Conf.LogPath {
				obj.Status = StatusNormal
				break
			}
		}

		if obj.Status == StatusDelete {
			obj.ExitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}

	TMgr.TailObjs = tailObjs
	return
}

func CreateNewTask(conf CollectConf) {

	obj := &TailObj{
		Conf:     conf,
		ExitChan: make(chan int, 1),
	}

	tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})

	if errTail != nil {
		logs.Error("collect filename[%s] failed, err:%v", conf.LogPath, errTail)
		return
	}

	obj.Tail = tails
	TMgr.TailObjs = append(TMgr.TailObjs, obj)

	go ReadFromTail(obj)

}


func ReadFromTail(tailObj *TailObj) {
	for true {
		select {
		case line, ok := <-tailObj.Tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n", tailObj.Tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.Conf.Topic,
			}

			TMgr.MsgChan <- textMsg
		case <-tailObj.ExitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.Conf)
			return

		}
	}
}
