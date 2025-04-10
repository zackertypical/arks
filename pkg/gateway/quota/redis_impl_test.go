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

package quota

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

func TestRedisQuotaService_IncrUsage(t *testing.T) {
	client := setupRedis()

	service := NewRedisQuotaService(client, "test")
	ctx := context.Background()

	tests := []struct {
		name     string
		requests []*QuotaRequest
		wantErr  bool
	}{
		{
			name: "single request",
			requests: []*QuotaRequest{
				{
					Identifier: []Label{
						{
							Key:   "namespace",
							Value: "arks",
						},
						{
							Key:   "quotaname",
							Value: "q1",
						},
						{
							Key:   "type", // prompt, response, total etc.
							Value: "prompt",
						},
					},
					Request: 10,
					Limit:   100,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple requests",
			requests: []*QuotaRequest{
				{
					Identifier: []Label{
						{
							Key:   "namespace",
							Value: "arks",
						},
						{
							Key:   "quotaname",
							Value: "q2",
						},
						{
							Key:   "type", // prompt, response, total etc.
							Value: "prompt",
						},
					},
					Request: 10,
					Limit:   100,
				},
				{
					Identifier: []Label{
						{
							Key:   "namespace",
							Value: "arks",
						},
						{
							Key:   "quotaname",
							Value: "q2",
						},
						{
							Key:   "type", // prompt, response, total etc.
							Value: "response",
						},
					},
					Request: 20,
					Limit:   200,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.IncrUsage(ctx, tt.requests)
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

// func TestRedisQuotaService_GetUsage(t *testing.T) {
// 	client := setupRedis(t)
// 	defer mr.Close()

// 	service := NewRedisQuotaService(client, "test")
// 	ctx := context.Background()

//
// 	requests := []*QuotaRequest{
// 		{
// 			Identifier: "user1",
// 			Request:    10,
// 			Limit:      100,
// 		},
// 		{
// 			Identifier: "user2",
// 			Request:    90,
// 			Limit:      80,
// 		},
// 	}

//
// 	err := service.IncrUsage(ctx, requests)
// 	assert.NoError(t, err)

//
// 	results, err := service.GetUsage(ctx, requests)
// 	assert.NoError(t, err)
// 	assert.Len(t, results, 2)

//
// 	assert.Equal(t, int64(10), results[0].CurrentUsage)
// 	assert.False(t, results[0].OverLimit)
// 	assert.Equal(t, int64(90), results[1].CurrentUsage)
// 	assert.True(t, results[1].OverLimit)
// }

// func TestRedisQuotaService_SetUsage(t *testing.T) {
// 	client := setupRedis(t)
// 	defer mr.Close()

// 	service := NewRedisQuotaService(client, "test")
// 	ctx := context.Background()

// 	tests := []struct {
// 		name     string
// 		requests []*QuotaRequest
// 		wantErr  bool
// 	}{
// 		{
// 			name: "set single usage",
// 			requests: []*QuotaRequest{
// 				{
// 					Identifier: "user1",
// 					Request:    50,
// 					Limit:      100,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "set multiple usage",
// 			requests: []*QuotaRequest{
// 				{
// 					Identifier: "user1",
// 					Request:    60,
// 					Limit:      100,
// 				},
// 				{
// 					Identifier: "user2",
// 					Request:    30,
// 					Limit:      200,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := service.SetUsage(ctx, tt.requests)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
// 			assert.NoError(t, err)

//
// 			for _, req := range tt.requests {
// 				key := "test:" + req.Identifier
// 				val, err := client.Get(ctx, key).Int64()
// 				assert.NoError(t, err)
// 				assert.Equal(t, req.Request, val)
// 			}
// 		})
// 	}
// }

func TestRedisQuotaService_ConcurrentAccess(t *testing.T) {
	client := setupRedis()

	service := NewRedisQuotaService(client, "test")
	ctx := context.Background()
	keygen := NewCacheKeyGenerator("test")

	concurrency := 10
	done := make(chan bool)

	request := &QuotaRequest{
		Identifier: []Label{
			{
				Key:   "namespace",
				Value: "arks",
			},
			{
				Key:   "quotaname",
				Value: "q2",
			},
			{
				Key:   "type", // prompt, response, total etc.
				Value: "response",
			},
		},
		Request: 1,
		Limit:   100,
	}

	for i := 0; i < concurrency; i++ {
		go func() {
			err := service.IncrUsage(ctx, []*QuotaRequest{request})
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}

	key := keygen.Generate(request)
	val, err := client.Get(ctx, key).Int64()
	assert.NoError(t, err)
	assert.Equal(t, int64(concurrency), val)

	client.Del(ctx, key)
}
