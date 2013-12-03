package helpers

import (
	"bytes"
	"container/ring"
	"fmt"
	"strings"
)

func RingNewInt64(size int) *ring.Ring {
	r := ring.New(size + 1)
	for i := 0; i < r.Len(); i++ {
		r.Value = int64(0)
		r = r.Next()
	}
	return r
}

func RingToStringInt64(r *ring.Ring, delim string, delta bool) string {
	bb := new(bytes.Buffer)
	prev := r.Value.(int64)
	r = r.Next()
	for i := 1; i < r.Len(); i++ {
		if delta {
			if r.Value.(int64) == int64(0) {
				bb.WriteString(fmt.Sprintf("0%s", delim))
			} else {
				d := r.Value.(int64) - prev
				bb.WriteString(fmt.Sprintf("%d%s", d, delim))
			}
			prev = r.Value.(int64)
		} else {
			bb.WriteString(fmt.Sprintf("%d%s", r.Value.(int64), delim))
		}
		r = r.Next()
	}
	return strings.TrimSuffix(bb.String(), delim)
}
