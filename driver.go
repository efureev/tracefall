package traceFall

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"sort"
	"sync"
)

type Driver interface {
	Open(map[string]string) (interface{}, error)
	Send(log *Log) (ResponseCmd, error)
	RemoveThread(id uuid.UUID) (ResponseCmd, error)
	RemoveByTags(tags Tags) (ResponseCmd, error)
	GetLog(id uuid.UUID) (ResponseLog, error)
	GetThread(id uuid.UUID) (ResponseThread, error)
	Truncate(ind string) (ResponseCmd, error)
}

var (
	drivers   = make(map[string]Driver)
	driversMu sync.RWMutex
)

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("log: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("log: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]Driver)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

type DB struct {
	connector Connector
	stop      func() // stop cancels the connection opener and the session resetter.
}

func (d *DB) conn(ctx context.Context) (interface{}, error) {
	return d.connector.Connect(ctx)
}

func (d *DB) Send(log *Log) (ResponseCmd, error) {
	return d.connector.Driver().Send(log)
}

func (d *DB) RemoveThread(id uuid.UUID) (ResponseCmd, error) {
	return d.connector.Driver().RemoveThread(id)
}

func (d *DB) RemoveByTags(tags Tags) (ResponseCmd, error) {
	return d.connector.Driver().RemoveByTags(tags)
}

func (d *DB) Get(id uuid.UUID) (ResponseLog, error) {
	return d.connector.Driver().GetLog(id)
}

func (d *DB) GetThread(id uuid.UUID) (ResponseThread, error) {
	return d.connector.Driver().GetThread(id)
}

func (d *DB) Truncate(ind string) (ResponseCmd, error) {
	return d.connector.Driver().Truncate(ind)
}

func (d *DB) Driver() Driver {
	return d.connector.Driver()
}

type Connector interface {
	Connect(ctx context.Context) (interface{}, error)
	Driver() Driver
}

type drvConnector struct {
	params map[string]string
	driver Driver
}

func (t drvConnector) Connect(_ context.Context) (interface{}, error) {
	return t.driver.Open(t.params)
}

func (t drvConnector) Driver() Driver {
	return t.driver
}

func Open(driverName string, connectParams map[string]string) (*DB, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("tracer: unknown driver %q (forgotten import?)", driverName)
	}

	db, err := OpenDB(drvConnector{params: connectParams, driver: driveri})
	if err != nil {
		return db, err
	}
	return db, nil
}

func OpenDB(c drvConnector) (*DB, error) {
	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{
		connector: c,
		stop:      cancel,
	}

	_, err := db.conn(ctx)

	return db, err
}
