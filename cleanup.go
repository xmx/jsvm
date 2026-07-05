package jsvm

import (
	"io"
	"log/slog"
	"sort"
	"sync"
)

// CleanHandle controls the lifecycle of a registered resource.
// Close closes the resource and removes it from the VM's cleanup list.
// Unregister removes the resource without closing it.
type CleanHandle interface {
	io.Closer

	// Unregister removes the resource from the VM's cleanup list.
	// Returns the original io.Closer, or nil if it was already removed.
	Unregister() io.Closer
}

type cleaner interface {
	register(io.Closer) int64
	unregister(id int64) io.Closer
	execute()
}

type cleanerMap struct {
	log     *slog.Logger
	serial  int64
	mutex   sync.Mutex
	closers map[int64]io.Closer
}

func newCleanerMap(log *slog.Logger) cleaner {
	return &cleanerMap{
		log:     log,
		closers: make(map[int64]io.Closer),
	}
}

func (cm *cleanerMap) register(c io.Closer) int64 {
	cm.mutex.Lock()
	cm.serial++
	id := cm.serial
	cm.closers[id] = c
	cm.mutex.Unlock()

	cm.log.Debug("vm cleanup register", "id", id)

	return id
}

func (cm *cleanerMap) unregister(id int64) io.Closer {
	cm.mutex.Lock()
	cls := cm.closers[id]
	delete(cm.closers, id)
	cm.mutex.Unlock()

	if cls == nil {
		cm.log.Debug("vm cleanup unregister nil", "id", id)
	} else {
		cm.log.Debug("vm cleanup unregister", "id", id)
	}

	return cls
}

func (cm *cleanerMap) execute() {
	cm.mutex.Lock()
	closers := cm.closers
	cm.closers = make(map[int64]io.Closer)
	cm.mutex.Unlock()

	ids := make([]int64, 0, len(closers))
	for id := range closers {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})

	for _, id := range ids {
		c := closers[id]
		if c == nil {
			continue
		}

		if err := c.Close(); err == nil {
			cm.log.Debug("vm cleanup close success", "id", id)
		} else {
			cm.log.Warn("vm cleanup close error", "id", id, "err", err)
		}
	}
}

type cleanHandle struct {
	id  int64
	cln cleaner
}

func (ch *cleanHandle) Close() error {
	if c := ch.Unregister(); c != nil {
		return c.Close()
	}

	return nil
}

func (ch *cleanHandle) Unregister() io.Closer {
	return ch.cln.unregister(ch.id)
}
