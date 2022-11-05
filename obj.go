package coolpipeline

import "sync"

/**
------------------------------------------------
Created on 2022-11-05 17:47
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
**/

// ========================== 对象池 ==========================
type objPool struct {
	*sync.Pool
	limitChan chan int
}

func newObjPool(limitNum int, newFunc func() any) *objPool {
	return &objPool{
		Pool:      &sync.Pool{New: newFunc},
		limitChan: make(chan int, limitNum),
	}
}

func (o *objPool) Get() any {
	o.limitChan <- 1
	return o.Pool.Get()
}

func (o *objPool) Put(x any) {
	o.Pool.Put(x)
	<-o.limitChan
}
