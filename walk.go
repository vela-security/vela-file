package file

import (
	"context"
	"github.com/vela-security/vela-public/lua"
	"reflect"
)

const (
	Run int = iota + 1
	Init
	Err
)

var (
	walkTypeof = reflect.TypeOf((*walk)(nil)).String()
)

type walk struct {
	lua.ProcEx
	cfg    *walkConfig
	output chan info
	ctx    context.Context
	stop   context.CancelFunc
	offset int
	dirs   int32
	files  int32
}

func newWalk(cfg *walkConfig) *walk {
	w := &walk{cfg: cfg}
	return w
}

func (w *walk) Name() string {
	return w.cfg.name
}

func (w *walk) Type() string {
	return walkTypeof
}

func (w *walk) Close() error {
	w.stop()
	w.V(lua.PTClose)
	close(w.output)
	return nil
}

func (w *walk) pretreatment() {
	ctx, stop := context.WithCancel(context.Background())
	w.ctx = ctx
	w.stop = stop
	w.offset = 0
	w.output = make(chan info, 64)
}
func (w *walk) Start() error {
	w.pretreatment()
	xEnv.Spawn(0, w.handle)
	xEnv.Spawn(0, w.scan)
	return nil
}
