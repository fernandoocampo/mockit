package webservers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fernandoocampo/webclient"
	"github.com/stretchr/testify/assert"
)

func TestMockingWebServerWithHandlerObject(t *testing.T) {
	cases := map[string]struct {
		path          string
		timeoutMillis int
		isError       error
		handler       *handlerMock
		want          *webclient.Response
	}{
		"success_with_data": {
			path:    "/people",
			handler: newHandlerMock(200, `[{name: "vane", age:  43}]`),
			want: &webclient.Response{
				StatusCode: 200,
				Data:       []byte(`[{name: "vane", age:  43}]`),
			},
		},
		"success_without_data": {
			path:    "/people",
			handler: newHandlerMock(200, ""),
			want: &webclient.Response{
				StatusCode: 200,
				Data:       []byte(""),
			},
		},
		"request_with_timeout": {
			path:          "/people",
			timeoutMillis: 200,
			handler:       newHandlerMockWithTimeout(100, 200, ""),
			isError:       errors.New("context deadline exceeded"),
		},
	}

	for name, data := range cases {
		name, data := name, data
		t.Run(name, func(st *testing.T) {
			st.Parallel()

			// GIVEN
			// here we create the mocked web server
			mockServer := httptest.NewServer(data.handler)
			defer mockServer.Close()
			ctx := context.TODO()
			var cancel context.CancelFunc
			if data.timeoutMillis != 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(data.timeoutMillis)*time.Millisecond)
				defer cancel()
			}
			// WHEN
			got, err := webclient.New(mockServer.URL).Get(ctx, data.path)
			// THEN
			assertError(st, err)
			assert.Equal(st, data.want, got)
		})
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), err.Error())

		return
	}

	assert.NoError(t, err)
}

type handlerMock struct {
	resp    []byte
	code    int
	sleepms int
}

func newHandlerMock(code int, resp string) *handlerMock {
	var respdata []byte
	if resp != "" {
		respdata = []byte(resp)
	}
	return &handlerMock{
		code: code,
		resp: respdata,
	}
}

func newHandlerMockWithTimeout(code, sleepms int, resp string) *handlerMock {
	var respdata []byte
	if resp != "" {
		respdata = []byte(resp)
	}
	return &handlerMock{
		code:    code,
		resp:    respdata,
		sleepms: sleepms,
	}
}

func (w *handlerMock) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if w.sleepms != 0 {
		time.Sleep(time.Duration(w.sleepms) * time.Millisecond)
	}
	rw.WriteHeader(w.code)
	rw.Write(w.resp)
}
