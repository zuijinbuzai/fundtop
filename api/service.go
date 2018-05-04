package api

import (
	"sync"
	"zuijinbuzai/go1/fundtop/api/types"
)

type Service struct {
	downloadHeapMu 	sync.Mutex
	newDownloads   	chan struct{}

	downloadHeap	*downloadChunkHeap
}

func (s *Service) Close() error {
	return nil
}

func New() *Service {
	s := &Service{
		newDownloads:	make(chan struct{}, 1),
		downloadHeap:	new(downloadChunkHeap),
	}
	return s
}

func (s *Service) GetAllFundData() (map[string]*types.Fund, error) {
	return nil, nil
}