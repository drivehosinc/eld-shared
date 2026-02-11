package metadata

import (
	"context"
	"strings"

	"github.com/drivehosinc/eld-shared/logger"
	"github.com/drivehosinc/eld-shared/postgresql"
	"google.golang.org/grpc/metadata"
)

// ToContext stores Metadata in the context and also sets logger context
// values (request_id, is_debug) so that logging works automatically.
func ToContext(ctx context.Context, md *Metadata) context.Context {
	ctx = context.WithValue(ctx, metadataKey{}, md)

	if md.RequestId != "" {
		ctx = logger.WithRequestID(ctx, md.RequestId)
	}
	if md.IsDebug {
		ctx = logger.WithDebug(ctx)
	}

	if md.Replica {
		ctx = postgresql.WithReplica(ctx)
	}

	return ctx
}

// FromContext retrieves Metadata from the context.
// Returns nil if no metadata is stored.
func FromContext(ctx context.Context) *Metadata {
	md, _ := ctx.Value(metadataKey{}).(*Metadata)
	return md
}

// ExtractFromIncoming reads gRPC incoming metadata from the context,
// parses it into a Metadata struct, stores it in the context, and
// sets logger context values.
func ExtractFromIncoming(ctx context.Context) (context.Context, *Metadata) {
	md := &Metadata{}

	grpcMD, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md.RequestId = firstValue(grpcMD, KeyRequestID)
		md.IsDebug = firstValue(grpcMD, KeyIsDebug) == "true"
		md.Replica = firstValue(grpcMD, KeyReplica) == "true"
		md.CompanyIds = grpcMD.Get(KeyCompanyIDs)

		// company_ids may be sent as a single comma-separated value
		if len(md.CompanyIds) == 1 && strings.Contains(md.CompanyIds[0], ",") {
			md.CompanyIds = strings.Split(md.CompanyIds[0], ",")
		}
	}

	return ToContext(ctx, md), md
}

// InjectToOutgoing appends Metadata fields into the gRPC outgoing
// metadata so they propagate to downstream services.
func InjectToOutgoing(ctx context.Context, md *Metadata) context.Context {
	pairs := make([]string, 0, 8)

	if md.RequestId != "" {
		pairs = append(pairs, KeyRequestID, md.RequestId)
	}
	if md.IsDebug {
		pairs = append(pairs, KeyIsDebug, "true")
	}
	if md.Replica {
		pairs = append(pairs, KeyReplica, "true")
	}
	for _, id := range md.CompanyIds {
		pairs = append(pairs, KeyCompanyIDs, id)
	}

	if len(pairs) == 0 {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

// firstValue returns the first value for a gRPC metadata key, or "".
func firstValue(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
