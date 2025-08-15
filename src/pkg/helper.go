package pkg

type HeaderCtx map[string]string

var HeaderCtxKey = struct{ Name string }{Name: "header_ctx_key"}

type DefaultHeader struct {
	UserId          *uint   `json:"user_id"`
	Currency        *string `json:"currency"`
	Language        *string `json:"language"`
	Timezone        *string `json:"timezone"`
	VisibilityLevel *int8   `json:"visibility_level"`
	FeatureLevel    *int    `json:"feature_level"`
}
