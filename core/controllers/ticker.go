package controllers

import (
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"time"
)

var LoopChan = make(chan CountType, 1000)

type CountType struct {
	UserId int64
	NodeId int64
	T      int // 1,2,3
}

func SendToLoop(userId, nodeId int64, t int) {
	flog.Log.Debugf("Ticker SendToLoop: %v", userId)
	LoopChan <- CountType{
		UserId: userId,
		NodeId: nodeId,
		T:      t,
	}
}

func LoopCount() {
	flog.Log.Debugf("Ticker start")
	for {
		select {
		case v := <-LoopChan:
			flog.Log.Debugf("Ticker Count: %v", v)
			if v.T == 1 {
				if v.UserId != 0 {
					err := model.CountContentAll(v.UserId)
					if err != nil {
						flog.Log.Errorf("Ticker Count all content err: %s", err.Error())
					}
				}
			} else if v.T == 2 {
				if v.UserId != 0 && v.NodeId != 0 {
					err := model.CountContentOneNode(v.UserId, v.NodeId)
					if err != nil {
						flog.Log.Errorf("Ticker Count node content err: %s", err.Error())
					}
				}
			} else if v.T == 3 {
				if v.UserId != 0 {
					err := model.CountContentCool(v.UserId)
					if err != nil {
						flog.Log.Errorf("Ticker Count all content cool err: %s", err.Error())
					}
				}
			}

		case <-time.After(1 * time.Second):
			flog.Log.Debugf("Ticker...")

		}
	}
}
