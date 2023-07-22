package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github.com/cockroachdb/errors"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/unrolled/render"
	"net/http"
	"strings"
)

type Authorizer struct {
	renderer *render.Render
	config   *pkg.Config
	logger   *pkg.Logger
	apiKeys  [][]byte
}

func NewAuthorizer(renderer *render.Render, config *pkg.Config, logger *pkg.Logger) *Authorizer {
	apiKeys := make([][]byte, 0)
	for _, value := range config.AppAPIKeys {
		decodedKey, err := hex.DecodeString(value)
		if err != nil {
			panic("failed to decode API keys")
		}
		apiKeys = append(apiKeys, decodedKey)
	}
	return &Authorizer{renderer, config, logger, apiKeys}
}

func (s *Authorizer) isAuthorized(rq *http.Request) bool {
	return true
}

func (s *Authorizer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		apiKey, err := s.getBearerToken(rq)
		if err != nil {
			s.logger.Warn().Err(err).Send()
			_ = s.renderer.JSON(rw, http.StatusUnauthorized, struct {
				Message string `json:"message"`
			}{"Permission denied"})
			return
		}
		if !s.isTokenValid(apiKey, s.apiKeys) {
			s.logger.Warn().Msgf("invalid authorization attempt with key: %s", apiKey)
			_ = s.renderer.JSON(rw, http.StatusUnauthorized, struct {
				Message string `json:"message"`
			}{"Permission denied"})
			return
		}
		next.ServeHTTP(rw, rq)
	})
}

func (s *Authorizer) isTokenValid(rawKey string, availableKeys [][]byte) bool {
	hash := sha256.Sum256([]byte(rawKey))
	key := hash[:]
	for _, value := range availableKeys {
		contentEqual := subtle.ConstantTimeCompare(value, key) == 1
		if contentEqual {
			return true
		}
	}
	return false
}

func (s *Authorizer) getBearerToken(rq *http.Request) (string, error) {
	var token string
	authHeader := rq.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return token, errors.New("malformed auth token")
	}
	token = splitToken[1]
	return token, nil
}
