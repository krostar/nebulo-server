package auth

// Token store a token and his expiration date
type Token struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
