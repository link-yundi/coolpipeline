package coolpipeline

import (
	"sync"
)

/*
*
------------------------------------------------
Created on 2022-11-04 17:34
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
*
*/
type Worker func(in any) (out any)

type Pipeline struct {
	entry       Worker // 步骤1
	final       Worker // 最后的步骤
	workflows   []Worker
	workingChan chan int // 正在工作的数量
	inWg        *sync.WaitGroup
	inCache     chan any // 参数缓存
	exitFlag    bool
	threadWg    *sync.WaitGroup
}

func NewPipelines(parallelSize int, workers ...Worker) *Pipeline {
	newWorkers := make([]Worker, 0)
	for _, w := range workers {
		if w != nil {
			newWorkers = append(newWorkers, w)
		}
	}
	var entry, final Worker
	if len(newWorkers) > 0 {
		entry = newWorkers[0]
		final = newWorkers[len(newWorkers)-1]
	}
	pipeline := &Pipeline{
		entry:       entry,
		final:       final,
		workflows:   newWorkers,
		workingChan: make(chan int, parallelSize),
		inWg:        &sync.WaitGroup{},
		threadWg:    &sync.WaitGroup{},
		inCache:     make(chan any, parallelSize),
	}
	return pipeline
}

// 开始第一步的任务
func (pl *Pipeline) start(in any) {
	pl.threadWg.Add(1)
	for !pl.exitFlag {
		if len(pl.workflows) < 1 {
			return
		}
		var (
			out any
		)
		out = pl.entry(in)
		for _, w := range pl.workflows[1:] {
			in = out
			out = w(in)
		}
		pl.inWg.Done()
		in = <-pl.inCache
	}
	pl.threadWg.Done()
	<-pl.workingChan
}

func (pl *Pipeline) AddTask(ins ...any) {
	for _, in := range ins {
		pl.inWg.Add(1)
		select {
		case pl.workingChan <- 1:
			go pl.start(in)
		default:
			pl.inCache <- in
		}
	}
	pl.inWg.Wait()
	close(pl.inCache)
	pl.exitFlag = true
	pl.threadWg.Wait()
	close(pl.workingChan) //
}
