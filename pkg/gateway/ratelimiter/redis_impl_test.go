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
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedis() redis.UniversalClient {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return client
}

func TestRedis_IncrUsage(t *testing.T) {
	client := setupRedis()

	service := NewRedisRateLimter(client, "test")
	ctx := context.Background()

	tests := []struct {
		name     string
		requests []*RateLimitRequest
		wantErr  bool
	}{
		{
			name: "single request",
			requests: []*RateLimitRequest{
				{
					Identifier: []Label{
						{
							Key:   "user",
							Value: "123",
						},
						{
							Key:   "model",
							Value: "m1",
						},
					},
					RuleName: "rpd",
					Request:  10,
					Limit:    100,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple requests",
			requests: []*RateLimitRequest{
				{
					Identifier: []Label{
						{
							Key:   "user",
							Value: "123",
						},
						{
							Key:   "model",
							Value: "m1",
						},
					},
					RuleName: "rpm",
					Request:  10,
					Limit:    100,
				},
				{
					Identifier: []Label{
						{
							Key:   "user",
							Value: "123",
						},
						{
							Key:   "model",
							Value: "m1",
						},
					},
					RuleName: "tpm",
					Request:  20,
					Limit:    200,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DoLimit(ctx, tt.requests)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// TODO: check values
			// for _, req := range tt.requests {

			// }
		})
	}
}
