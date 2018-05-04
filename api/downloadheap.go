package api

import (
	"sync"
	"container/heap"
	"zuijinbuzai/fundtop/api/types"
)

type downloadHeap struct {
	heap         	downloadChunkHeap
	newDownloads   	chan struct{}
	mu          	sync.Mutex
}

func (dh *downloadHeap) managedPush(dc *types.Fund)  {
	dh.mu.Lock()
	dh.heap.Push(dc)
	dh.mu.Unlock()
}

func (dh *downloadHeap) managedPop() (dc *types.Fund) {
	dh.mu.Lock()
	if len(dh.heap) > 0 {
		dc = heap.Pop(&dh.heap).(*types.Fund)
	}
	dh.mu.Unlock()
	return dc
}

type downloadChunkHeap []*types.Fund

func (dch downloadChunkHeap) Len() int {
	return len(dch)
}

func (dch downloadChunkHeap) Less(i, j int) bool {
	return dch[i].Code < dch[j].Code
}

func (dch downloadChunkHeap) Swap(i, j int) {
	dch[i], dch[j] = dch[j], dch[i]
}

func (dch *downloadChunkHeap) Push(x interface{}) {
	*dch = append(*dch, x.(*types.Fund))
}

func (dch *downloadChunkHeap) Pop() interface{} {
	old := *dch
	n := len(old)
	x := old[n-1]
	*dch = old[0 : n-1]
	return x
}