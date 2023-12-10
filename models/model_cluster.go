package models

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/core"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Cluster 群集的模型
type Cluster struct {
	ClusterID   uint   `xorm:"'cluster_id' notnull int pk autoincr"                      valid:"-"                                json:"cluster_id"  gqlgen:"-"`        //
	UUID        string `xorm:"'uuid' notnull char(36) unique(unique_1)"                  valid:"-"                                json:"uuid"        gqlgen:"UUID"`     //
	Host        string `xorm:"'host' notnull varchar(150) unique(unique_2)"              valid:"required,length(1|100),alphanum"  json:"host"        gqlgen:"Host"`     //
	IP          string `xorm:"'ip' notnull varchar(15) unique(unique_3)"                 valid:"required,int,range(0|4294967295)" json:"ip"          gqlgen:"Ip"`       //
	Port        uint16 `xorm:"'port' notnull smallint unique(unique_2) unique(unique_3)" valid:"required,port"                    json:"port"        gqlgen:"Port"`     //
	Alias       string `xorm:"'alias' notnull varchar(75) unique(unique_4)"              valid:"required,runelength(1|100)"       json:"alias"       gqlgen:"Alias"`    //
	User        string `xorm:"'user' notnull varchar(50)"                                valid:"required,length(1|50),alphanum"   json:"user"        gqlgen:"User"`     //
	Password    []byte `xorm:"'password' notnull varbinary(48)"                          valid:"required"                         json:"-"           gqlgen:"-"`        // 双向加密
	FingerPrint []byte `xorm:"'fingerprint' notnull varbinary(16)"                       valid:"required"                         json:"-"           gqlgen:"-"`        //
	Status      uint8  `xorm:"'status' notnull tinyint"                                  valid:"required,matches(^[0-9]$)"        json:"status"      gqlgen:"Status"`   //
	Version     int    `xorm:"'version'"                                                 valid:"-"                                json:"version"     gqlgen:"-"`        //
	UpdateAt    uint   `xorm:"'update_at' notnull int"                                   valid:"-"                                json:"update_at"   gqlgen:"UpdateAt"` //
	CreateAt    uint   `xorm:"'create_at' notnull int"                                   valid:"-"                                json:"create_at"   gqlgen:"CreateAt"` //
}

// TableName 结构体到数据库表名称的映射
func (m *Cluster) TableName() string {
	return "mm_clusters"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Cluster) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Host = strings.TrimSpace(m.Host)
	m.IP = strings.TrimSpace(m.IP)
	m.User = strings.TrimSpace(m.User)
	m.FingerPrint = m.Sha1(fmt.Sprintf("%s-%s-%d-%s-%s",
		m.Host,
		m.IP,
		m.Port,
		m.User,
		m.Password))
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Cluster) BeforeUpdate() {
	m.Host = strings.TrimSpace(m.Host)
	m.IP = strings.TrimSpace(m.IP)
	m.User = strings.TrimSpace(m.User)
	fp := m.Sha1(fmt.Sprintf("%s-%s-%d-%s-%s",
		m.Host,
		m.IP,
		m.Port,
		m.User,
		m.Password))
	if bytes.Compare(fp, m.FingerPrint) != 0 {
		m.FingerPrint = fp
	}
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Cluster) AfterSet(colName string, cell xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Cluster) String() string {
	return fmt.Sprintf("uuid: %s, host: %s, alias: %s, ip: %s, port: %d, user: %s, status: %d",
		m.UUID,
		m.Host,
		m.Alias,
		m.IP,
		m.Port,
		m.User,
		m.Status,
	)
}

// Sha1 对字符串进行sha1 计算
func (m *Cluster) Sha1(data string) []byte {
	t := sha1.New()
	io.WriteString(t, data)
	return t.Sum(nil)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Cluster) IsNode() {}

// 创建时间
func (m *Cluster) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Cluster) GetUpdateAt() *uint {
	return &m.UpdateAt
}

// IsSearchable GraphQL的基类需要实现的接口，暂时不动
func (Cluster) IsSearchable() {}

