package models

// 4.1.1. "iss" (Issuer) Claim
// 4.1.2. "sub" (Subject) Claim
// 4.1.3. "aud" (Audience) Claim
// 4.1.4. "exp" (Expiration Time) Claim
// 4.1.5. "nbf" (Not Before) Claim
// 4.1.6. "iat" (Issued At) Claim
// 4.1.7. "jti" (JWT ID) Claim
type Claims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
	Iat int64  `json:"iat"`
	Jti string `json:"jti"`
}
