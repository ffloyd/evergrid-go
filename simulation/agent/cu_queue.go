package agent

import "github.com/ffloyd/evergrid-go/global/types"

// CUQueue is an internal ControlUnit's queue
type cuQueue struct {
	workersQueues map[string]*jobQueue
}

func newCUQueue(workerNames []string) *cuQueue {
	result := &cuQueue{
		workersQueues: make(map[string]*jobQueue),
	}

	for _, name := range workerNames {
		result.workersQueues[name] = &jobQueue{
			queue: make([]types.JobInfo, 0, 10),
		}
	}

	return result
}

func (cuq *cuQueue) forWorker(workerName string) *jobQueue {
	return cuq.workersQueues[workerName]
}

type jobQueue struct {
	queue []types.JobInfo
}

func (q *jobQueue) push(job types.JobInfo) {
	q.queue = append(q.queue, job)
}

func (q *jobQueue) pop() *types.JobInfo {
	if len(q.queue) == 0 {
		return nil
	}

	var result types.JobInfo
	result, q.queue = q.queue[0], q.queue[1:len(q.queue)]
	return &result
}
