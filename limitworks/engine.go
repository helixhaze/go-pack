package limitworks

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

func NewConcurrencyEngine(maxConcurrencyWorks int) WorkEngine {

	return &LimitConcurrencyWorkEngine{
		currentActiveWorks: new(int64),
		Concurrency:        maxConcurrencyWorks,
		waitingWorks:       make(WorksQueue, 0),
	}
}

type WorkID int64
type WorksQueue []Work
type Work interface {
	ID() WorkID
	F() error
}

type WorkEngine interface {
	DoJob(simpleJob func() error) Work
	DoWork(work Work)
}

var worksQueueMux = sync.Mutex{}

func (wq *WorksQueue) Push(w Work) {

	worksQueueMux.Lock()
	defer worksQueueMux.Unlock()

	*wq = append(*wq, w)
}

func (wq *WorksQueue) Shift() (w Work) {

	worksQueueMux.Lock()
	defer worksQueueMux.Unlock()

	n := len(*wq)
	if n == 0 {
		return nil
	}

	old := *wq
	firstWork := old[0]

	if n > 1 {
		*wq = old[1:]
	} else {
		*wq = old[0:0]
	}

	//fmt.Printf("dequeued work from wating: %v\n", firstWork.ID())
	return firstWork
}

type ConcurrencyWork struct {
	_ID WorkID
	_F  func() error
}

func (fw *ConcurrencyWork) ID() WorkID {
	return fw._ID
}

func (fw *ConcurrencyWork) F() error {
	return fw._F()
}

type LimitConcurrencyWorkEngine struct {
	currentActiveWorks *int64
	Concurrency        int
	waitingWorks       WorksQueue
}

func (fwp *LimitConcurrencyWorkEngine) wakeUpWaitingWork() {

	if w := fwp.waitingWorks.Shift(); w != nil {
		fwp.DoWork(w)
	}
}

func (fwp *LimitConcurrencyWorkEngine) DoWork(w Work) {

	if atomic.LoadInt64(fwp.currentActiveWorks) >= int64(fwp.Concurrency) {
		fwp.waitingWorks.Push(w)
		return
	}
	atomic.AddInt64(fwp.currentActiveWorks, 1)
	//fmt.Printf("[work %v start] current active works: %d\n", w.ID(), atomic.LoadInt64(fwp.currentActiveWorks))

	go func() {
		_ = w.F()
		atomic.AddInt64(fwp.currentActiveWorks, -1)
		//fmt.Printf("[work %v end] current active works: %d\n", w.ID(), atomic.LoadInt64(fwp.currentActiveWorks))
		go fwp.wakeUpWaitingWork()
	}()
}

func (fwp *LimitConcurrencyWorkEngine) DoJob(simpleJob func() error) Work {

	work := &ConcurrencyWork{
		_ID: WorkID(rand.Int63()),
		_F:  simpleJob,
	}

	fwp.DoWork(work)
	return work
}
