package auth

type Claims struct {
	Subject  string   `json:"sub"`
	Scope    string   `json:"scope"`
	Audience []string `json:"aud"`
}
