package logs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestVariousLogger(t *testing.T) {
	l := New(LevelDebug)
	sl := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Log("logger")
	l.With("key", "value").
		With("map", map[string]any{
			"struct": struct {
				Foo     string
				Bar     int
				Slice   []int
				Pointer *int
			}{
				Foo:     "foo",
				Bar:     123,
				Slice:   []int{1, 2, 3},
				Pointer: &[]int{1, 2, 3}[0],
			},
		}).
		With("slice", []int{1, 2, 3}).
		With("strings", strings.Join([]string{"foo", "bar", "baz"}, " -> ")).
		Info("logger")

	t.Log("slog")
	sl.With("key", "value").
		With("map", map[string]any{
			"struct": struct {
				Foo     string
				Bar     int
				Slice   []int
				Pointer *int
			}{
				Foo:     "foo",
				Bar:     123,
				Slice:   []int{1, 2, 3},
				Pointer: &[]int{1, 2, 3}[0],
			},
		}).
		With("slice", []int{1, 2, 3}).
		With("strings", strings.Join([]string{"foo", "bar", "baz"}, " -> ")).
		Info("slog")
}

func TestWithLoop(t *testing.T) {
	wg := sync.WaitGroup{}
	l := New(LevelDebug)
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			funcName := fmt.Sprintf("fund-%d", i)
			ll := l.With("func", funcName)
			ll.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}

func TestWithsLoop(t *testing.T) {
	wg := sync.WaitGroup{}
	l := New(LevelDebug)
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			funcName := fmt.Sprintf("fund-%d", i)
			l = l.With(map[string]any{"func": funcName})
			l.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}
