package server

import "context"

type errorHolderKey struct{}

type ErrorHolder struct {
	Err error
}

func WithErrorHolder(ctx context.Context) (context.Context, *ErrorHolder) {
	holder := &ErrorHolder{}
	return context.WithValue(ctx, errorHolderKey{}, holder), holder
}

func SetError(ctx context.Context, err error) {
	holder, ok := ctx.Value(errorHolderKey{}).(*ErrorHolder)
	if !ok || holder == nil {
		return
	}
	holder.Err = err
}

func GetError(ctx context.Context) error {
	holder, ok := ctx.Value(errorHolderKey{}).(*ErrorHolder)
	if !ok || holder == nil {
		return nil
	}
	return holder.Err
}
