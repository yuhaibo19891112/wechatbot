package services

import (
	"github.com/869413421/wechatbot/database"
	"github.com/869413421/wechatbot/database/model"
	"gorm.io/gorm"
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

func (s *RulesConfigService) GetNewsRulesConfig() *model.RulesConfig {
	var config model.RulesConfig
	affected := s.getDB().Table(model.RulesConfigName).Where("rule_type=1").First(&config).RowsAffected
	if affected > 0 {
		return &config
	}
	return nil
}
