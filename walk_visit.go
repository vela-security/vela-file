package file

import (
	"fmt"
	audit "github.com/vela-security/vela-audit"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
)

func (w *walk) dir(path string, f fs.FileInfo) error {
	atomic.AddInt32(&w.dirs, 1)
	w.output <- newInfo(path, f, nil)

	if !w.cfg.dir {
		return filepath.SkipDir
	}

	return nil
}

func (w *walk) newVisit(deep int) filepath.WalkFunc {
	return func(path string, fd fs.FileInfo, err error) error {
		if w.IsClose() {
			return fmt.Errorf("file.walk.%s over", w.Name())
		}

		if err != nil {
			return err
		}

		fi := newInfo(path, fd, err)
		if !w.Match(fi) {
			return nil
		}

		if fd.IsDir() {
			return w.dir(path, fd)
		}
		atomic.AddInt32(&w.files, 1)
		w.output <- fi
		return nil

	}
}

func (w *walk) visit(fi info) {
	if w.IsClose() {
		return
	}

	if !w.Match(fi) {
		return
	}

	w.output <- fi
}

func (w *walk) add(root string, dirs []fs.FileInfo, depth int) {
	n := len(dirs)
	if n == 0 {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			xEnv.Errorf("%s panic %v", w.Name(), e)
		}
	}()

	for i := 0; i < n; i++ {
		fi := dirs[i]
		path := filepath.Join(root, fi.Name())

		if !fi.IsDir() {
			w.visit(newInfo(path, fi, nil))
			continue
		}

		if w.cfg.deep == 0 || w.cfg.deep >= depth {
			w.readDirNames(wkItem{path: path, depth: depth})
			continue
		}
	}
}

func (w *walk) readDirNames(v wkItem) error {
	if w.IsClose() {
		return nil
	}

	stat, err := os.Stat(v.path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		w.visit(newInfo(v.path, stat, nil))
		return nil
	}

	if w.cfg.dir {
		w.visit(newInfo(v.path, stat, nil))
	}

	dirs, err := ioutil.ReadDir(v.path)
	if err != nil {
		return err
	}

	w.add(v.path, dirs, v.depth+1)
	return nil
}

func (w *walk) Ignore(f info) bool {
	return match(f.path, w.cfg.ignore)
}

func (w *walk) Match(f info) bool {
	if w.Ignore(f) {
		return false
	}

	n := len(w.cfg.filter)
	if n == 0 {
		return true
	}

	return match(f.path, w.cfg.filter)
}

func (w *walk) scan() {
	n := len(w.cfg.path)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		select {
		case <-w.ctx.Done():
			return

		default:
			path := w.cfg.path[i]
			w.readDirNames(wkItem{path: path, depth: 1})

			//case item := <-w.input:
			//	w.readDirNames(item)
			//}
		}
	}
}

func (w *walk) handle() {
	if w.cfg.pipe.Len() == 0 {
		return
	}

	for {
		select {
		case <-w.ctx.Done():
			return

		case v := <-w.output:
			w.cfg.pipe.Do(v, w.cfg.co, func(err error) {
				audit.Errorf("file.%s pipe call fail %v", w.Name(), err)
			})
		}
	}
}
