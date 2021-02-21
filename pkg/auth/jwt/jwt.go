package jwt

import (
	"errors"
	"falcon-seed/pkg/auth/rbac"
	"falcon-seed/pkg/config"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var minSecretLen = 128

type User struct {
	Username    string
	Email       string
	AccessLevel rbac.AccessRole
}

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	// Secret key used for signing.
	secret string
	key    []byte

	// Duration for which the jwt token is valid.
	ttlMinutes int
	ttl        time.Duration

	// JWT signing algorithm
	algoMethod string
	algo       jwt.SigningMethod
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	accessToken  Token
	refreshToken Token
}

func NewToken(algo, secret string, ttlMinutes, minSecretLength int) (Token, error) {
	token, err := validateTokenParams(Token{
		secret:     secret,
		ttlMinutes: ttlMinutes,
		algoMethod: algo,
	}, minSecretLength)

	if err != nil {
		return Token{}, err
	}

	return token, err
}

// New generates new JWT service necessary for auth middleware
func New(config *config.JWT) (Service, error) {
	accessToken, err := NewToken(
		config.AccessToken.SigningAlgorithm,
		os.Getenv("JWT_SECRET"),
		config.AccessToken.DurationMinutes,
		config.AccessToken.MinSecretLength,
	)
	if err != nil {
		return Service{}, err
	}

	refreshToken, err := NewToken(
		config.RefreshToken.SigningAlgorithm,
		os.Getenv("JWT_SECRET"),
		config.RefreshToken.DurationMinutes,
		config.RefreshToken.MinSecretLength,
	)
	if err != nil {
		return Service{}, err
	}

	return Service{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}, nil
}

func validateTokenParams(token Token, minSecretLength int) (Token, error) {
	if minSecretLength > 0 {
		minSecretLen = minSecretLength
	}

	if len(token.secret) < minSecretLen {
		return Token{}, fmt.Errorf("jwt secret length is %v, which is less than required %v",
			len(token.secret), minSecretLen)
	}
	signingMethod := jwt.GetSigningMethod(token.algoMethod)
	if signingMethod == nil {
		return Token{}, fmt.Errorf("invalid jwt signing method: %s", token.algoMethod)
	}

	parsedToken := token
	parsedToken.key = []byte(token.secret)
	parsedToken.algo = signingMethod
	parsedToken.ttl = time.Duration(token.ttlMinutes) * time.Minute

	return parsedToken, nil
}

// ParseAccessToken parses access token from Authorization header
func (service Service) ParseAccessToken(authHeader string) (*jwt.Token, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("no token")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if service.accessToken.algo != token.Method {
			return nil, errors.New("bad token")
		}
		return service.accessToken.key, nil
	})
}

// ParseRefreshToken parses refresh token from Authorization header
func (service Service) ParseRefreshToken(refreshToken string) (*jwt.Token, error) {
	return nil, nil
}

// GenerateToken generates new JWT token and populates it with user data
func (service Service) GenerateAccessToken(user User) (string, error) {
	accessToken, err := jwt.NewWithClaims(service.accessToken.algo, jwt.MapClaims{
		"user":  user.Username,
		"email": user.Email,
		"role":  user.AccessLevel,
		"exp":   time.Now().Add(service.accessToken.ttl).Unix(),
	}).SignedString(service.accessToken.key)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (service Service) GenerateRefreshToken(user string) (string, error) {
	refreshToken, err := jwt.NewWithClaims(service.refreshToken.algo, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(service.refreshToken.ttl).Unix(),
	}).SignedString(service.refreshToken.key)

	if err != nil {
		return "", err
	}

	return refreshToken, err
}
