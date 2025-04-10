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

import "context"

type RateLimterInterface interface {
	// CheckLimit Check rate limit when request comes
	CheckLimit(ctx context.Context, requests []*RateLimitRequest) ([]*RateLimitResponse, error)
	// DoLimit Execute rate limit when request returns
	DoLimit(ctx context.Context, requests []*RateLimitRequest) error
	// GetRuleByName Get rate limit rule
	GetRuleByName(name string) *LimitRuleConfig
}

// TODO: Example config, to be fetched from config center
func newLimitConfig() *LimitConfig {
	config := NewLimitConfig()

	// Add predefined rules
	config.AddRule(&LimitRuleConfig{
		Name:        "rpm",
		Type:        TypeRequest,
		Period:      1,
		Unit:        UnitMinute,
		Description: "Requests per minute",
	})

	config.AddRule(&LimitRuleConfig{
		Name:        "rpd",
		Type:        TypeRequest,
		Period:      1,
		Unit:        UnitDay,
		Description: "Requests per day",
	})

	config.AddRule(&LimitRuleConfig{
		Name:        "tpm",
		Type:        TypeToken,
		Period:      1,
		Unit:        UnitMinute,
		Description: "Tokens per minute",
	})

	config.AddRule(&LimitRuleConfig{
		Name:        "tpd",
		Type:        TypeToken,
		Period:      1,
		Unit:        UnitDay,
		Description: "Tokens per day",
	})

	return config
}
