/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ratelimiter

import (
	"encoding/json"
	"fmt"
	"time"
)

// pkg/gateway/limiter/types.go
type LimitUnit string

const (
	UnitSecond LimitUnit = "second"
	UnitMinute LimitUnit = "minute"
	UnitHour   LimitUnit = "hour"
	UnitDay    LimitUnit = "day"
)

func UnitToSeconds(unit LimitUnit) int64 {
	switch unit {
	case UnitSecond:
		return 1
	case UnitMinute:
		return 60
	case UnitHour:
		return 3600
	case UnitDay:
		return 86400
	}
	panic(fmt.Sprintf("unit %s not valid, should not get here", unit))
}

type LimitType string

const (
	TypeRequest LimitType = "request"
	TypeToken   LimitType = "token"
	// New types can be easily added in the future
)

// Rate limit rule configuration
type LimitRuleConfig struct {
	Name        string    `json:"name"`                  // Rule name, e.g., "rpm", "tpd"
	Type        LimitType `json:"type"`                  // Limit type
	Period      int64     `json:"period"`                // Time period count
	Unit        LimitUnit `json:"unit"`                  // Time unit
	Description string    `json:"description,omitempty"` // Rule description
}

// pkg/gateway/limiter/config.go
type LimitConfig struct {
	rules map[string]*LimitRuleConfig
}

func NewLimitConfig() *LimitConfig {
	return &LimitConfig{
		rules: make(map[string]*LimitRuleConfig),
	}
}

func (c *LimitConfig) AddRule(rule *LimitRuleConfig) {
	c.rules[rule.Name] = rule
}

func (c *LimitConfig) GetRule(name string) *LimitRuleConfig {
	return c.rules[name]
}

// Key-value pair for resource identification
type Label struct {
	Key   string
	Value string
}

// RateLimit description
type RateLimitRequest struct {
	Identifier []Label
	RuleName   string // Use rule name
	Limit      int64  // Limit value
	Request    int64  // Request value
}

// RateLimit response
type RateLimitResponse struct {
	RuleName     string    `json:"ruleName"`
	OverLimit    bool      `json:"overLimit"`
	CurrentUsage int64     `json:"currentUsage"`
	LimitMax     int64     `json:"limitMax"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (r *RateLimitResponse) JSON() string {
	json, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(json)
}
