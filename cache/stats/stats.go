package stats

import "sync/atomic"

//功能拓展，解决信息量化

type Stats struct {
	Requests uint64
	Hits     uint64
	Misses   uint64
	DBLoads  uint64
}

func (s *Stats) RecordRequest() {
	atomic.AddUint64(&s.Requests, 1)
}

func (s *Stats) RecordHit() {
	atomic.AddUint64(&s.Hits, 1)
}

func (s *Stats) RecordMiss() {
	atomic.AddUint64(&s.Misses, 1)
}

func (s *Stats) RecordDBLoad() {
	atomic.AddUint64(&s.DBLoads, 1)
}
