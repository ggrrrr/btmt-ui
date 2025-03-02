package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

type (
	WaitFunc    func(ctx context.Context) error
	CleanupFunc func()

	Waiter interface {
		Add(fns ...WaitFunc)
		Wait() error
		Context() context.Context
		CancelFunc() context.CancelFunc
		AddCleanup(fns ...CleanupFunc)
	}

	waiter struct {
		ctx          context.Context
		fns          []WaitFunc
		cancel       context.CancelFunc
		cleanupFuncs []CleanupFunc
	}

	waiterCfg struct {
		parentCtx    context.Context
		catchSignals bool
	}
)

func New(options ...WaiterOption) Waiter {
	cfg := &waiterCfg{
		parentCtx:    context.Background(),
		catchSignals: false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		fns:          []WaitFunc{},
		cleanupFuncs: []CleanupFunc{},
	}
	w.ctx, w.cancel = context.WithCancel(cfg.parentCtx)
	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

func (w *waiter) Add(fns ...WaitFunc) {
	w.fns = append(w.fns, fns...)
}

func (w waiter) Wait() (err error) {
	g, ctx := errgroup.WithContext(w.ctx)

	// Here we wait for OS signal or for root ctx cancel call
	g.Go(func() error {
		<-ctx.Done()
		log.Log().Info("got kill signal")
		w.cancel()
		for _, fn := range w.cleanupFuncs {
			cleanupFunc := fn
			cleanupFunc()
		}
		return nil
	})

	// Here we are starting all Wait functions in goroutine
	for _, fn := range w.fns {
		waitFn := fn
		g.Go(func() error { return waitFn(ctx) })
	}

	log.Log().Info("started")
	return g.Wait()
}

func (w waiter) Context() context.Context {
	return w.ctx
}

func (w waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}

func (w *waiter) AddCleanup(fns ...CleanupFunc) {
	w.cleanupFuncs = append(w.cleanupFuncs, fns...)
}
