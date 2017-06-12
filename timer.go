package timewheel

import "time"

var (
	publicTimer *Timer
)

// Timer Timer
type Timer struct {
	Second TimeWheel
	Minute TimeWheel
	Hour   TimeWheel
}

// Add will 添加时间点到时间轮中
func (t *Timer) Add(seed uint, ts ...*TimerSlice) {
	for _, item := range ts {
		if last := deletedTimerSlice.LastTime(item.id); last >= item.insertTime {
			// 此id插入时间比删除时间小
			continue
		}

		second := item.Second
		if seed > 0 {
			if second < 60 {
				second %= seed
			} else if second < 3600 {
				second = (second + item.SecondOffset) % seed
			} else if second < 86400 {
				second = (second + item.SecondOffset + item.MinuteOffset*60) % seed
			}
		} else {
			item.SecondOffset = publicTimer.Second.Index
			item.MinuteOffset = publicTimer.Minute.Index
		}

		if second < 60 {
			//插入到秒轮
			index := (t.Second.Index + second) % 60
			t.Second.AddToIndex(uint(index), item)
		} else if second < 3600 {
			//插入到分钟轮
			index := (t.Minute.Index + second/60) % 60
			t.Minute.AddToIndex(uint(index), item)
		} else if second < 86400 {
			//插入到小时轮
			index := (t.Hour.Index + second/3600) % 24
			t.Hour.AddToIndex(uint(index), item)
		}
	}
}

// PutTimer will 启动Timer
func PutTimer(second uint, repeat bool, id uint64, e interface{}, callBack CallBackType) {
	ts := &TimerSlice{
		CallBack:     callBack,
		Second:       second,
		SecondOffset: 0,
		MinuteOffset: 0,
		Repeat:       repeat,
		id:           id,
		insertTime:   time.Now().UnixNano(),
		e:            e,
	}
	if publicTimer == nil {
		publicTimer = newTimer()
	}
	publicTimer.Add(0, ts)

	return
}

// RemoveTimer will 删除Timer
func RemoveTimer(id uint64) bool {
	deletedTimerSlice.Add(id)
	return true
}

func newTimer() *Timer {
	t := &Timer{
		Second: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Minute: TimeWheel{TS: make([][]*TimerSlice, 60), Index: 0},
		Hour:   TimeWheel{TS: make([][]*TimerSlice, 24), Index: 0},
	}

	go func(refTimer *Timer) {
		// 时间轮询
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 走一秒
				refTimer.Second.Tick()
				if refTimer.Second.Index == 0 {
					// 走一分钟
					refTimer.Minute.Tick()
					tmp := refTimer.Minute.CurTimerSliceAndClear()
					refTimer.Add(60, tmp...)
					if refTimer.Minute.Index == 0 {
						// 走一小时
						refTimer.Hour.Tick()
						tmp := refTimer.Hour.CurTimerSliceAndClear()
						refTimer.Add(3600, tmp...)
					}
				}

				//处理当前秒
				tmp := refTimer.Second.CurTimerSliceAndClear()
				go func(t *Timer, items []*TimerSlice) {
					for i := 0; i < len(items); i++ {
						item := items[i]
						if last := deletedTimerSlice.LastTime(item.id); last >= item.insertTime {
							// 此id插入时间比删除时间小,已经被删除
							continue
						}
						go func(t *Timer, instance *TimerSlice) {
							instance.CallBack(instance.e)
							if instance.Repeat {
								t.Add(0, instance)
							}
						}(t, item)
					}
				}(refTimer, tmp)

			}
		}
	}(t)
	return t
}
