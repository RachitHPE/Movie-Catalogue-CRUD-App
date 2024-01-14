package server

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWTParameters - params to configure JWT
type JWTParameters struct {
	Algorithm     string
	AccessKey     []byte
	AccessKeyTTL  int
	RefreshKey    []byte
	RefreshKeyTTL int
	PrivKeyECDSA  *ecdsa.PrivateKey
	PubKeyECDSA   *ecdsa.PublicKey
	PrivKeyRSA    *rsa.PrivateKey
	PubKeyRSA     *rsa.PublicKey

	Audience string
	Issuer   string
	AccNbf   int
	RefNbf   int
	Subject  string
}

// JWTParams - exported variables
var JWTParams JWTParameters

// MyCustomClaims ...
type MyCustomClaims struct {
	AuthID uint64 `json:"authID,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	Scope  string `json:"scope,omitempty"`
}

// JWTClaims ...
type JWTClaims struct {
	MyCustomClaims
	jwt.RegisteredClaims
}

// JWTPayload ...
type JWTPayload struct {
	AccessJWT   string `json:"accessJWT,omitempty"`
	RefreshJWT  string `json:"refreshJWT,omitempty"`
	TwoAuth     string `json:"twoFA,omitempty"`
	RecoveryKey string `json:"recoveryKey,omitempty"`
}

// JWT - validate access token
func JWTConfiguration() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtPayload JWTPayload
		var val string
		var vals []string

		// first try to read the cookie
		accessJWT, err := c.Cookie("accessJWT")
		// accessJWT is available in the cookie
		if err == nil {
			jwtPayload.AccessJWT = accessJWT
			errVal := verifyClaims(c, jwtPayload)

			if errVal != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "failed to verify claims")

				return
			}
		}

		// accessJWT is not available in the cookie
		// try to read the Authorization header
		val = c.Request.Header.Get("Authorization")
		if len(val) == 0 || !strings.Contains(val, "Bearer") {
			// no vals or no bearer found
			c.AbortWithStatusJSON(http.StatusUnauthorized, "token missing")
			return
		}
		vals = strings.Split(val, " ")
		// Authorization: Bearer {access} => length is 2
		// Authorization: Bearer {access} {refresh} => length is 3
		if len(vals) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "token missing")
			return
		}

		jwtPayload.AccessJWT = vals[1]
	}
}

func verifyClaims(c *gin.Context, jwtPayload JWTPayload) error {
	token, err := jwt.ParseWithClaims(jwtPayload.AccessJWT, &JWTClaims{}, ValidateAccessJWT)
	if err != nil {
		return fmt.Errorf("error is : %w", err)
	}

	_, ok := token.Claims.(*JWTClaims)
	if !ok {
		return fmt.Errorf("error is : %w", err)
	}

	c.Next()

	return nil
}

// ValidateHMACAccess - validate hash based access token
func ValidateHMACAccess(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JWTParams.AccessKey, nil
}

// ValidateECDSA - validate elliptic curve digital signature algorithm based token
func ValidateECDSA(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JWTParams.PubKeyECDSA, nil
}

// ValidateRSA - validate Rivest–Shamir–Adleman cryptosystem based token
func ValidateRSA(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JWTParams.PubKeyRSA, nil
}

// ValidateAccessJWT - verify the access JWT's signature, and validate its claims
func ValidateAccessJWT(token *jwt.Token) (interface{}, error) {
	alg := JWTParams.Algorithm

	switch alg {
	case "HS256", "HS384", "HS512":
		return ValidateHMACAccess(token)
	case "ES256", "ES384", "ES512":
		return ValidateECDSA(token)
	case "RS256", "RS384", "RS512":
		return ValidateRSA(token)
	default:
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
}

// GetJWT - issue new tokens
func GetJWT(customClaims MyCustomClaims, tokenType string) (string, string, error) {
	var (
		key          []byte
		privKeyECDSA *ecdsa.PrivateKey
		privKeyRSA   *rsa.PrivateKey
		ttl          int
		nbf          int
	)

	if tokenType == "access" {
		key = JWTParams.AccessKey
		ttl = JWTParams.AccessKeyTTL
		nbf = JWTParams.AccNbf
	}
	// Create the Claims
	claims := JWTClaims{
		MyCustomClaims{
			AuthID: customClaims.AuthID,
			Email:  customClaims.Email,
			Role:   customClaims.Role,
			Scope:  customClaims.Scope,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(ttl))),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    JWTParams.Issuer,
			Subject:   JWTParams.Subject,
		},
	}

	if JWTParams.Audience != "" {
		claims.Audience = []string{JWTParams.Audience}
	}
	if nbf > 0 {
		claims.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(nbf)))
	}

	var token *jwt.Token
	alg := JWTParams.Algorithm

	switch alg {
	case "HS256":
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	case "HS384":
		token = jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	case "HS512":
		token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	case "ES256":
		privKeyECDSA = JWTParams.PrivKeyECDSA
		token = jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	case "ES384":
		privKeyECDSA = JWTParams.PrivKeyECDSA
		token = jwt.NewWithClaims(jwt.SigningMethodES384, claims)
	case "ES512":
		privKeyECDSA = JWTParams.PrivKeyECDSA
		token = jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	case "RS256":
		privKeyRSA = JWTParams.PrivKeyRSA
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	case "RS384":
		privKeyRSA = JWTParams.PrivKeyRSA
		token = jwt.NewWithClaims(jwt.SigningMethodRS384, claims)
	case "RS512":
		privKeyRSA = JWTParams.PrivKeyRSA
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	default:
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	}

	var jwtValue string
	var err error

	// HMAC
	if alg == "HS256" || alg == "HS384" || alg == "HS512" {
		jwtValue, err = token.SignedString(key)
		if err != nil {
			return "", "", err
		}
	}

	// ECDSA
	//
	// ES256
	// prime256v1: X9.62/SECG curve over a 256 bit prime field, also known as P-256 or NIST P-256
	// widely used, recommended for general-purpose cryptographic operations
	// openssl ecparam -name prime256v1 -genkey -noout -out private-key.pem
	// openssl ec -in private-key.pem -pubout -out public-key.pem
	//
	// ES384
	// secp384r1: NIST/SECG curve over a 384 bit prime field
	// openssl ecparam -name secp384r1 -genkey -noout -out private-key.pem
	// openssl ec -in private-key.pem -pubout -out public-key.pem
	//
	// ES512
	// secp521r1: NIST/SECG curve over a 521 bit prime field
	// openssl ecparam -name secp521r1 -genkey -noout -out private-key.pem
	// openssl ec -in private-key.pem -pubout -out public-key.pem
	if alg == "ES256" || alg == "ES384" || alg == "ES512" {
		jwtValue, err = token.SignedString(privKeyECDSA)
		if err != nil {
			return "", "", err
		}

	}

	// RSA
	//
	// RS256
	// openssl genpkey -algorithm RSA -out private-key.pem -pkeyopt rsa_keygen_bits:2048
	// openssl rsa -in private-key.pem -pubout -out public-key.pem
	//
	// RS384
	// openssl genpkey -algorithm RSA -out private-key.pem -pkeyopt rsa_keygen_bits:3072
	// openssl rsa -in private-key.pem -pubout -out public-key.pem
	//
	// RS512
	// openssl genpkey -algorithm RSA -out private-key.pem -pkeyopt rsa_keygen_bits:4096
	// openssl rsa -in private-key.pem -pubout -out public-key.pem
	if alg == "RS256" || alg == "RS384" || alg == "RS512" {
		jwtValue, err = token.SignedString(privKeyRSA)
		if err != nil {
			return "", "", err
		}
	}

	return jwtValue, claims.ID, nil
}

func Cors() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// First, we add the headers with need to enable CORS.
		// Make sure to adjust these headers to your needs.
		ginCtx.Header("Access-Control-Allow-Origin", "*")
		// ginCtx.Header("Access-Control-Allow-Methods", "*")
		// ginCtx.Header("Access-Control-Allow-Headers", "*")
		ginCtx.Header("Content-Type", "application/json")
		ginCtx.Header("Vary", "Origin")

		ginCtx.Next()
	}
}
