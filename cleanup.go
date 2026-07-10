package jsvm

import (
	"cmp"
	"io"
	"slices"
	"sync"
)

// CleanHandle 控制已注册资源的生命周期。
// Close 关闭资源并从 VM 清理列表中移除。
// Unregister 仅从清理列表中移除资源，不关闭。
type CleanHandle interface {
	io.Closer

	// Unregister 从 VM 清理列表中移除资源。
	// 返回原始 io.Closer；若已被移除则返回 nil。
	Unregister() io.Closer
}

// cleanManager 管理一组可关闭资源。
// register 添加资源并返回唯一标识 id。
// unregister 按 id 移除资源但不关闭。
// closeAll 按注册顺序逆序关闭所有资源，此后禁止再次注册。
type cleanManager interface {
	// register 注册待释放资源，返回其唯一 id 和是否成功。
	register(c io.Closer) (id int64, succeed bool)

	// unregister 按 id 注销资源，返回原 io.Closer；若不存在则返回 nil。
	unregister(id int64) io.Closer

	// closeAll 逆序关闭所有已注册资源，并将管理器标记为已执行。
	closeAll()
}

type cleanMapManager struct {
	mutex sync.Mutex
	seq   int64
	done  bool
	elems map[int64]io.Closer
}

func (cm *cleanMapManager) register(c io.Closer) (int64, bool) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.done {
		return 0, false
	}
	if cm.elems == nil {
		cm.elems = make(map[int64]io.Closer, 4)
	}

	cm.seq++
	id := cm.seq
	cm.elems[id] = c

	return id, true
}

func (cm *cleanMapManager) unregister(id int64) io.Closer {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cl := cm.elems[id]
	delete(cm.elems, id)

	return cl
}

func (cm *cleanMapManager) closeAll() {
	cm.mutex.Lock()
	done := cm.done
	elems := cm.elems
	cm.done = true
	cm.elems = nil
	cm.mutex.Unlock()

	if done {
		return
	}

	ids := make([]int64, 0, len(elems))
	for id := range elems {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a, b int64) int {
		return cmp.Compare(b, a)
	})
	for _, id := range ids {
		if cl := elems[id]; cl != nil {
			_ = cl.Close()
		}
	}
}

type fallbackCleanHandle struct {
	id int64
	cm cleanManager
	fb io.Closer
}

func (fch *fallbackCleanHandle) Close() error {
	cl := fch.fb
	if cl == nil {
		cl = fch.cm.unregister(fch.id)
	}
	if cl != nil {
		return cl.Close()
	}

	return nil
}

func (fch *fallbackCleanHandle) Unregister() io.Closer {
	return fch.cm.unregister(fch.id)
}