// Connect 连接到群集的指定数据库
func (m *Cluster) Connect(name string, passwd func(c *Cluster) []byte) (engine *xorm.Engine, err error) {
L:
	for {
		format := "%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true"
		addr := fmt.Sprintf(format, m.User, string(passwd(m)), m.IP, m.Port, name)
		engine, err = xorm.NewEngine("mysql", addr)
		if err != nil {
			break L
		}
		engine.SetConnMaxLifetime(time.Duration(60) * time.Second)
		err = engine.DB().Ping()
		if err != nil {
			break L
		}

		break L
	}

	return
}

// Execute 执行语句
func (m *Cluster) Execute(dbname string, passwd func(c *Cluster) []byte, sql string) (err error) {
L:
	for {
		var engine *xorm.Engine
		if engine, err = m.Connect(dbname, passwd); err != nil {
			return err
		}

		_, err = engine.Exec(sql)

		break L
	}

	return
}

// Databases 获取群集上所有数据库的信息
func (m *Cluster) Databases(passwd func(c *Cluster) []byte) (databases []Database, err error) {
L:
	for {
		engine := &xorm.Engine{}
		if engine, err = m.Connect("information_schema", passwd); err != nil {
			break L
		}

		rows := []map[string]string{}
		sql := `
		SELECT SCHEMA_NAME,
		       DEFAULT_CHARACTER_SET_NAME,
		       DEFAULT_COLLATION_NAME
		  FROM SCHEMATA
		 WHERE SCHEMA_NAME NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
				;
				`
		if rows, err = engine.QueryString(sql); err != nil {
			break L
		}

		for _, row := range rows {
			database := Database{
				Name:    row["SCHEMA_NAME"],
				Charset: row["DEFAULT_CHARACTER_SET_NAME"],
				Collate: row["DEFAULT_COLLATION_NAME"],
			}
			databases = append(databases, database)
		}

		break L
	}

	return
}

