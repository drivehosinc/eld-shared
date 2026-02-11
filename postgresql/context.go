package postgresql

import "context"

type ctxKey struct{}

type mode int

const (
	modeMaster  mode = iota
	modeReplica
)

// WithMaster forces the query to use the master pool.
func WithMaster(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, modeMaster)
}

// WithReplica routes the query to the replica pool.
// Falls back to master if no replica is configured.
func WithReplica(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, modeReplica)
}

func modeFromCtx(ctx context.Context) mode {
	if m, ok := ctx.Value(ctxKey{}).(mode); ok {
		return m
	}
	return modeMaster
}
