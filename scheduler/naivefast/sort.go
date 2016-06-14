package naivefast

import "github.com/ffloyd/evergrid-go/global/types"

type byMFlopsDesc []types.WorkerInfo

func (sw byMFlopsDesc) Len() int {
	return len(sw)
}

func (sw byMFlopsDesc) Swap(i, j int) {
	sw[i], sw[j] = sw[j], sw[i]
}

func (sw byMFlopsDesc) Less(i, j int) bool {
	return sw[i].MFlops > sw[j].MFlops // сортируем по убыванию
}

type byQueueAsc []types.WorkerInfo

func (sw byQueueAsc) Len() int {
	return len(sw)
}

func (sw byQueueAsc) Swap(i, j int) {
	sw[i], sw[j] = sw[j], sw[i]
}

func (sw byQueueAsc) Less(i, j int) bool {
	return sw[i].QueueLength < sw[j].QueueLength // сортируем по убыванию
}
