package helpers

import (
	"testing"
	//"container/ring"
)

func TestRing(t *testing.T) {

	r := RingNewInt64(10)

	for i := 0; i < 20; i++ {
		r.Value = int64(20 * i)
		r = r.Next()
	}
	s := RingToStringInt64(r, ",", true)

	t.Errorf(s)
	s = RingToStringInt64(r, ",", false)

	t.Errorf(s)
	t.Fail()

}
