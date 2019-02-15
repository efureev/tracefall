package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/efureev/tracefall"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Params struct {
	Host, User, Password, DbName, TableName string
}

func (p *Params) set(params map[string]string) {
	p.Host = params[`host`]
	p.User = params[`user`]
	p.Password = params[`pwd`]
	p.DbName = params[`db`]
	p.TableName = params[`table`]
}

type DriverPostgres struct {
	params Params
}

func (d *DriverPostgres) initDb() *sql.DB {
	db, err := sql.Open("postgres", d.params.pgConnectionStr())
	if err != nil {
		panic(err)
	}

	return db
}

func (p Params) pgConnectionStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		p.User,
		p.Password,
		p.Host,
		p.DbName,
	)
}

func (d DriverPostgres) Send(l *tracefall.Log) (tracefall.ResponseCmd, error) {
	db := d.initDb()
	defer db.Close()

	query := `INSERT INTO "` + d.params.TableName + `" (
			"id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING "id";`

	resp := tracefall.NewResponse(l)

	stmt, err := db.Prepare(query)

	if err != nil {
		return *resp.SetError(err).ToCmd(), err
	}
	defer stmt.Close()

	var (
		parentID, errLog *string
		te               *int64
	)

	if l.Parent != nil {
		idStr := l.Parent.ID.String()
		parentID = &idStr
	} else {
		parentID = nil
	}

	if l.Error != nil {
		errStr := l.Error.Error()
		errLog = &errStr
	} else {
		errLog = nil
	}

	if l.TimeEnd != nil {
		teInt := l.TimeEnd.UnixNano()
		te = &teInt
	}

	row := db.QueryRow(query, l.ID.String(), l.Thread.String(), parentID, l.App, l.Name, l.Time.UnixNano(), te,
		l.Environment, pq.Array(l.Tags), l.Notes.ToJSON(), l.Data.ToJSON(), errLog, l.Result, l.Finish)

	var id string

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return *resp.SetError(err).ToCmd(), err
	case nil:
		return *resp.Success().SetID(id).ToCmd(), err
	default:
		panic(err)
	}
}

func (d DriverPostgres) RemoveThread(id uuid.UUID) (tracefall.ResponseCmd, error) {
	db := d.initDb()
	defer db.Close()

	query := `DELETE FROM "` + d.params.TableName + `" WHERE thread = $1`

	resp := tracefall.NewResponse(id)

	_, err := db.Exec(query, id.String())
	if err != nil {
		return *resp.SetError(err).ToCmd(), err
	}

	return *resp.Success().ToCmd(), nil
}

func (d DriverPostgres) RemoveByTags(tags tracefall.Tags) (tracefall.ResponseCmd, error) {
	db := d.initDb()
	defer db.Close()

	query := `DELETE FROM "` + d.params.TableName + `" WHERE $1 <@ "tags"`

	resp := tracefall.NewResponse(tags)

	_, err := db.Exec(query, pq.Array(tags))
	if err != nil {
		return *resp.SetError(err).ToCmd(), err
	}

	return *resp.Success().ToCmd(), nil
}

func (d DriverPostgres) getListResult(rows *sql.Rows) ([]*tracefall.Log, error) {
	var logList []*tracefall.Log

	for rows.Next() {
		var (
			l                   = tracefall.Log{}
			idStr, threadStr    string
			parentPtr, errorPtr *string
			notesStr, dataStr   []byte
			ts                  int64
			te                  *int64
			t                   pq.StringArray
		)
		err := rows.Scan(&idStr, &threadStr, &parentPtr, &l.App, &l.Name, &ts, &te, &l.Environment, &t, &notesStr, &dataStr, &errorPtr, &l.Result, &l.Finish)
		if err != nil {
			return nil, err
		}

		uid, err := uuid.FromString(idStr)
		if err != nil {
			return nil, err
		}
		l.ID = uid
		thid, err := uuid.FromString(threadStr)
		if err != nil {
			return nil, err
		}
		l.Thread = thid

		if parentPtr != nil {
			pid, err := uuid.FromString(*parentPtr)
			if err != nil {
				return nil, err
			}
			l.SetParentID(pid)
		}

		l.Data.FromJSON(dataStr)
		l.Notes.FromJSON(notesStr)
		l.Tags = tracefall.Tags(t)

		l.Time = time.Unix(0, ts)

		if te != nil {
			t := time.Unix(0, *te)
			l.TimeEnd = &t
		}

		if errorPtr != nil {
			l.Error = errors.New(*errorPtr)
		}

		logList = append(logList, &l)
	}

	return logList, nil
}

