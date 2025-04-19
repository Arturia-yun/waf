package models

import (
	"errors"
	"time"
)

// Rule 表示WAF规则
type Rule struct {
	Id        int64     `xorm:"pk autoincr" json:"id"`
	RuleType  string    `xorm:"varchar(50) notnull" json:"rule_type"`
	RuleItem  string    `xorm:"varchar(500) notnull" json:"rule_item"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

// TableName 返回表名
func (r *Rule) TableName() string {
	return "waf_rule"
}

// GetRules 获取所有规则
func GetRules() ([]*Rule, error) {
	rules := make([]*Rule, 0)
	err := Engine.Find(&rules)
	return rules, err
}

// GetRulesByType 根据类型获取规则
func GetRulesByType(ruleType string) ([]*Rule, error) {
	rules := make([]*Rule, 0)
	err := Engine.Where("rule_type = ?", ruleType).Find(&rules)
	return rules, err
}

// GetRuleById 根据ID获取规则
func GetRuleById(id int64) (*Rule, error) {
	rule := &Rule{Id: id}
	has, err := Engine.Get(rule)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("规则不存在")
	}
	return rule, nil
}

// CreateRule 创建新规则
func CreateRule(rule *Rule) error {
	_, err := Engine.Insert(rule)
	return err
}

// UpdateRule 更新规则
func UpdateRule(rule *Rule) error {
	_, err := Engine.ID(rule.Id).Update(rule)
	return err
}

// DeleteRule 删除规则
func DeleteRule(id int64) error {
	_, err := Engine.ID(id).Delete(&Rule{})
	return err
}

// GetAllRuleTypes 获取所有规则类型
func GetAllRuleTypes() []string {
	return []string{
		"url",
		"args",
		"blackip",
		"whiteip",
		"whiteUrl",
		"useragent",
		"cookie",
		"post",
		"header",
	}
}