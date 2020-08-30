/*
	doonce 扩展
*/
package socket

import "sync"

type DoOnce struct {
	sync.Once
	done bool
}

func (d *DoOnce) Do(f func()) {
	d.Once.Do(f)
	d.done = true
}

func (d *DoOnce) IsDone() bool {
	return d.done
}
