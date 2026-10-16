package models

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

// Cluster 群集的模型
type Cluster struct {
	ClusterID   uint   `db:"cluster_id"                      valid:"-"                                json:"cluster_id"  gqlgen:"-"`                         //
	UUID        string `db:"uuid"                  valid:"-"                                json:"uuid"        gqlgen:"UUID"`                                //
	Host        string `db:"host"              valid:"required,length(1|100),alphanum"  json:"host"        gqlgen:"Host"`                                    //
	IP          string `db:"ip"                 valid:"required,int,range(0|4294967295)" json:"ip"          gqlgen:"Ip"`                                     //
	Port        uint16 `db:"port" valid:"required,port"                    json:"port"        gqlgen:"Port"`                                                 //
	Alias       string `db:"alias"              valid:"required,runelength(1|100)"       json:"alias"       gqlgen:"Alias"`                                  //
	User        string `db:"user"                                valid:"required,length(1|50),alphanum"   json:"user"        gqlgen:"User"`                  //
	Type        string `db:"type"                                 valid:"required,in(mysql|postgres)"      json:"type"        gqlgen:"Type"`                 // 目标库类型
	Password    []byte `db:"password"                          valid:"required"                         json:"-"           gqlgen:"-"`                       // 双向加密
	FingerPrint []byte `db:"fingerprint"                       valid:"required"                         json:"-"           gqlgen:"-"`                       //
	Status      uint8  `db:"status"                                  valid:"required,matches(^[0-9]$)"        json:"status"      gqlgen:"Status"`            //
	Version     int    `db:"version"                                                 valid:"-"                                json:"version"     gqlgen:"-"` //
	UpdateAt    uint   `db:"update_at"                                   valid:"-"                                json:"update_at"   gqlgen:"UpdateAt"`      //
	CreateAt    uint   `db:"create_at"                                   valid:"-"                                json:"create_at"   gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Cluster) TableName() string {
	return "mm_clusters"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Cluster) BeforeInsert() {
	// GraphQL 尚未暴露 type，缺省按 MySQL 记录，避免写入空串。
	if m.Type == "" {
		m.Type = string(domain.ClusterTypeMySQL)
	}
	m.UUID = uuid.New().String()
	m.Host = strings.TrimSpace(m.Host)
	m.IP = strings.TrimSpace(m.IP)
	m.User = strings.TrimSpace(m.User)
	m.FingerPrint = m.fingerprintHash(fmt.Sprintf("%s-%s-%d-%s-%s",
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
	fp := m.fingerprintHash(fmt.Sprintf("%s-%s-%d-%s-%s",
		m.Host,
		m.IP,
		m.Port,
		m.User,
		m.Password))
	if !bytes.Equal(fp, m.FingerPrint) {
		m.FingerPrint = fp
	}
	m.UpdateAt = uint(time.Now().Unix())
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

// fingerprintHash 对群集连接信息做指纹散列（SHA-256）。
func (m *Cluster) fingerprintHash(data string) []byte {
	t := sha256.New()
	_, _ = io.WriteString(t, data)
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

// Database 数据库信息的结构体
type Database struct {
	Name    string
	Charset string
	Collate string
}