// Metadata 获取群集上某一个具体的数据库的元数据信息
func (m *Cluster) Metadata(database string, passwd func(c *Cluster) []byte) (tables map[string][]*core.Table, err error) {
L:
	for {
		var engine *xorm.Engine

		engine, err = m.Connect("information_schema", passwd)
		if err != nil {
			break L
		}
		defer engine.Close()

		args := []interface{}{}
		sql := ""
		if database == "*" {
			sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, " +
				"`AUTO_INCREMENT`, `TABLE_COMMENT`, `TABLE_COLLATION` " +
				"FROM `INFORMATION_SCHEMA`.`TABLES` " +
				"WHERE 1 = ? AND `ENGINE` IN ('MyISAM', 'InnoDB', 'TokuDB', 'RocksDB') " +
				"ORDER BY `TABLE_SCHEMA`;"
			args = []interface{}{1}
		} else {
			sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `ENGINE`, `TABLE_ROWS`, " +
				"`AUTO_INCREMENT`, `TABLE_COMMENT`, `TABLE_COLLATION` " +
				"FROM `INFORMATION_SCHEMA`.`TABLES` " +
				"WHERE `TABLE_SCHEMA` = ? AND `ENGINE` IN ('MyISAM', 'InnoDB', 'TokuDB', 'RocksDB') " +
				"ORDER BY `TABLE_SCHEMA`;"
			args = []interface{}{database}
		}
		rows, err := engine.DB().Query(sql, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		tables = make(map[string][]*core.Table, 0)
		for rows.Next() {
			table := core.NewEmptyTable()
			var db, name, engine, tableRows, comment, charset, collate string
			var autoIncr *string
			err = rows.Scan(&db, &name, &engine, &tableRows, &autoIncr, &comment, &collate)
			if err != nil {
				return nil, err
			}

			charset = strings.Split(collate, "_")[0]
			table.Name = name
			table.Comment = comment
			table.StoreEngine = engine
			table.Charset = charset
			table.Collate = collate
			if _, ok := tables[db]; ok {
				tables[db] = append(tables[db], table)
			} else {
				tables[db] = []*core.Table{table}
			}
		}

		allColumns, _ := m.columns(engine, database)
		allIndexes, _ := m.indexes(engine, database)
		for db, tbls := range tables {
			for _, tbl := range tbls {
				k := fmt.Sprintf("%s.%s", db, tbl.Name)
				cols := allColumns[k]
				inds := allIndexes[k]
				for _, col := range cols {
					tbl.AddColumn(col)
					tbl.Indexes = inds
				}
				for _, idx := range inds {
					for _, name := range idx.Cols {
						if col := tbl.GetColumn(name); col != nil {
							col.Indexes[idx.Name] = idx.Type
						} else {
							return nil, fmt.Errorf("Unknown col %s in index %v of table %v, columns %v", name, idx.Name, tbl.Name, tbl.ColumnsSeq())
						}
					}
				}
			}
		}
		break
	}

	return
}

func (m *Cluster) columns(engine *xorm.Engine, database string) (cols map[string]map[string]*core.Column, err error) {
	args := []interface{}{}
	sql := ""
	if database == "*" {
		sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `IS_NULLABLE`, " +
			"`COLUMN_DEFAULT`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`,`COLUMN_COMMENT` " +
			"FROM `INFORMATION_SCHEMA`.`COLUMNS` " +
			"WHERE 1 = ? AND `TABLE_SCHEMA` NOT IN ('mysql', 'sys', 'information_schema') " +
			"ORDER BY `TABLE_SCHEMA`, `TABLE_NAME`;"
		args = []interface{}{1}
	} else {
		sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `IS_NULLABLE`, " +
			"`COLUMN_DEFAULT`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`,`COLUMN_COMMENT` " +
			"FROM `INFORMATION_SCHEMA`.`COLUMNS` " +
			"WHERE `TABLE_SCHEMA` = ? AND `TABLE_SCHEMA` NOT IN ('mysql', 'sys', 'information_schema') " +
			"ORDER BY `TABLE_SCHEMA`, `TABLE_NAME`;"
		args = []interface{}{database}
	}

	rows, err := engine.DB().Query(sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	cols = make(map[string]map[string]*core.Column)
	for rows.Next() {
		col := new(core.Column)
		col.Indexes = make(map[string]int)

		var db, table, columnName, isNullable, colType, colKey, extra, comment string
		var colDefault *string
		err = rows.Scan(&db, &table, &columnName, &isNullable, &colDefault, &colType, &colKey, &extra, &comment)
		if err != nil {
			return nil, err
		}
		col.Name = strings.Trim(columnName, "` ")
		col.Comment = comment
		if "YES" == isNullable {
			col.Nullable = true
		}

		if colDefault != nil {
			col.Default = *colDefault
			if col.Default == "" {
				col.DefaultIsEmpty = true
			}
		}

		cts := strings.Split(colType, "(")
		colType = strings.ToUpper(cts[0])
		var len1, len2 int
		if len(cts) == 2 {
			idx := strings.Index(cts[1], ")")
			if colType == core.Enum && cts[1][0] == '\'' { //enum
				options := strings.Split(cts[1][0:idx], ",")
				col.EnumOptions = make(map[string]int)
				for k, v := range options {
					v = strings.TrimSpace(v)
					v = strings.Trim(v, "'")
					col.EnumOptions[v] = k
				}
			} else if colType == core.Set && cts[1][0] == '\'' {
				options := strings.Split(cts[1][0:idx], ",")
				col.SetOptions = make(map[string]int)
				for k, v := range options {
					v = strings.TrimSpace(v)
					v = strings.Trim(v, "'")
					col.SetOptions[v] = k
				}
			} else {
				lens := strings.Split(cts[1][0:idx], ",")
				len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
				if err != nil {
					return nil, err
				}
				if len(lens) == 2 {
					len2, err = strconv.Atoi(lens[1])
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if colType == "FLOAT UNSIGNED" {
			colType = "FLOAT"
		}
		if colType == "DOUBLE UNSIGNED" {
			colType = "DOUBLE"
		}
		col.Length = len1
		col.Length2 = len2
		if _, ok := core.SqlTypes[colType]; ok {
			col.SQLType = core.SQLType{Name: colType, DefaultLength: len1, DefaultLength2: len2}
		} else {
			return nil, fmt.Errorf("Unknown colType %v", colType)
		}

		if colKey == "PRI" {
			col.IsPrimaryKey = true
		}
		if colKey == "UNI" {
			//col.is
		}

		if extra == "auto_increment" {
			col.IsAutoIncrement = true
		}

		if col.SQLType.IsText() || col.SQLType.IsTime() {
			if col.Default != "" {
				col.Default = "'" + col.Default + "'"
			} else {
				if col.DefaultIsEmpty {
					col.Default = "''"
				}
			}
		}
		k := fmt.Sprintf("%s.%s", db, table)
		if _, ok := cols[k]; ok {
			cols[k][col.Name] = col
		} else {
			cols[k] = map[string]*core.Column{col.Name: col}
		}
	}

	return
}

func (m *Cluster) indexes(engine *xorm.Engine, database string) (indexes map[string]map[string]*core.Index, err error) {
	args := []interface{}{}
	sql := ""
	if database == "*" {
		sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` " +
			"FROM `INFORMATION_SCHEMA`.`STATISTICS` " +
			"WHERE 1 = ? AND `TABLE_SCHEMA` NOT IN ('mysql', 'sys', 'information_schema') " +
			"ORDER BY `TABLE_SCHEMA`, `TABLE_NAME`;"
		args = []interface{}{1}
	} else {
		sql = "SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` " +
			"FROM `INFORMATION_SCHEMA`.`STATISTICS` " +
			"WHERE `TABLE_SCHEMA` = ? AND `TABLE_SCHEMA` NOT IN ('mysql', 'sys', 'information_schema') " +
			"ORDER BY `TABLE_SCHEMA`, `TABLE_NAME`;"
		args = []interface{}{database}
	}

	rows, err := engine.DB().Query(sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	indexes = make(map[string]map[string]*core.Index, 0)
	for rows.Next() {
		var indexType int
		var db, table, indexName, colName, nonUnique string
		err = rows.Scan(&db, &table, &indexName, &nonUnique, &colName)
		if err != nil {
			return nil, err
		}

		if indexName == "PRIMARY" {
			continue
		}

		if "YES" == nonUnique || nonUnique == "1" {
			indexType = core.IndexType
		} else {
			indexType = core.UniqueType
		}

		colName = strings.Trim(colName, "` ")

		var index *core.Index
		var ok bool
		if index, ok = indexes[fmt.Sprintf("%s.%s", db, table)][indexName]; !ok {
			index = new(core.Index)
			index.Type = indexType
			index.Name = indexName
			k := fmt.Sprintf("%s.%s", db, table)
			if indexes[k] == nil {
				indexes[k] = map[string]*core.Index{indexName: index}
			} else {
				indexes[k][indexName] = index
			}
		}
		index.AddColumn(colName)
	}

	return
}

// Stat 仿照os.Stat方法，用于获得一个数据库的信息，数据库不存在则报错
func (m *Cluster) Stat(name string, passwd func(c *Cluster) []byte) (database *Database, err error) {
	databases := []Database{}
	databases, err = m.Databases(passwd)
	if err != nil {
		return
	}

	for _, elem := range databases {
		if elem.Name == name {
			database = &elem
			break
		}
	}

	if database == nil {
		err = fmt.Errorf("Unknown database '%s'", name)
	}

	return
}

// Repack 因为XORM默认没有暴露表结构体上的列信息，所以这里重新封包，方便序列化
func (m *Cluster) Repack(tables []*core.Table) (L []*Table) {
	for _, t := range tables {
		L = append(L, &Table{
			Name:          t.Name,
			Columns:       t.Columns(),
			Indexes:       t.Indexes,
			PrimaryKeys:   t.PrimaryKeys,
			AutoIncrement: t.AutoIncrement,
			Created:       t.Created,
			Updated:       t.Updated,
			Deleted:       t.Deleted,
			Version:       t.Version,
			StoreEngine:   t.StoreEngine,
			Charset:       t.Charset,
			Comment:       t.Comment,
		})
	}
	return
}

// Database 数据库信息的结构体
type Database struct {
	Name    string
	Charset string
	Collate string
}

// Table 重新封包向外直接暴露Columns
type Table struct {
	Name          string
	Columns       []*core.Column
	Indexes       map[string]*core.Index
	PrimaryKeys   []string
	AutoIncrement string
	Created       map[string]bool
	Updated       string
	Deleted       string
	Version       string
	StoreEngine   string
	Charset       string
	Collate       string
	Comment       string
}
