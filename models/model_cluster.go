package models

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
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
func (m *Cluster) Metadata(name string, passwd func(c *Cluster) []byte) (tables []*core.Table, err error) {
L:
	for {
		var engine *xorm.Engine

		engine, err = m.Connect(name, passwd)
		if err != nil {
			break L
		}
		defer engine.Close()

		if _, err = m.Stat(name, passwd); err != nil {
			break L
		}

		tables, err = engine.DBMetas()
		if err != nil {
			break L
		}

		break
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
