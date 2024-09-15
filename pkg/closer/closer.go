package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

// Add - добавляет `func() error` в globalCloser
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait ...
func Wait() {
	globalCloser.Wait()
}

// CloseAll ...
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer ...
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New - возвращает Closer, если из ОС придет сигнал, то Closer вызовет CloseAll
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return c
}

// Add - метод объекта Closer, добавляет в объект функцию
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait - ожидает, пока не поступит сигнал в канал done
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll - вызывает все функции закрытия из объекта Closer
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
