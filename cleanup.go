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

	// Unregister 从 VM 清理列表中移除资源，不触发关闭。
	// 返回原始 io.Closer，若资源已不存在则返回 nil。
	Unregister() io.Closer
}

// cleanManager 管理一组可关闭资源。
// register 添加资源并返回唯一标识 id。
// unregister 按 id 移除资源但不关闭。
// closeAll 按注册顺序逆序关闭所有资源，此后禁止再次注册。
type cleanManager interface {
	// register 注册待释放资源，返回其唯一 id 和是否成功。
	// 若管理器已执行清理（done=true），则返回 (0, false)。
	register(c io.Closer) (id int64, succeed bool)

	// unregister 按 id 注销资源，返回原 io.Closer；若不存在则返回 nil。
	unregister(id int64) io.Closer

	// closeAll 逆序关闭所有已注册资源，并将管理器标记为已执行。
	closeAll()
}

// cleanMapManager 是 cleanManager 的基于 map 实现。
// 使用自增 id 作为 map key，线程安全，支持并发注册/注销/关闭。
type cleanMapManager struct {
	mutex sync.Mutex          // 保护以下所有字段的并发访问
	seq   int64               // 自增序列号，用作资源的唯一 id
	done  bool                // 是否已执行 closeAll（执行后禁止再注册）
	elems map[int64]io.Closer // id -> 资源的映射
}

// register 注册资源并返回递增 id。
// 若 cleanAll 已调用则返回 (0, false) 表示注册失败。
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

// unregister 按 id 移除并返回资源；若 id 不存在返回 nil。
func (cm *cleanMapManager) unregister(id int64) io.Closer {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cl := cm.elems[id]
	delete(cm.elems, id)

	return cl
}

// closeAll 按 id 降序（逆注册顺序）依次调用所有资源的 Close 方法。
// 执行完毕后设置 done=true，并清空 elems，防止重复执行或后续注册。
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

	// 收集所有 id 并降序排列，确保后进先出（LIFO）的关闭顺序
	ids := make([]int64, 0, len(elems))
	for id := range elems {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a, b int64) int {
		return cmp.Compare(b, a) // 降序，id 大的先关闭
	})
	for _, id := range ids {
		if cl := elems[id]; cl != nil {
			_ = cl.Close()
		}
	}
}

// fallbackCleanHandle 是 CleanHandle 的实现。
// 正常情况（id!=0）通过 cleanManager 管理；
// 降级情况（fb!=nil）直接持有资源，绕过管理器。
type fallbackCleanHandle struct {
	id int64        // 资源在 cleanManager 中的 id（正常情况）
	cm cleanManager // 资源管理器引用
	fb io.Closer    // 降级持有的资源（管理器不可用时直接用此字段）
}

// Close 关闭资源。
// 若 fb 不为 nil，直接关闭 fb；
// 否则从管理器中移除并关闭对应资源。
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

// Unregister 仅从管理器中注销资源，不执行关闭，返回原 io.Closer。
func (fch *fallbackCleanHandle) Unregister() io.Closer {
	return fch.cm.unregister(fch.id)
}
