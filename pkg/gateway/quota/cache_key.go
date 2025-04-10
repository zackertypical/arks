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
	"bytes"
	"sync"
)

type CacheKeyGenerator struct {
	Prefix string
	pool   sync.Pool // Use buffer pool to improve performance
}

func NewCacheKeyGenerator(prefix string) *CacheKeyGenerator {
	return &CacheKeyGenerator{
		Prefix: prefix,
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (g *CacheKeyGenerator) Generate(desc *QuotaRequest) string {
	b := g.pool.Get().(*bytes.Buffer)
	defer g.pool.Put(b)
	b.Reset()

	// Write prefix
	b.WriteString(g.Prefix)
	b.WriteByte(':')

	// Write identifier
	for _, kv := range desc.Identifier {
		b.WriteString(kv.Key)
		b.WriteByte('=')
		b.WriteString(kv.Value)
		b.WriteByte('.')
	}

	return b.String()
}
