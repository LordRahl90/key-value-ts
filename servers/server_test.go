package servers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"key-value-ts/domains/entities"
	"key-value-ts/domains/storage/files"
	"key-value-ts/requests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	path   = "testdata/test/"
	server *Server

	errNoMockInitialized = errors.New("no mock setup for function")
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	s, err := files.New(path)
	if err != nil {
		panic(err)
	}
	server = New(s)
	code = m.Run()
}

func TestPutSequence(t *testing.T) {
	rec := httptest.NewRecorder()
	req := requests.Request{
		Key:       "my-key",
		Timestamp: 101,
		Value:     "My Value",
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	require.NoError(t, err)
	server.router.ServeHTTP(rec, request)

	require.NotNil(t, rec)
	assert.Equal(t, http.StatusCreated, rec.Code)
	exp := `{"message":"sequence saved successfully"}`
	assert.Equal(t, exp, rec.Body.String())
}

func TestPutSequence_WithInvalidJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := []byte(`{{"key":"my-key","value":"My Value","timestamp":101},"`)
	request, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(req))
	require.NoError(t, err)
	server.router.ServeHTTP(rec, request)

	require.NotNil(t, rec)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	exp := `{"error":"invalid character '{' looking for beginning of object key string"}`
	assert.Equal(t, exp, rec.Body.String())
}

func TestPutSequence_WithError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := requests.Request{
		Key:       "my-key",
		Timestamp: 101,
		Value:     "My Value",
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	require.NoError(t, err)

	storer := mockStorer{
		SaveFunc: func(e entities.Sequence) error {
			return fmt.Errorf("cannot store record")
		},
	}

	svr := New(&storer)
	svr.router.ServeHTTP(rec, request)

	require.NotNil(t, rec)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	exp := `{"error":"cannot store record"}`
	assert.Equal(t, exp, rec.Body.String())
}

func TestPutSequence_NoMockInitialized(t *testing.T) {
	rec := httptest.NewRecorder()
	req := requests.Request{
		Key:       "my-key",
		Timestamp: 101,
		Value:     "My Value",
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	require.NoError(t, err)

	storer := mockStorer{}

	svr := New(&storer)
	svr.router.ServeHTTP(rec, request)

	require.NotNil(t, rec)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	exp := `{"error":"no mock setup for function"}`
	assert.Equal(t, exp, rec.Body.String())
}

func TestGetSequence(t *testing.T) {
	rec := httptest.NewRecorder()
	req := requests.Request{
		Key:       "my-key",
		Timestamp: 101,
		Value:     "My Value",
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	require.NoError(t, err)
	server.router.ServeHTTP(rec, request)

	require.NotNil(t, rec)
	assert.Equal(t, http.StatusCreated, rec.Code)
	exp := `{"message":"sequence saved successfully"}`
	assert.Equal(t, exp, rec.Body.String())

	path := fmt.Sprintf("/?key=%s&timestamp=%d", "my-key", 101)

	request, err = http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	getRec := httptest.NewRecorder()

	server.router.ServeHTTP(getRec, request)
	require.NotNil(t, rec)

	require.Equal(t, http.StatusOK, getRec.Code)
	exp = `{"value":"My Value"}`
	assert.Equal(t, exp, getRec.Body.String())
}

func TestGetSequence_NonExistent(t *testing.T) {
	rec := httptest.NewRecorder()
	path := fmt.Sprintf("/?key=%s&timestamp=%d", "my-key", 110)

	request, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	getRec := httptest.NewRecorder()

	server.router.ServeHTTP(getRec, request)
	require.NotNil(t, rec)

	require.Equal(t, http.StatusOK, getRec.Code)
	exp := `{"value":""}`
	assert.Equal(t, exp, getRec.Body.String())
}

func TestGetSequence_InvalidTimestamp(t *testing.T) {
	rec := httptest.NewRecorder()

	path := fmt.Sprintf("/?key=%s&timestamp=timestamp", "my-key")
	request, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	getRec := httptest.NewRecorder()

	server.router.ServeHTTP(getRec, request)
	require.NotNil(t, rec)

	require.Equal(t, http.StatusBadRequest, getRec.Code)
	exp := `{"error":"strconv.Atoi: parsing \"timestamp\": invalid syntax"}`
	assert.Equal(t, exp, getRec.Body.String())
}

func TestGetSequence_WithError(t *testing.T) {
	rec := httptest.NewRecorder()

	path := fmt.Sprintf("/?key=%s&timestamp=%d", "my-key", 101)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	getRec := httptest.NewRecorder()

	storer := &mockStorer{
		GetFunc: func(key string, timestamp int64) (string, error) {
			return "", fmt.Errorf("cannot get value for some reason")
		},
	}
	svr := New(storer)

	svr.router.ServeHTTP(getRec, request)
	require.NotNil(t, rec)

	require.Equal(t, http.StatusInternalServerError, getRec.Code)
	exp := `{"error":"cannot get value for some reason"}`
	assert.Equal(t, exp, getRec.Body.String())
}

func TestGetSequence_NoMockInitialized(t *testing.T) {
	rec := httptest.NewRecorder()

	path := fmt.Sprintf("/?key=%s&timestamp=%d", "my-key", 101)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	getRec := httptest.NewRecorder()

	storer := &mockStorer{}
	svr := New(storer)

	svr.router.ServeHTTP(getRec, request)
	require.NotNil(t, rec)

	require.Equal(t, http.StatusInternalServerError, getRec.Code)
	exp := `{"error":"no mock setup for function"}`
	assert.Equal(t, exp, getRec.Body.String())
}

type mockStorer struct {
	SaveFunc func(e entities.Sequence) error
	GetFunc  func(key string, timestamp int64) (string, error)
}

// Save mock save function
func (ms *mockStorer) Save(e entities.Sequence) error {
	if ms.SaveFunc == nil {
		return errNoMockInitialized
	}
	return ms.SaveFunc(e)
}

// Get mock Get function
func (ms *mockStorer) Get(key string, ts int64) (string, error) {
	if ms.GetFunc == nil {
		return "", errNoMockInitialized
	}
	return ms.GetFunc(key, ts)
}
