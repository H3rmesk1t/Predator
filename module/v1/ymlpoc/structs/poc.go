package structs

import "gopkg.in/yaml.v2"

type SetMapSlice = yaml.MapSlice
type PayloadsMapSlice = yaml.MapSlice

type Payloads struct {
	Continue bool             `yaml:"continue,omitempty"`
	Payloads PayloadsMapSlice `yaml:"payloads"`
}

type Detail struct {
	Author      string   `yaml:"author"`
	Links       []string `yaml:"links"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Tags        string   `yaml:"tags"`
}

type Poc struct {
	Name       string       `yaml:"name"`
	Transport  string       `yaml:"transport"`
	Set        SetMapSlice  `yaml:"set"`
	Payloads   Payloads     `yaml:"payloads"`
	Rules      RuleMapSlice `yaml:"rules"`
	Expression string       `yaml:"expression"`
	Detail     Detail       `yaml:"detail"`
}

type RuleRequest struct {
	Cache           bool              `yaml:"cache"`
	Raw             string            `yaml:"raw"`
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers"`
	Body            string            `yaml:"body"`
	FollowRedirects bool              `yaml:"follow_redirects"`
	Content         string            `yaml:"content"`
	ReadTimeout     string            `yaml:"read_timeout"`
	ConnectionID    string            `yaml:"connection_id"`
}

type Rule struct {
	Request    RuleRequest   `yaml:"request"`
	Expression string        `yaml:"expression"`
	Output     yaml.MapSlice `yaml:"output"`
	Order      int
}

type RuleAlias struct {
	Request    RuleRequest   `yaml:"request"`
	Expression string        `yaml:"expression"`
	Output     yaml.MapSlice `yaml:"output"`
}

type RuleMapItem struct {
	Key   string
	Value Rule
}

type RuleMapSlice []RuleMapItem

var ORDER = 0

func (r *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ruleAlias RuleAlias
	if err := unmarshal(&ruleAlias); err != nil {
		return err
	}

	r.Request = ruleAlias.Request
	r.Expression = ruleAlias.Expression
	r.Output = ruleAlias.Output
	r.Order = ORDER
	ORDER += 1

	return nil
}

func (rm *RuleMapSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ORDER = 0
	tempMap := make(map[string]Rule, 1)
	err := unmarshal(&tempMap)
	if err != nil {
		return err
	}

	newRuleSlice := make([]RuleMapItem, len(tempMap))
	for roleName, role := range tempMap {
		if role.Order < len(tempMap) {
			newRuleSlice[role.Order] = RuleMapItem{
				Key:   roleName,
				Value: role,
			}
		}
	}

	*rm = RuleMapSlice(newRuleSlice)

	return nil
}
