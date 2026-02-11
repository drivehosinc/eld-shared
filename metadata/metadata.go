package metadata

type Metadata struct {
	RequestId  string   `json:"request_id"`
	IsDebug    bool     `json:"is_debug"`
	Replica    bool     `json:"replica"`
	CompanyIds []string `json:"company_ids"`
}


