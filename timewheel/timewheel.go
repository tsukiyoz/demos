package timewheel

import (
	"container/list"
	"log"
	"time"
)

type taskPos struct {
	pos int
	ele *list.Element
}

type task struct {
	delay  time.Duration
	key    string
	circle int
	fn     func()
	cancel func()
}

type TimeWheel struct {
	interval   time.Duration
	ticker     *time.Ticker
	curSlotPos int
	slotNum    int
	slots      []*list.List
	m          map[string]*taskPos
	addCh      chan *task
	cancelCh   chan string
	stopCh     chan struct{}
}

func New(interval time.Duration, slowNum int) *TimeWheel {
	tw := &TimeWheel{
		interval:   interval,
		ticker:     time.NewTicker(interval),
		curSlotPos: 0,
		slotNum:    slowNum,
		slots:      make([]*list.List, slowNum),
		m:          make(map[string]*taskPos),
		addCh:      make(chan *task),
		cancelCh:   make(chan string),
		stopCh:     make(chan struct{}),
	}

	for i := 0; i < slowNum; i++ {
		tw.slots[i] = list.New()
	}

	return tw
}

func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.execTask()
		case t := <-tw.addCh:
			tw.addTask(t)
		case key := <-tw.cancelCh:
			tw.cancelTask(key)
		case <-tw.stopCh:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) execTask() {
	l := tw.slots[tw.curSlotPos]
	tw.curSlotPos = (tw.curSlotPos + 1) % tw.slotNum
	tw.scanList(l)
}

func (tw *TimeWheel) scanList(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)

		if task.circle > 0 {
			task.circle--
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					// log
					log.Println("Recover from timewheel task executing, err:", err)
				}
			}()

			task.fn()
		}()

		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(tw.m, task.key)
		}

		e = next
	}
}

func (tw *TimeWheel) posAndCicle(delay time.Duration) (pos, circle int) {
	delaySecond := int(delay.Seconds())
	intervalSecond := int(tw.interval.Seconds())

	pos = (tw.curSlotPos + delaySecond/intervalSecond) % tw.slotNum
	circle = (delaySecond / intervalSecond) / tw.slotNum
	return
}

func (tw *TimeWheel) addTask(t *task) {
	var pos int
	pos, t.circle = tw.posAndCicle(t.delay)

	ele := tw.slots[pos].PushBack(t)
	if t.key != "" {
		if _, ok := tw.m[t.key]; ok {
			tw.cancelTask(t.key)
		}
		tw.m[t.key] = &taskPos{pos: pos, ele: ele}
	}
}

func (tw *TimeWheel) cancelTask(key string) {
	taskPos, ok := tw.m[key]
	if !ok {
		return
	}

	tw.slots[taskPos.pos].Remove(taskPos.ele)
	delete(tw.m, key)

	taskPos.ele.Value.(*task).cancel()
}

// ------------- external api -------------

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *TimeWheel) Stop() {
	tw.stopCh <- struct{}{}
}

func (tw *TimeWheel) Add(delay time.Duration, key string, callback func()) {
	if delay < 0 {
		return
	}
	t := task{
		delay: delay,
		key:   key,
		fn:    callback,
	}
	tw.addCh <- &t
}

func (tw *TimeWheel) Cancel(key string) {
	tw.cancelCh <- key
}
