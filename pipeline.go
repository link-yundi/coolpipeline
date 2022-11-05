package coolpipeline

/**
------------------------------------------------
Created on 2022-11-04 17:34
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
**/

type Worker func(in any) (out any)

type workFlow struct {
	outChan chan any
	worker  Worker
}

func newWorkFlow(worker Worker) *workFlow {
	return &workFlow{
		outChan: make(chan any),
		worker:  worker,
	}
}

type Pipeline struct {
	entry     *workFlow // 步骤1
	final     *workFlow // 最后的步骤
	workflows []*workFlow
	pool      *parallelPool
	//wg        *sync.WaitGroup
}

func newPipeline(workers ...Worker) *Pipeline {
	var entry, final *workFlow
	workflowList := make([]*workFlow, 0)
	for _, w := range workers {
		wf := newWorkFlow(w)
		workflowList = append(workflowList, wf)
	}
	if len(workflowList) > 0 {
		entry = workflowList[0]
		final = workflowList[len(workflowList)-1]
	}
	pipeline := &Pipeline{
		entry:     entry,
		final:     final,
		workflows: workflowList,
	}
	pipeline.listen()
	return pipeline
}

// 启动监听
func (pl *Pipeline) listen() {
	if len(pl.workflows) < 1 {
		return
	}
	delivery := func(curCh chan any, nextWf *workFlow) {
		for d := range curCh {
			nextWf.outChan <- nextWf.worker(d)
		}
	}
	for index, wf := range pl.workflows[:len(pl.workflows)-1] {
		// 监听
		nextWf := pl.workflows[index+1]
		go delivery(wf.outChan, nextWf)
	}
}

// 开始第一步的任务
func (pl *Pipeline) start(d any) {
	if len(pl.workflows) < 1 {
		return
	}
	pl.entry.outChan <- pl.entry.worker(d)
}
