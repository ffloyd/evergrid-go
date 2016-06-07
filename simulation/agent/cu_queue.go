package agent

import "github.com/ffloyd/evergrid-go/global/types"

// CUQueue is an internal ControlUnit's queue
type cuQueue struct {
	workersQueues map[string]*jobQueue
}

func newCUQueue(workers []*types.WorkerInfo) *cuQueue {
	result := &cuQueue{
		workersQueues: make(map[string]*jobQueue),
	}

	for _, info := range workers {
		result.workersQueues[string(info.UID)] = &jobQueue{
			queue: make([]types.JobInfo, 0, 10),
			info:  info,
		}
	}

	return result
}

func (cuq *cuQueue) forWorker(workerName string) *jobQueue {
	return cuq.workersQueues[workerName]
}

type jobQueue struct {
	queue []types.JobInfo
	info  *types.WorkerInfo
}

func (q *jobQueue) push(job types.JobInfo) {
	q.queue = append(q.queue, job)
	q.info.QueueLength++
}

func (q *jobQueue) pop() *types.JobInfo {
	if len(q.queue) == 0 {
		return nil
	}

	var result types.JobInfo
	result, q.queue = q.queue[0], q.queue[1:len(q.queue)]
	q.info.QueueLength--
	return &result
}
