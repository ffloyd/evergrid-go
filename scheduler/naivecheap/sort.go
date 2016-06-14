package naivecheap

import "github.com/ffloyd/evergrid-go/global/types"

type byPriceAsc []types.WorkerInfo

func (sw byPriceAsc) Len() int {
	return len(sw)
}

func (sw byPriceAsc) Swap(i, j int) {
	sw[i], sw[j] = sw[j], sw[i]
}

func (sw byPriceAsc) Less(i, j int) bool {
	return sw[i].PricePerTick < sw[j].PricePerTick // сортируем по убыванию
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
