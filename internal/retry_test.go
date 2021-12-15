package internal

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/crochee/lirity/routine"
)

type Manager struct {
	model map[string]Executor
}

func (m *Manager) register(v Executor) error {
	var (
		vt   = reflect.TypeOf(v)
		name string
	)
	switch vt.Kind() {
	case reflect.Struct:
		name = vt.String()
	case reflect.Ptr:
		name = vt.Elem().String()
	default:
		return fmt.Errorf("not support %s", vt.String())
	}
	if _, ok := m.model[name]; ok {
		return fmt.Errorf("%s is enable", name)
	}
	m.model[name] = v
	return nil
}

func (m *Manager) Register(executors ...Executor) error {
	for _, executor := range executors {
		if err := m.register(executor); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Run(name string) error {
	v, ok := m.model[name]
	if !ok {
		return fmt.Errorf("must register model %s", name)
	}
	return v.Copy().Run(context.Background(), nil)
}

type Executor interface {
	Copy() Executor
	Run(ctx context.Context, data []byte) error
}

type test struct {
}

func (t test) Copy() Executor {
	return t
}

func (t test) Run(ctx context.Context, data []byte) error {
	fmt.Print("90\t")
	return nil
}

type test1 struct {
	i uint
}

func (t *test1) Copy() Executor {
	tmp := *t
	return &tmp
}

func (t *test1) Run(ctx context.Context, data []byte) error {
	t.i++
	fmt.Print("91\t")
	fmt.Printf("%v \t", t)
	return nil
}

type multiTest struct {
	list []Executor
}

func (m *multiTest) Copy() Executor {
	tmp := &multiTest{list: make([]Executor, 0, len(m.list))}
	for _, e := range m.list {
		tmp.list = append(tmp.list, e.Copy())
	}
	return tmp
}

func (m *multiTest) Run(ctx context.Context, data []byte) error {
	fmt.Println(len(m.list))
	g := routine.NewGroup(ctx)
	for _, e := range m.list {
		tmp := e
		g.Go(func(ctx context.Context) error {
			return tmp.Run(ctx, data)
		})
	}
	return g.Wait()
}

func TestRetry(t *testing.T) {
	m := &Manager{model: map[string]Executor{}}
	if err := m.Register(test{}, &test1{}, &multiTest{list: []Executor{test{}, &test1{}}}); err != nil {
		t.Fatal(err)
	}
	t.Log(m.Run("internal.test"))
	t.Log(m.Run("internal.test1"))
	t.Log(m.Run("internal.multiTest"))
	t.Log(m.Run("internal.multiTest"))
}
