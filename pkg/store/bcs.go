// Copyright (c) The BCS Authors.
// Licensed under the Apache License 2.0.

package store

import (
	"context"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/thanos-io/thanos/pkg/store/storepb"
)

// LabelMatchKey
const LabelMatchKey = ctxKey(1)

// CopyLabelMatchContext copies the necessary trace context from given source context to target context.
func CopyLabelMatchContext(trgt, src context.Context) context.Context {
	v, ok := LabelMatchValue(src)
	if !ok {
		return trgt
	}
	return WithLabelMatchValue(trgt, v)
}

// WithLabelMatchValue 设置值
func WithLabelMatchValue(ctx context.Context, matches [][]*labels.Matcher) context.Context {
	return context.WithValue(ctx, LabelMatchKey, matches)
}

// LabelMatchValue 获取值，获取变量, 修改matcher, 支持 namespace 级别过滤
func LabelMatchValue(ctx context.Context) ([][]*labels.Matcher, bool) {
	v, ok := ctx.Value(LabelMatchKey).([][]*labels.Matcher)
	return v, ok
}

// MakeLabelMatchSeriesRequest
func MakeLabelMatchSeriesRequest(ctx context.Context, r *storepb.SeriesRequest) []*storepb.SeriesRequest {
	matchValues, ok := LabelMatchValue(ctx)
	if !ok {
		return []*storepb.SeriesRequest{r}
	}

	reqs := make([]*storepb.SeriesRequest, 0, len(matchValues))
	for _, v := range matchValues {
		storeMatchers, _ := storepb.PromMatchersToMatchers(v...)
		newReq := *r
		newReq.Matchers = append(newReq.Matchers, storeMatchers...)
		reqs = append(reqs, &newReq)
	}

	return reqs
}
