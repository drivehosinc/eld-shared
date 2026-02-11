package metadata

// gRPC metadata keys used on the wire.
const (
	KeyRequestID  = "request_id"
	KeyIsDebug    = "is_debug"
	KeyReplica    = "replica"
	KeyCompanyIDs = "company_ids"
)

// context key for storing Metadata in context.Value.
type metadataKey struct{}
