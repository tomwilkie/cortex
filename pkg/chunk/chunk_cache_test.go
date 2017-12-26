package chunk

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
	promchunk "github.com/weaveworks/cortex/pkg/prom1/storage/local/chunk"
	"golang.org/x/net/context"
)

type mockMemcache struct {
	sync.RWMutex
	contents map[string][]byte
}

func newMockMemcache() *mockMemcache {
	return &mockMemcache{
		contents: map[string][]byte{},
	}
}

func (m *mockMemcache) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	m.RLock()
	defer m.RUnlock()
	result := map[string]*memcache.Item{}
	for _, k := range keys {
		if c, ok := m.contents[k]; ok {
			result[k] = &memcache.Item{
				Value: c,
			}
		}
	}
	return result, nil
}

func (m *mockMemcache) Set(item *memcache.Item) error {
	m.Lock()
	defer m.Unlock()
	m.contents[item.Key] = item.Value
	return nil
}

func TestChunkCache(t *testing.T) {
	c := Cache{
		memcache: newMockMemcache(),
	}

	const (
		chunkLen = 13 * 3600 // in seconds
	)

	// put 100 chunks from 0 to 99
	keys := []string{}
	chunks := []Chunk{}
	for i := 0; i < 100; i++ {
		ts := model.TimeFromUnix(int64(i * chunkLen))
		promChunk, _ := promchunk.New().Add(model.SamplePair{
			Timestamp: ts,
			Value:     model.SampleValue(i),
		})
		chunk := NewChunk(
			Descriptor{
				UserID:  userID,
				From:    ts,
				Through: ts.Add(chunkLen),
				Metric: model.Metric{
					model.MetricNameLabel: "foo",
					"bar": "baz",
				},
				Fingerprint: model.Fingerprint(1),
				Encoding:    promChunk[0].Encoding(),
			},
			promChunk[0],
		)

		buf, err := chunk.Encode()
		require.NoError(t, err)

		key := chunk.Descriptor().ExternalKey()
		err = c.StoreChunk(context.Background(), key, buf)
		require.NoError(t, err)

		keys = append(keys, key)
		chunks = append(chunks, chunk)
	}

	for i := 0; i < 100; i++ {
		index := rand.Intn(len(keys))
		key := keys[index]

		desc, err := parseExternalKey(userID, key)
		require.NoError(t, err)

		found, missing, err := c.FetchChunkData(context.Background(), []Descriptor{desc})
		require.NoError(t, err)
		require.Empty(t, missing)
		require.Len(t, found, 1)
		require.Equal(t, chunks[index], found[0])
	}

	// test getting them all
	descs := []Descriptor{}
	for i := 0; i < len(keys); i++ {
		desc, err := parseExternalKey(userID, keys[i])
		require.NoError(t, err)
		descs = append(descs, desc)
	}
	found, missing, err := c.FetchChunkData(context.Background(), descs)
	require.NoError(t, err)
	require.Empty(t, missing)
	require.Len(t, found, len(keys))
	require.Equal(t, chunks, found)
}
