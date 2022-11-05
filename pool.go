package coolpipeline

import (
	"sync"
)

/**
------------------------------------------------
Created on 2022-11-05 12:04
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
**/

type Pool struct {
	*parallelPool
	pipelinePool *sync.Pool
	taskPool     *sync.Pool      // 任务池
	taskInChan   chan any        // 参数通道，任务会不断监听该通道，根据读到的参数传递给任务并执行任务
	wg           *sync.WaitGroup // 任务计数器
}

// 流水线池
func NewPool(size int, workers ...Worker) *Pool {
	p := &Pool{
		taskInChan: make(chan any, size),
		wg:         &sync.WaitGroup{},
	}
	// 并发池
	parallelPool := newParallelPool(size)
	// 对象池
	pipelinePool := &sync.Pool{
		New: func() any {
			return newPipeline(workers...)
		},
	}
	taskPool := &sync.Pool{
		New: func() any {
			pt := &pipelineTask{
				Pipeline: pipelinePool.Get().(*Pipeline),
				In:       nil,
			}
			go p.success(pt)
			return pt
		},
	}
	p.parallelPool = parallelPool
	p.pipelinePool = pipelinePool
	p.taskPool = taskPool
	go p.listen()
	return p
}

// 启动监听任务数据
func (p *Pool) listen() {
	for in := range p.taskInChan {
		task := p.taskPool.Get().(*pipelineTask)
		// 监听任务是否成功
		task.In = in
		// 添加任务
		p.parallelPool.addTask(task)
	}
}

// 监听任务是否成功
func (p *Pool) success(task *pipelineTask) {
	for range task.Pipeline.final.outChan {
		p.wg.Done()
		p.taskPool.Put(task)
	}
}

// 等待任务完成
func (p *Pool) Wait() {
	p.wg.Wait()
}

// 关闭
func (p *Pool) Close() {
	close(p.taskInChan)
}

// 添加任务
func (p *Pool) AddTask(ins ...any) {
	for _, in := range ins {
		p.wg.Add(1)
		p.taskInChan <- in
	}
}
