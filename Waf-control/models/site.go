package models

import (
	"log"
	"time"
)

// Site 表示WAF站点配置
// debuglevel: debug, info, notice, warn, error, crit, alert, emerg
// ssl: on, off
type Site struct {
	Id          int64     `json:"id"`
	SiteName    string    `xorm:"unique" json:"site_name"`
	Port        int       `json:"port"`
	BackendAddr []string  `json:"backend_addr"`
	UnrealAddr  []string  `json:"unreal_addr"`
	Ssl         string    `xorm:"varchar(10) notnull default 'off'" json:"ssl"`
	DebugLevel  string    `xorm:"varchar(10) notnull default 'error'" json:"debug_level"`
	LastChange  time.Time `xorm:"updated" json:"last_change"`
	Version     int       `xorm:"version" json:"version"` // 乐观锁
}

// TableName 返回表名
func (s *Site) TableName() string {
	return "waf_site"
}

// ListSite 获取所有站点
func ListSite() (sites []Site, err error) {
	sites = make([]Site, 0)
	err = Engine.Find(&sites)
	log.Println(err, sites)
	return sites, err
}

// GetSiteById 根据ID获取站点
func GetSiteById(id int64) (*Site, error) {
	site := &Site{Id: id}
	has, err := Engine.Get(site)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return site, nil
}

// ListSiteById 根据ID获取站点列表
func ListSiteById(Id int64) (sites []Site, err error) {
	sites = make([]Site, 0)
	err = Engine.Id(Id).Find(&sites)
	log.Println(err, sites)
	return sites, err
}

// NewSite 创建新站点
func NewSite(siteName string, Port int, BackendAddr []string, UnrealAddr []string, SSL string, DebugLevel string) (err error) {
	if SSL == "" {
		SSL = "off"
	}
	if DebugLevel == "" {
		DebugLevel = "error"
	}

	_, err = Engine.Insert(&Site{SiteName: siteName, Port: Port, BackendAddr: BackendAddr, UnrealAddr: UnrealAddr, Ssl: SSL, DebugLevel: DebugLevel})
	return err
}

// UpdateSite 更新站点
func UpdateSite(Id int64, SiteName string, Port int, BackendAddr []string, UnrealAddr []string, SSL string, DebugLevel string) (err error) {
	if SSL == "" {
		SSL = "off"
	}
	if DebugLevel == "" {
		DebugLevel = "error"
	}

	site := new(Site)
	Engine.Id(Id).Get(site)
	site.SiteName = SiteName
	site.Port = Port
	site.BackendAddr = BackendAddr
	site.UnrealAddr = UnrealAddr
	site.Ssl = SSL
	site.DebugLevel = DebugLevel
	_, err = Engine.Id(Id).Update(site)
	return err
}

// DelSite 删除站点
func DelSite(id int64) (err error) {
	_, err = Engine.Delete(&Site{Id: id})
	return err
}