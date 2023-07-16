package fs

import (
	"fmt"

	"github.com/rjeczalik/notify"
	rt "github.com/arnodel/golua/runtime"
)

type pathWatcher struct{
	path string
	callback *rt.Closure
	paused bool
	started bool
	ud *rt.UserData
	notifyChan chan notify.EventInfo
}

func (w *pathWatcher) start() {
	if w.callback == nil || w.started {
		return
	}

	w.started = true
	w.notifyChan = make(chan notify.EventInfo)
	notify.Watch(w.path, w.notifyChan, notify.All)

	go func() {
		for notif := range w.notifyChan {
			ev := notif.Event().String()
			path := notif.Path()

			_, err := rt.Call1(rtmm.MainThread(), rt.FunctionValue(w.callback), rt.StringValue(ev), rt.StringValue(path))
			if err != nil {
				// TODO: throw error
			}
		}
	}()
}

func (w *pathWatcher) stop() {
	w.started = false
	notify.Stop(w.notifyChan)
}

func watcherStart(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	pw, err := watcherArg(c, 0)
	if err != nil {
		return nil, err
	}

	pw.start()

	return c.Next(), nil
}

func watcherStop(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	pw, err := watcherArg(c, 0)
	if err != nil {
		return nil, err
	}

	pw.stop()

	return c.Next(), nil
}

func newWatcher(path string, callback *rt.Closure) *pathWatcher {
	pw := &pathWatcher{
		path: path,
		callback: callback,
	}
	pw.ud = watcherUserData(pw)
	pw.start()

	return pw
}

func watcherArg(c *rt.GoCont, arg int) (*pathWatcher, error) {
	j, ok := valueToWatcher(c.Arg(arg))
	if !ok {
		return nil, fmt.Errorf("#%d must be a watcher", arg + 1)
	}

	return j, nil
}

func valueToWatcher(val rt.Value) (*pathWatcher, bool) {
	u, ok := val.TryUserData()
	if !ok {
		return nil, false
	}

	j, ok := u.Value().(*pathWatcher)
	return j, ok
}

func watcherUserData(j *pathWatcher) *rt.UserData {
	watcherMeta := rtmm.Registry(watcherMetaKey)
	return rt.NewUserData(j, watcherMeta.AsTable())
}
