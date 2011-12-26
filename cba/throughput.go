package cba

import (
	"fmt"
	"github.com/hanwen/termite/stats"
	"time"
)

func (st *Store) initThroughputSampler() {
	st.throughput = stats.NewPeriodicSampler(time.Second, 60, func() stats.Sample {
		return &ThroughputSample{received: st.bytesReceived, served: st.bytesServed} 
	})
}

type ThroughputSample struct {
	served, received stats.MemCounter
}

func (s *ThroughputSample) CopySample() stats.Sample {
	t := *s
	return &t
}

func (s *ThroughputSample) String() string {
	return fmt.Sprintf("received %v, sent %v", s.received, s.served)
}

func (s *ThroughputSample) SubSample(r stats.Sample) {
	t := r.(*ThroughputSample)
	s.served -= t.served
	s.received -= t.received
}

func (s *ThroughputSample) AddSample(r stats.Sample) {
	t := r.(*ThroughputSample)
	s.served += t.served
	s.received += t.received
}

func (s *ThroughputSample) TableHeader() string {
	return "<tr><th>received</th><th>served</th></tr>"
}

func (s *ThroughputSample) TableRow() string {
	return fmt.Sprintf("<tr><td>%v</td><td>%v</td></tr>", s.received, s.served)
}

func (st *Store) ThroughputStats() []stats.Sample {
	return st.throughput.Diffs()
}

func (st *Store) addThroughput(received, served int64) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.bytesReceived += stats.MemCounter(received)
	st.bytesServed += stats.MemCounter(served)
}

