package services

import (
	"github.com/869413421/wechatbot/database"
	"github.com/869413421/wechatbot/database/model"
	"gorm.io/gorm"
	"log"
)

type RulesConfigService struct {
}

func NewRulesConfigService() *RulesConfigService {
	return &RulesConfigService{}
}

// 获取数据库连接
func (s *RulesConfigService) getDB() *gorm.DB {
	return database.GetDB()
}

func (s *RulesConfigService) GetRulesConfig(ruleType string) *model.RulesConfig {
	var config model.RulesConfig
	affected := s.getDB().Table(model.RulesConfigName).Where("rule_type=" + ruleType).First(&config).RowsAffected
	if affected > 0 {
		return &config
	}
	return nil
}

func (s *RulesConfigService) UpdateRulesConfig(ruleType string, timeCron string, user string, group string, content string) {
	if timeCron == "" && user == "" && group == "" && content == ""{
		return
	}
	log.Printf("update: %s, %s, %s, %s", timeCron, user, group, content)
	updateCfg := model.RulesConfig{}
	if timeCron != "" {
		updateCfg.SendTime = timeCron
	}
	if user != "" {
		updateCfg.SendUser = user
	}
	if group != "" {
		updateCfg.SendGroup = group
	}
	if content != "" {
		updateCfg.Content = content
	}
	s.getDB().Table(model.RulesConfigName).Where("rule_type=" + ruleType).Updates(&updateCfg)
}
