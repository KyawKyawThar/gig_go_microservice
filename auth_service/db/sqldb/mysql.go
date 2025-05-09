package sqldb

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"reflect"
	"strings"
	"time"

	glog "auth_service/logs"
	"gorm.io/gorm"
)

type fieldInfo struct {
	columnName string
	valueType  string
}

var sqlDB *SqlDB

// DBConfig holds all database configuration
type DBConfig struct {
	SQLAddr     string
	EnableLog   bool
	MaxConn     int
	IdleConn    int
	MaxLifeTime int
}

type SqlDB struct {
	db *gorm.DB
}

// This line of code is a compile-time check to ensure that the SqlDB struct implements
//the Querier interface. Let me break down how it works:
// What it does
// The statement var _ Querier = (*SqlDB)(nil) does the following:
//
// 1. It declares a variable _ (blank identifier) of type Querier
// 2. It assigns a nil pointer of type *SqlDB to this variable
// 3. If *SqlDB doesn't implement all methods required by the Querier interface, the code won't compile

var _ Querier = (*SqlDB)(nil)

// DBFactory interface for database creation
type DBFactory interface {
	CreateDB() (*SqlDB, error)
}

// MySQLFactory implements DBFactory for MySQL
type MySQLFactory struct {
	Config *DBConfig
}

// NewMySQLFactory creates a new MySQL factory instance
func NewMySQLFactory(config DBConfig) *MySQLFactory {
	return &MySQLFactory{
		Config: &config,
	}
}

// CreateDB implements the DBFactory interface
func (f *MySQLFactory) CreateDB() (*SqlDB, error) {

	var logLevel logger.LogLevel

	if f.Config.EnableLog {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	glog.NewLogger().Infof("database address: %s", f.Config.SQLAddr)

	config := &gorm.Config{
		Logger: logger.New(glog.NewGormZapLogger(), logger.Config{
			SlowThreshold: 1000 * time.Microsecond,
			Colorful:      true,
			LogLevel:      logLevel,
		}),
	}

	tyDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN: f.Config.SQLAddr,
	}), config)

	if err != nil {

		tmpStr := fmt.Sprintf("conect=%s the sqldb, err: %v", f.Config.SQLAddr, err)
		glog.NewLogger().Errorf("%s", tmpStr)
		return nil, fmt.Errorf("%s", tmpStr)
	}

	db, err := tyDB.DB()

	if err != nil {
		tmpStr := fmt.Sprintf("conect=%s the sqldb, err: %v", f.Config.SQLAddr, err)
		glog.NewLogger().Errorf("%s", tmpStr)
		return nil, fmt.Errorf("%s", tmpStr)
	}

	db.SetMaxOpenConns(f.Config.MaxConn)
	db.SetMaxIdleConns(f.Config.IdleConn)
	db.SetConnMaxLifetime(time.Duration(f.Config.MaxLifeTime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(f.Config.MaxLifeTime) * time.Second)

	if err := db.Ping(); err != nil {
		tmpStr := fmt.Sprintf("Ping the db=%s err: %v", f.Config.SQLAddr, err)
		glog.NewLogger().Error(tmpStr)
		return nil, fmt.Errorf("%s", tmpStr)
	}
	return &SqlDB{db: tyDB}, nil
}

func (sql *SqlDB) MigrateDB(model any) error {

	return sql.db.AutoMigrate(model)
}

// InitSQLDB initializes the database using factory pattern
func InitSQLDB(sqlAddr string, enableLog bool, maxConn, idleConn, maxLifeTime int) error {

	config := DBConfig{
		SQLAddr:     sqlAddr,
		EnableLog:   enableLog,
		MaxConn:     maxConn,
		IdleConn:    idleConn,
		MaxLifeTime: maxLifeTime,
	}

	factory := NewMySQLFactory(config)
	db, err := factory.CreateDB()
	if err != nil {
		return err
	}

	sqlDB = db

	return nil
}

