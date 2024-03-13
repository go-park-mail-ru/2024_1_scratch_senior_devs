package utils

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type pagingParams struct {
	count  int
	offset int
}

var defaultPagingParams = pagingParams{
	count:  10,
	offset: 0,
}

type Option interface {
	apply(params *pagingParams)
}

type funcPagingParams struct {
	f func(params *pagingParams)
}

func (f *funcPagingParams) apply(params *pagingParams) {
	f.f(params)
}

func newFuncPagingParams(f func(params *pagingParams)) *funcPagingParams {
	return &funcPagingParams{
		f: f,
	}
}

func WithCustomCount(count int) Option {
	return newFuncPagingParams(func(o *pagingParams) {
		o.count = count
	})
}

func WithCustomOffset(offset int) Option {
	return newFuncPagingParams(func(o *pagingParams) {
		o.offset = offset
	})
}

func GetParams(r *http.Request, params ...Option) (int, int, error) {
	defaultParams := defaultPagingParams
	for _, param := range params {
		param.apply(&defaultParams)
	}

	strCount := r.URL.Query().Get("count")
	if strCount != "" {
		count, err := strconv.Atoi(strCount)
		if err != nil {
			return 0, 0, errors.Wrap(err, "invalid count param")
		}
		if count > 0 {
			defaultParams.count = count
		}
	}

	strOffset := r.URL.Query().Get("offset")
	if strOffset != "" {
		offset, err := strconv.Atoi(strOffset)
		if err != nil {
			return 0, 0, errors.Wrap(err, "invalid offset param")
		}
		if offset >= 0 {
			defaultParams.offset = offset
		}
	}

	return defaultParams.count, defaultParams.offset, nil
}
