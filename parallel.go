package coolpipeline

/**
------------------------------------------------
Created on 2022-11-04 18:40
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
**/

type pipelineTask struct {
	Pipeline *Pipeline
	In       any
}

// ========================== 协程池 ==========================
type parallelPool struct {
	size  int
	tasks chan *pipelineTask // 任务
	num   chan int           // 正在工作数量
}

func newParallelPool(size int) *parallelPool {
	return &parallelPool{
		size:  size,
		tasks: make(chan *pipelineTask),
		num:   make(chan int, size),
	}
}

// 往协程池中添加任务
func (pp *parallelPool) addTask(task *pipelineTask) {
	select {
	case pp.tasks <- task:
	case pp.num <- 1:
		go pp.worker(task)
	}
}

// 执行任务
func (pp *parallelPool) worker(task *pipelineTask) {
	for {
		task.Pipeline.start(task.In)
		task = <-pp.tasks
	}
}
