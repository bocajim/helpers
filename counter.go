package helpers

import (
	"container/ring"
	"sync"
	"time"
)

type counterValue struct {
	current   int64
	minAgo    int64
	minRing   *ring.Ring
	hourAgo   int64
	hourRing  *ring.Ring
}

var counterRwl sync.RWMutex
var counter map[string]*counterValue

func CounterInitialize() {
	counter = make(map[string]*counterValue)
	go CounterAggregator()
}

func CounterAggregator() {
 	var ctr int
	for {
		select {
			case <-time.After(1 * time.Minute):
				for _,c:=range counter {
				
					if c.current>c.minAgo {
						c.minAgo = c.current
					} else {
						c.minAgo = int64(0)
					}
					c.minRing.Value=c.current
					c.minRing=c.minRing.Next()
					
					if ctr%60==0 {
						if c.current>c.hourAgo {
							c.hourAgo = c.current
						} else {
							c.hourAgo = int64(0)
						}
						c.hourRing.Value=c.current
						c.hourRing=c.hourRing.Next()
					}
				}
				
				ctr++
				break
		}
	}
					
}

func CounterIncrement(bucket string, val int64) {
	if len(bucket)==0 || val==0 { return }
	counterRwl.Lock()
	if b,f:=counter[bucket];f {
		b.current+=val
	} else {
		counter[bucket]=&counterValue{val,0,RingNewInt64(60),0,RingNewInt64(72)}
	}
	counterRwl.Unlock()
}

func CounterSet(bucket string, val int64) {
	if len(bucket)==0 { return }
	counterRwl.Lock()
	if b,f:=counter[bucket];f {
		b.current=val
	} else {
		counter[bucket]=&counterValue{val,0,RingNewInt64(60),0,RingNewInt64(72)}
	}
	counterRwl.Unlock()
}

func CounterQuery(bucket string) (int64,int64,int64) {
	if len(bucket)==0 { return 0,0,0 }
	counterRwl.RLock()
	v,f:=counter[bucket]
	counterRwl.RUnlock()
	if !f {
		return 0,0,0
	}
	return v.current,v.minAgo,v.hourAgo
}

func CounterQueryEx(bucket string, delta bool) (int64,int64,string,int64,string) {
	if len(bucket)==0 { return 0,0,"",0,"" }
	counterRwl.RLock()
	v,f:=counter[bucket]
	counterRwl.RUnlock()
	if !f {
		return 0,0,"",0,""
	}
	return v.current,v.minAgo,RingToStringInt64(v.minRing,",",delta),v.hourAgo,RingToStringInt64(v.hourRing,",",delta)
}

func CounterList() map[string]int64 {
	counterRwl.RLock()
	nm := make(map[string]int64)
	for k, v := range counter {
	    nm[k] = v.current
	}
	counterRwl.RUnlock()
	return nm
}