func (d DriverPostgres) getListLogJSONResult(rows *sql.Rows) ([]*tracefall.LogJSON, error) {
	var list []*tracefall.LogJSON

	for rows.Next() {
		var (
			l                   = tracefall.LogJSON{}
			idStr, threadStr    string
			parentPtr, errorPtr *string
			notesStr, dataStr   []byte
			ts                  int64
			te                  *int64
			t                   pq.StringArray
		)
		err := rows.Scan(&idStr, &threadStr, &parentPtr, &l.App, &l.Name, &ts, &te, &l.Environment, &t, &notesStr, &dataStr, &errorPtr, &l.Result, &l.Finish)
		if err != nil {
			return nil, err
		}

		uid, err := uuid.FromString(idStr)
		if err != nil {
			return nil, err
		}
		l.ID = uid
		thid, err := uuid.FromString(threadStr)
		if err != nil {
			return nil, err
		}
		l.Thread = thid

		l.Parent = parentPtr

		l.Data.FromJSON(dataStr)
		l.Notes.FromJSON(notesStr)
		l.Tags = tracefall.Tags(t)

		l.Time = ts
		l.TimeEnd = te
		l.Error = errorPtr

		list = append(list, &l)
	}

	return list, nil
}

func (d DriverPostgres) getListByThread(id uuid.UUID) ([]*tracefall.LogJSON, error) {
	query := `SELECT "id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish" FROM "` + d.params.TableName + `" WHERE "thread"=$1`

	db := d.initDb()
	defer db.Close()

	rows, err := db.Query(query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.getListLogJSONResult(rows)
}

func (d DriverPostgres) GetLastRootList(limit int) ([]*tracefall.Log, error) {
	query := `SELECT "id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish" 
		FROM "` + d.params.TableName + `"
		WHERE parent IS NULL
		ORDER BY time 
		LIMIT $1`

	db := d.initDb()
	defer db.Close()

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.getListResult(rows)
}

func (d DriverPostgres) GetLastThreadList(limit int) ([]*tracefall.Log, error) {
	query := `SELECT "id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish"
		FROM "` + d.params.TableName + `" t
		where t.thread IN (SELECT "id" pid
			FROM "` + d.params.TableName + `"
			WHERE parent is null
			ORDER BY time DESC
			LIMIT $1)`

	db := d.initDb()
	defer db.Close()

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return d.getListResult(rows)
}

func (d DriverPostgres) GetThread(id uuid.UUID) (tracefall.ResponseThread, error) {
	resp := tracefall.NewResponse(id)
	list, err := d.getListByThread(id)
	if err != nil {
		return *resp.SetError(err).ToThread(tracefall.ThreadFromList(list)), err
	}

	return *resp.Success().ToThread(tracefall.ThreadFromList(list)), nil
}

/*
func (d DriverPostgres) GetLog(id uuid.UUID) (tracefall.ResponseLog, error) {
	query := `SELECT "id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish" FROM "` + d.params.TableName + `" WHERE "id"=$1`

	var (
		l                   = tracefall.Log{}
		idStr, threadStr    string
		parentPtr, errorPtr *string
		notesStr, dataStr   string
		ts                  int64
		te                  *int64
		t                   pq.StringArray
	)

	db := d.initDb()
	defer db.Close()

	resp := tracefall.NewResponse(id)

	row := db.QueryRow(query, id)
	switch err := row.Scan(&idStr, &threadStr, &parentPtr, &l.App, &l.Name, &ts, &te, &l.Environment, &t, &notesStr, &dataStr, &errorPtr, &l.Result, &l.Finish); err {
	case sql.ErrNoRows:
		return *resp.SetError(errors.New(`not found`)).ToLog(nil), nil
	case nil:
		uid, err := uuid.FromString(idStr)
		if err != nil {
			return *resp.SetError(err).ToLog(nil), nil
		}
		l.ID = uid
		thid, err := uuid.FromString(threadStr)
		if err != nil {
			return *resp.SetError(err).ToLog(nil), nil
		}
		l.Thread = thid

		if parentPtr != nil {
			pid, err := uuid.FromString(*parentPtr)
			if err != nil {
				return *resp.SetError(err).ToLog(nil), nil
			}
			l.SetParentID(pid)
		}

		l.Data.FromJSON(dataStr)
		l.Notes.FromJSON(notesStr)
		l.Tags = tracefall.Tags(t)

		l.Time = time.Unix(0, ts)

		if te != nil {
			t := time.Unix(0, *te)
			l.TimeEnd = &t
		}

		if errorPtr != nil {
			l.Error = errors.New(*errorPtr)
		}

		return *resp.Success().ToLog(), nil
	default:
		return *resp.SetError(err).ToLog(nil), nil
	}
}
*/

func (d DriverPostgres) GetLog(id uuid.UUID) (tracefall.ResponseLog, error) {
	query := `SELECT "id", "thread", "parent", "app", "name", "time", "time_end", "env", "tags", "notes", "data", "error", "result", "finish" FROM "` + d.params.TableName + `" WHERE "id"=$1`

	var (
		l                   = tracefall.LogJSON{}
		idStr, threadStr    string
		parentPtr, errorPtr *string
		notesStr, dataStr   []byte
		ts                  int64
		te                  *int64
		t                   pq.StringArray
	)

	db := d.initDb()
	defer db.Close()

	resp := tracefall.NewResponse(id)

	row := db.QueryRow(query, id)
	switch err := row.Scan(&idStr, &threadStr, &parentPtr, &l.App, &l.Name, &ts, &te, &l.Environment, &t, &notesStr, &dataStr, &errorPtr, &l.Result, &l.Finish); err {
	case sql.ErrNoRows:
		e := errors.New(`not found`)
		return *resp.SetError(e).ToLog(nil), e
	case nil:
		uid, err := uuid.FromString(idStr)
		if err != nil {
			return *resp.SetError(err).ToLog(nil), nil
		}
		l.ID = uid
		thid, err := uuid.FromString(threadStr)
		if err != nil {
			return *resp.SetError(err).ToLog(nil), nil
		}
		l.Thread = thid

		l.Parent = parentPtr

		l.Data.FromJSON(dataStr)
		l.Notes.FromJSON(notesStr)
		l.Tags = tracefall.Tags(t)

		l.Time = ts
		l.TimeEnd = te

		l.Error = errorPtr

		return *resp.Success().ToLog(&l), nil
	default:
		return *resp.SetError(err).ToLog(nil), nil
	}
}

// Create table for tracer
func (d DriverPostgres) CreateTable() error {
	db := d.initDb()
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS "` + d.params.TableName + `" (
  id          UUID primary key,
  thread      UUID NOT NULL,
  parent      UUID NULL,
  app         VARCHAR(100) NOT NULL,
  name        VARCHAR(255) NOT NULL,
  time        bigint NOT NULL,
  time_end    bigint NULL,
  env         VARCHAR(50) default 'dev',
  tags        text[],
  notes       jsonb default '[]',
  data        jsonb default '[]',
  error       text NULL,
  result      boolean NOT NULL default false,
  finish      boolean NOT NULL default false,
  created     timestamp without time zone default now()
);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Create table for tracer
func (d DriverPostgres) InstallIndex() error {
	db := d.initDb()
	defer db.Close()

	query := `	
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_time_idx" ON "` + d.params.TableName + `"("time");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_finish_idx" ON "` + d.params.TableName + `"("finish");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_result_idx" ON "` + d.params.TableName + `"("result");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_result_idx" ON "` + d.params.TableName + `"("result");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_env_idx" ON "` + d.params.TableName + `"("env");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_app_idx" ON "` + d.params.TableName + `"("app");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_thread_idx" ON "` + d.params.TableName + `"("thread");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_parent_idx" ON "` + d.params.TableName + `"("parent");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_data_idx" ON "` + d.params.TableName + `" USING GIN ("data");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_notes_idx" ON "` + d.params.TableName + `" USING GIN ("notes");
	CREATE INDEX IF NOT EXISTS "` + d.params.TableName + `_tags_idx" ON "` + d.params.TableName + `" USING GIN ("tags");
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Erase table
func (d DriverPostgres) DropTable() error {
	db := d.initDb()
	defer db.Close()

	query := `DROP TABLE IF EXISTS ` + d.params.TableName + `;`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d DriverPostgres) Truncate(ind string) (tracefall.ResponseCmd, error) {
	db := d.initDb()
	defer db.Close()

	resp := tracefall.NewResponse(ind).GenerateID()

	if ind == `` {
		ind = d.params.TableName
	}
	query := `TRUNCATE TABLE ` + ind + `;`

	_, err := db.Exec(query)
	if err != nil {
		return *resp.SetError(err).ToCmd(), err
	}

	return *resp.Success().ToCmd(), nil
}

func (d *DriverPostgres) Open(params map[string]string) (interface{}, error) {
	d.params.set(params)
	db := d.initDb()

	defer db.Close()

	if err := db.Ping(); err != nil {
		e := errors.New(fmt.Sprintf("Couldn't ping postgre database (%s): error: %s", d.params.DbName, err))
		panic(e)
	}

	err := d.CreateTable()

	if err != nil {
		panic(err)
	}

	err = d.InstallIndex()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	tracefall.Register("postgres", &DriverPostgres{})
}

func GetConnParams(host, db, table, user, pwd string) map[string]string {
	return map[string]string{`host`: host, `db`: db, `table`: table, `user`: user, `pwd`: pwd}
}
