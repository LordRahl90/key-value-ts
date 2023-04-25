package files

import (
	"log"
	"os"
	"testing"
	"time"

	"key-value-ts/domains/entities"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestSave(t *testing.T) {
	cs, err := New("testdata/test")
	require.NoError(t, err)
	require.NotNil(t, cs)

	e := entities.Sequence{
		Key:       "my-key",
		Timestamp: int64(time.Now().Unix()),
		Value:     "hello-world",
	}

	err = cs.Save(e)
	require.NoError(t, err)
}

func TestSaveMultipleDir(t *testing.T) {
	path := "testdata/test/kvstore/"
	cs, err := New(path)
	require.NoError(t, err)
	require.NotNil(t, cs)

	e := entities.Sequence{
		Key:       "my-key",
		Timestamp: time.Now().Unix(),
		Value:     "hello-world",
	}

	err = cs.Save(e)
	require.NoError(t, err)

	_, err = os.Stat(path + "my-key.csv")
	require.NoError(t, err)
	require.NoError(t, os.RemoveAll(path))
}

func TestSave_AssertAppending(t *testing.T) {
	path := "testdata/test/appends/"
	cs, err := New(path)
	require.NoError(t, err)
	require.NotNil(t, cs)

	for i := 0; i <= 10; i++ {
		et := entities.Sequence{
			Key:       "append_key",
			Timestamp: time.Now().Unix(),
			Value:     gofakeit.Name(),
		}
		require.NoError(t, cs.Save(et))
	}

	_, err = os.Stat(path + "append_key.csv")
	require.NoError(t, err)

	require.NoError(t, os.RemoveAll(path))
}

func TestGetValue_NonExistingKey(t *testing.T) {
	path := "./testdata/demo"
	key := "non-existing-key"
	ts := int64(1682201632)
	cs, err := New(path)

	require.NoError(t, err)
	require.NotNil(t, cs)

	val, err := cs.Get(key, ts)
	require.NotNil(t, err)
	assert.Empty(t, val)
}

func TestGetValue(t *testing.T) {
	path := "./testdata/demo"
	key := "read"
	ts := int64(427852973007343884)
	cs, err := New(path)

	require.NoError(t, err)
	require.NotNil(t, cs)

	val, err := cs.Get(key, ts)
	require.NoError(t, err)
	assert.Equal(t, "131-259-5334", val)
}

func BenchmarkSave(b *testing.B) {
	path := "testdata/bench/"
	cs, err := New(path)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		e := entities.Sequence{
			Key:       "my-key",
			Timestamp: time.Now().Unix() + int64(gofakeit.Uint64()),
			Value:     gofakeit.PhoneFormatted(),
		}
		if err := cs.Save(e); err != nil {
			panic(err)
		}
	}
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}

func BenchmarkGet(b *testing.B) {
	table := []struct {
		key, name string
		timestamp int64
	}{
		{
			name:      "farthest-key",
			key:       "read",
			timestamp: 6152146711361107039,
		},
		{
			name:      "non-existent-key",
			key:       "read",
			timestamp: 10001,
		},
		{
			name:      "last-key",
			key:       "read",
			timestamp: 560053162071168685,
		},
		{
			name:      "early-key",
			key:       "read",
			timestamp: 5165520820317048201,
		},
	}

	path := "testdata/demo/"
	cs, err := New(path)
	if err != nil {
		panic(err)
	}

	for _, v := range table {
		b.Run(v.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := cs.Get(v.key, v.timestamp)
				if err != nil {
					log.Fatal(err)
				}
			}
		})
	}

}
