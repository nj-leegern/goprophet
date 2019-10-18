package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"net/url"
	"sync"
	"time"
)

/*
	sqlx agent
*/

var (
	dialects     sync.Map
	ErrSqlxAgent = errors.New("init sqlx agent failed")
)

type SqlxAgent struct {
	DB *sqlx.DB
	m  sync.Mutex
}

/* 实例化mysql代理 */
func NewMysqlAgent(ops ...Option) (*SqlxAgent, error) {
	option := parseOptions(ops...)
	dsName := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", option.username, option.password, option.hostname, option.database, url.QueryEscape("Asia/Shanghai"))
	return newInstance(DIALECT_MYSQL, dsName, int64(option.maxIdleConns), int64(option.maxOpenConns), option.connMaxLifetime.Nanoseconds())
}

/* 实例化sqlite3代理 */
func NewSqliteAgent(ops ...Option) (*SqlxAgent, error) {
	option := parseOptions(ops...)
	return newInstance(DIALECT_SQLITE3, option.database, int64(option.maxIdleConns), int64(option.maxOpenConns), option.connMaxLifetime.Nanoseconds())
}

/* 实例化postgres代理 */
func NewPostgresAgent(ops ...Option) (*SqlxAgent, error) {
	option := parseOptions(ops...)
	dsName := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full", option.username, option.password, option.hostname, option.database)
	return newInstance(DIALECT_PG, dsName, int64(option.maxIdleConns), int64(option.maxOpenConns), option.connMaxLifetime.Nanoseconds())
}

// 实例化sqlxAgent
func newInstance(dialectName, dsName string, args ...int64) (*SqlxAgent, error) {
	actual, loaded := dialects.LoadOrStore(dialectName, &SqlxAgent{})
	if dialect, ok := actual.(*SqlxAgent); ok {
		if loaded && dialect.DB != nil {
			return dialect, nil
		} else {
			dialect.m.Lock()
			defer dialect.m.Unlock()

			db, err := sqlx.Connect(dialectName, dsName)
			if err != nil {
				return nil, err
			}
			if len(args) >= 1 && args[0] > 0 {
				db.SetMaxIdleConns(int(args[0]))
			}
			if len(args) >= 2 && args[1] > 0 {
				db.SetMaxOpenConns(int(args[1]))
			}
			if len(args) >= 3 && args[2] > 0 {
				db.SetConnMaxLifetime(time.Duration(args[2]) * time.Nanosecond)
			}
			dialect.DB = db
		}
	}
	return nil, ErrSqlxAgent
}

// 解析配置项
func parseOptions(ops ...Option) options {
	defaultOps := defaultDBOptions()
	for _, apply := range ops {
		apply(&defaultOps)
	}
	return defaultOps
}
