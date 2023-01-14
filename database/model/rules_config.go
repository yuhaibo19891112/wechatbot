package model

const RulesConfigName = "rules_confg"

type RulesConfig struct {
	ID           uint
	RuleType     int
	SendTimeType int
	SendTime     string
	SendUser     string
	SendGroup    string
}
