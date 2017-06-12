package timewheel

import (
	"math"
	"sync"
	"time"
)

var (
	deletedTimerSlice DeleteTimerSlice
)

// dtsSize define DeleteTimerSlice Size,必须是2的N次方
const dtsSize uint64 = 0x10

func init() {
	deletedTimerSlice = make(DeleteTimerSlice, dtsSize)
	for i := 0; i < len(deletedTimerSlice); i++ {
		deletedTimerSlice[i].data = make(map[uint64]int64, 2048)
	}
}

// DeleteTimerSlice DeleteTimerSlice
type DeleteTimerSlice []DeleteTimer

// DeleteTimer 已经被删除的时间片
type DeleteTimer struct {
	sync.RWMutex
	data map[uint64]int64
}

// LastTime LastTime
func (d DeleteTimerSlice) LastTime(key uint64) int64 {
	index := key & (dtsSize - 1)
	d[index].RLock()
	defer d[index].RUnlock()
	if item, ok := d[index].data[key]; ok {
		return item
	}
	return math.MinInt64
}

// Add Add
func (d DeleteTimerSlice) Add(key uint64) {
	index := key & (dtsSize - 1)
	d[index].Lock()
	defer d[index].Unlock()
	d[index].data[key] = time.Now().UnixNano()
}

// Delete Delete
func (d DeleteTimerSlice) Delete(key uint64) {
	index := key & (dtsSize - 1)
	d[index].Lock()
	defer d[index].Unlock()
	delete(d[index].data, key)
}
