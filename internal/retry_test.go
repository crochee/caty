package internal

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type Manager struct {
	model map[string]reflect.Type
}

func (m *Manager) register(v Executor) error {
	vt := reflect.TypeOf(v)
	var (
		name      string
		typeValue reflect.Type
	)
	switch vt.Kind() {
	case reflect.Struct:
		name = vt.String()
		typeValue = vt
	case reflect.Ptr:
		name = vt.Elem().String()
		typeValue = vt.Elem()
	default:
		return fmt.Errorf("not support %s", vt.String())
	}
	if _, ok := m.model[name]; ok {
		return fmt.Errorf("%s is enable", name)
	}
	m.model[name] = typeValue
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
	var instance Executor
	if instance, ok = reflect.New(v).Interface().(Executor); !ok {
		return fmt.Errorf("%s must implement Executor interface", name)
	}
	return instance.Run(context.Background(), nil)
}

type Executor interface {
	Run(ctx context.Context, data []byte) error
}

type test struct {
	i interface{}
}

func (t test) Run(ctx context.Context, data []byte) error {
	fmt.Print(90)
	return nil
}

type test1 struct {
	i interface{}
}

func (t *test1) Run(ctx context.Context, data []byte) error {
	t.i = 91
	fmt.Print(91)
	fmt.Printf("%#v", t)
	return nil
}

func TestRetry(t *testing.T) {
	m := &Manager{model: map[string]reflect.Type{}}
	if err := m.Register(test{}, &test1{}); err != nil {
		t.Fatal(err)
	}
	t.Log(m.Run("internal.test"))
	t.Log(m.Run("internal.test1"))
}
