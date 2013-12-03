package helpers

import (
	"sync"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	CounterInitialize()

	CounterIncrement("abc123", 5)
	time.Sleep(time.Duration(1000))
	if c, _, _ := CounterQuery("abc123"); c != 5 {
		t.Errorf("CounterQuery() expected to get 5, but didn't.")
		t.Fail()
	}

	m := CounterList()
	v, f := m["abc123"]
	if !f {
		t.Errorf("CounterList() did not find expected item.")
		t.Fail()
	}
	if v != 5 {
		t.Errorf("CounterList() expected to get 5, but didn't.")
		t.Fail()
	}
}

func bcRoutine(b *testing.B, e chan bool) {
	for i := 0; i < b.N; i++ {
		CounterIncrement("abc123", 5)
		CounterIncrement("def456", 5)
		CounterIncrement("ghi789", 5)
		CounterIncrement("abc123", 5)
		CounterIncrement("def456", 5)
		CounterIncrement("ghi789", 5)
	}
	e <- true
}

func BenchmarkChannels(b *testing.B) {
	b.StopTimer()
	CounterInitialize()
	e := make(chan bool)
	b.StartTimer()

	go bcRoutine(b, e)
	go bcRoutine(b, e)
	go bcRoutine(b, e)
	go bcRoutine(b, e)
	go bcRoutine(b, e)

	<-e
	<-e
	<-e
	<-e
	<-e

}

var mux sync.Mutex
var m map[string]int

func bmIncrement(bucket string, value int) {
	mux.Lock()
	m[bucket] += value
	mux.Unlock()
}

func bmRoutine(b *testing.B, e chan bool) {
	for i := 0; i < b.N; i++ {
		bmIncrement("abc123", 5)
		bmIncrement("def456", 5)
		bmIncrement("ghi789", 5)
		bmIncrement("abc123", 5)
		bmIncrement("def456", 5)
		bmIncrement("ghi789", 5)
	}
	e <- true
}

func BenchmarkMutex(b *testing.B) {
	b.StopTimer()
	m = make(map[string]int)
	e := make(chan bool)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bmIncrement("abc123", 5)
		bmIncrement("def456", 5)
		bmIncrement("ghi789", 5)
		bmIncrement("abc123", 5)
		bmIncrement("def456", 5)
		bmIncrement("ghi789", 5)
	}

	go bmRoutine(b, e)
	go bmRoutine(b, e)
	go bmRoutine(b, e)
	go bmRoutine(b, e)
	go bmRoutine(b, e)

	<-e
	<-e
	<-e
	<-e
	<-e

}
