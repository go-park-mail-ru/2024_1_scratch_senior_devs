package paging

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestApply(t *testing.T) {
	params := &pagingParams{
		count:  10,
		offset: 0,
	}

	opt1 := &funcPagingParams{
		f: func(params *pagingParams) {
			params.count = 20
		},
	}

	opt2 := &funcPagingParams{
		f: func(params *pagingParams) {
			params.offset = 5
		},
	}

	options := []Option{opt1, opt2}
	for _, opt := range options {
		opt.apply(params)
	}

	assert.Equal(t, 20, params.count)
	assert.Equal(t, 5, params.offset)
}

func TestWithCustomCount(t *testing.T) {
	params := &pagingParams{
		count:  10,
		offset: 0,
	}

	opt := WithCustomCount(30)
	opt.apply(params)

	assert.Equal(t, 30, params.count)
	assert.Equal(t, 0, params.offset)
}

func TestWithCustomOffset(t *testing.T) {
	params := &pagingParams{
		count:  10,
		offset: 0,
	}

	opt := WithCustomOffset(15)
	opt.apply(params)

	assert.Equal(t, 10, params.count)
	assert.Equal(t, 15, params.offset)
}

func TestGetParams(t *testing.T) {
	tests := []struct {
		name   string
		r      *http.Request
		count  int
		offset int
		isErr  bool
	}{
		{
			name:   "Test_GetParams_Success",
			r:      &http.Request{URL: &url.URL{RawQuery: "count=10&offset=5"}},
			count:  10,
			offset: 5,
			isErr:  false,
		},
		{
			name:   "Test_GetParams_Fail_1",
			r:      &http.Request{URL: &url.URL{RawQuery: "count=invalid&offset=5"}},
			count:  0,
			offset: 0,
			isErr:  true,
		},
		{
			name:   "Test_GetParams_Fail_2",
			r:      &http.Request{URL: &url.URL{RawQuery: "count=10&offset=invalid"}},
			count:  0,
			offset: 0,
			isErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, offset, err := GetParams(tt.r)

			assert.Equal(t, tt.count, count)
			assert.Equal(t, tt.offset, offset)

			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