func GetBranchInsertSql(objs []any, tableName string) string {

	if len(objs) == 0 {
		return ""
	}

	// Get the type information once

	objType := reflect.TypeOf(objs[0])
	fieldNum := objType.NumField()

	// Pre-allocate slices with known capacity

	fields := make([]fieldInfo, 0, fieldNum)

	// Build field information in one pass

	for i := 0; i < fieldNum; i++ {
		field := objType.Field(i)

		fields[i].columnName = fmt.Sprintf("`%s`", GetColumnName(field.Tag.Get("gorm")))

		switch {
		case field.Type.Kind() == reflect.String:
			fields[i].valueType = "string"
		case strings.HasPrefix(field.Type.Name(), "uint"):
			fields[i].valueType = "uint"
		case strings.HasPrefix(field.Type.Name(), "int"):
			fields[i].valueType = "int"
		}
	}

	// Pre-allocate value list slice
	valueList := make([]string, len(objs))

	// Build values
	for i, obj := range objs {

		objV := reflect.ValueOf(obj)
		values := make([]string, fieldNum)

		for j, field := range fields {
			values[j] = GetFormatFeild(objV, j, field.valueType, "")
		}

		valueList[i] = "(" + strings.Join(values, ",") + ")"
	}
	// Extract column names for the SQL query
	columns := make([]string, fieldNum)
	for i, field := range fields {
		columns[i] = field.columnName
	}

	// Build final SQL
	return fmt.Sprintf("insert into `%s` (%s) values %s;",
		tableName,
		strings.Join(columns, ","),
		strings.Join(valueList, ","))

}

func Close() error {
	if sqlDB == nil || sqlDB.db == nil {
		return nil
	}
	db, err := sqlDB.db.DB()

	if err != nil {
		return err
	}
	return db.Close()

}

func (sql *SqlDB) Begin() *gorm.DB {

	tx := sql.db.Begin()
	return tx
}

func GetColumnName(jsonName string) string {
	for _, name := range strings.Split(jsonName, ";") {
		if strings.Index(name, "column") == -1 {
			continue
		}
		return strings.Replace(name, "column:", "", 1)
	}
	return ""
}

func GetFormatFeild(objV reflect.Value, index int, t string, sep string) string {
	switch t {
	case "string":
		return fmt.Sprintf("'%s'%s", objV.Field(index).String(), sep)
	case "uint":
		return fmt.Sprintf("%d%s", objV.Field(index).Uint(), sep)
	case "int":
		return fmt.Sprintf("%d%s", objV.Field(index).Int(), sep)
	default:
		return ""
	}
}

func (sql *SqlDB) UpdateDB(row any) error {

	err := sql.db.Save(row).Error
	if err != nil {
		return err
	}
	return nil
}

func (sql *SqlDB) InsertDB(row any) error {
	if sql.db == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Create(row).Error

}

func (sql *SqlDB) SingleUpdateDB(mods any, upMap map[string]any, query any, args ...any) error {
	if len(upMap) == 0 {
		return errors.New("update map cannot be empty")
	}

	if mods == nil || query == nil {
		return errors.New("model or update map cannot be nil")
	}

	return sql.db.Model(mods).Where(query, args...).Updates(upMap).Error
}

func (sql *SqlDB) SingleUpdateSc(mods any, query, oldValue, newValue any, field, fieldExpression string) error {
	if mods == nil {
		return errors.New("model cannot be nil")
	}
	if len(field) == 0 || len(fieldExpression) == 0 {
		return errors.New("field and fieldExpression cannot be empty")
	}

	if sqlDB == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Model(mods).Where(query, oldValue).Update(field, gorm.Expr(fieldExpression, newValue)).Error

}

func (sql *SqlDB) FetchOne(mods any, tableName string, cols []string, query any, args ...any) error {
	if sql == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Table(tableName).Select(cols).Where(query, args...).First(mods).Error
}

func (sql *SqlDB) FetchAll(mods any, tableName string, cols []string, query any, args ...any) error {
	if sql == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Table(tableName).Select(cols).Where(query, args...).Find(mods).Error
}

func (sql *SqlDB) FetchAllWithout(mods any, tableName string, cols []string) error {
	if sql == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Table(tableName).Select(cols).Find(mods).Error
}

func (sql *SqlDB) SingleUpdate(tableName string, query any, upMap map[string]interface{}, args ...interface{}) error {
	if sql == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Table(tableName).Where(query, args...).Updates(upMap).Error
}

func (sql *SqlDB) ExecSQLRepScan(mods any, query string, args ...interface{}) error {
	if sql == nil {
		return errors.New("sqlDB is not initialized")
	}
	return sql.db.Raw(query, args...).Scan(mods).Error
}
