package socket

import "testing"

func TestDonoce(t *testing.T) {
	var num int
	var once DoOnce

	if once.IsDone() {
		t.Fatalf("before done IsDone is setted")
	}

	for i := 0; i < 10; i++ {
		once.Do(func() {
			num++
		})
	}

	if num != 1 {
		t.Fatalf("donoce failed %v", num)
	}

	if !once.IsDone() {
		t.Fatalf("after done IsDone is not set")
	}
}
