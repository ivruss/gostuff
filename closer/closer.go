package closer

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Closer struct {
	mutex sync.Mutex
	once  sync.Once
	done  chan struct{}

	logger *zap.Logger

	gracefulCloseTimeout  time.Duration
	gracefulShutdownFuncs []func() error

	forceCloseTimeout  time.Duration
	forceShutdownFuncs []func() error
}

func NewCloser(logger *zap.Logger, gracefulCloseTimeout, forceCloseTimeout time.Duration, signals ...os.Signal) *Closer {
	closer := &Closer{
		mutex:                 sync.Mutex{},
		once:                  sync.Once{},
		done:                  make(chan struct{}),
		logger:                logger,
		gracefulCloseTimeout:  gracefulCloseTimeout,
		gracefulShutdownFuncs: nil,
		forceCloseTimeout:     forceCloseTimeout,
		forceShutdownFuncs:    nil,
	}

	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)
			<-ch
			signal.Stop(ch)
			closer.CloseEverything()
		}()
	}
	return closer
}

func (c *Closer) Run(funcs ...func() error) {
	for _, fn := range funcs {
		go func(f func() error) {
			if err := f(); err != nil {
				c.logger.Sugar().Errorf("closer func error: %v", err)
				c.CloseEverything()
			}
		}(fn)
	}
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) Add(f ...func() error) {
	c.mutex.Lock()
	c.gracefulShutdownFuncs = append(c.gracefulShutdownFuncs, f...)
	c.mutex.Unlock()
}

func (c *Closer) AddForce(f ...func() error) {
	c.mutex.Lock()
	c.gracefulShutdownFuncs = append(c.forceShutdownFuncs, f...)
	c.mutex.Unlock()
}

func (c *Closer) CloseEverything() {
	c.once.Do(func() {
		defer close(c.done)

		c.logger.Info("started shutdowns")
		c.waitAllCloseFunctions()
	})
}

func (c *Closer) waitAllCloseFunctions() {
	c.mutex.Lock()
	gracefulFuncs := c.gracefulShutdownFuncs
	forceFuncs := c.forceShutdownFuncs
	c.mutex.Unlock()

	waitFunc := func(funcs []func() error, timeout time.Duration) bool {
		wg := &sync.WaitGroup{}
		for _, f := range funcs {
			wg.Add(1)
			go func(f func() error) {
				if err := f(); err != nil {
					c.logger.Sugar().Errorf("close failed: %v", zap.Error(err))
				}
				wg.Done()
			}(f)
		}
		ch := make(chan os.Signal, 1)
		go func() {
			wg.Wait()
			close(ch)
		}()

		timer := time.NewTimer(c.gracefulCloseTimeout)
		defer timer.Stop()

		select {
		case <-ch:
			return true
		case <-timer.C:
			return false
		}
	}
	if ok := waitFunc(gracefulFuncs, c.gracefulCloseTimeout); !ok {
		c.logger.Error("graceful shutdown error: timeout limit exceeded")
		if ok := waitFunc(forceFuncs, c.forceCloseTimeout); !ok {
			c.logger.Error("force shutdown error: timeout limit exceeded")
		}
	}
}
