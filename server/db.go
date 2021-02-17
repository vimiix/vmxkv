package server

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/vimiix/vmxkv/internal/hotpool"
	"github.com/vimiix/vmxkv/internal/logging"
)

var (
	defaultHotPoolSize     = 512
	defaultMonitorDuration = 5 * time.Second
	defaultBPTreeDegree    = 4
)

type DBInterface interface {
	Put(key, value uint64) error
	Get(key uint64) (uint64, error)
	Del(key uint64) error
	Close()
	List(k1, k2 uint64, f func(uint64, uint64) bool) error
}

type DB struct {
	path     string
	version  int64
	file     *os.File
	fileLock *sync.Mutex
	tree     *BPTree
	hotPool  *hotpool.HotPool

	monitorCloseCh  chan int
	monitorDuration time.Duration
}

func (db *DB) Put(key, value uint64) (err error) {
	err = db.tree.Insert(key, value)
	if err != nil {
		logging.Error(fmt.Sprintf("put key:value [%d:%d] error: %s", key, value, err))
	}
	return
}

func (db *DB) Get(key uint64) (value uint64, err error) {
	if value, ok := db.hotPool.Get(key); ok {
		return value, nil
	}
	value, err = db.tree.Find(key)
	if err == nil {
		db.hotPool.Put(key, value)
	}
	return
}

func (db *DB) Del(key uint64) (err error) {
	db.hotPool.Delete(key)
	return db.tree.Delete(key)
}

func (db *DB) List(k1, k2 uint64, f func(uint64, uint64) bool) (err error) {
	return db.tree.RangeFind(k1, k2, f)
}

func (db *DB) Path() string {
	return db.path
}

func (db *DB) String() string {
	return fmt.Sprintf("DB<%q>", db.path)
}

func OpenDB(path string, mode os.FileMode) (db *DB, err error) {
	version := time.Now().Unix()
	db = &DB{
		path:            path,
		version:         version,
		hotPool:         hotpool.NewHotPool(defaultHotPoolSize),
		fileLock:        &sync.Mutex{},
		monitorCloseCh:  make(chan int),
		monitorDuration: defaultMonitorDuration,
	}

	if db.file, err = os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, mode); err != nil {
		return nil, err
	}
	var stat syscall.Statfs_t
	if err = syscall.Statfs(db.path, &stat); err != nil {
		return nil, err
	}
	blockSize := uint64(stat.Bsize)
	if blockSize == 0 {
		return nil, errors.New("db file blocksize is 0")
	}
	db.tree = NewBPTree(defaultBPTreeDegree)
	if fileInfo, err := db.file.Stat(); err != nil {
		return nil, err
	} else if fileInfo.Size() != 0 {
		if err = db.load(); err != nil {
			return nil, err
		}
	}

	go db.startMonitor()
	return
}

// Sync 同步数据到文件
// TODO bad practice
func (db *DB) Sync() (err error) {
	db.fileLock.Lock()
	defer db.fileLock.Unlock()
	if err = db.tree.Dump(db.file); err != nil {
		return
	}
	err = db.file.Sync()
	if err == nil {
		db.version = time.Now().Unix()
	}
	return
}

// load 从已有的数据库文件加载数据
func (db *DB) load() (err error) {
	db.tree, err = NewFromFile(db.file)
	return
}

// close 正常关闭数据库
func (db *DB) Close() {
	close(db.monitorCloseCh)

	if db.file != nil {
		db.Sync()
		db.file.Close()
	}
	return
}

// startMonitor 监控内存数据变化，并定期持久化到文件
func (db *DB) startMonitor() {
	for {
		select {
		case <-db.monitorCloseCh:
			return
		case <-time.Tick(db.monitorDuration):
			if db.version < db.tree.Version {
				db.Sync()
			}
		}
	}
}
