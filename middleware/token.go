package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenClaims struct {
	UserID   uuid.UUID
	Role     string
	Username string
	jwt.RegisteredClaims
}
type JWTServiceImpl interface {
	CreateToken(claims *TokenClaims) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	ValidateContext(next http.Handler) http.Handler
}
type JWTService struct {
	secretKey []byte
	zap       *zap.Logger
}

func NewJWTService(key []byte, zap *zap.Logger) *JWTService {
	return &JWTService{
		secretKey: key,
		zap:       zap,
	}
}

func (js *JWTService) CreateToken(claims *TokenClaims) (string, error) {
	newClaims := TokenClaims{
		UserID:   claims.UserID,
		Role:     claims.Role,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(js.secretKey)
	if err != nil {
		js.zap.Error(utils.ErrInternal.Error(), zap.Error(err))
		return "", fmt.Errorf("%v: failed to create token", utils.ErrInternal)
	}
	return tokenString, nil
}

func (js *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	newClaims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, newClaims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return js.secretKey, nil
	})
	if err != nil {
		js.zap.Error(utils.ErrInternal.Error(), zap.Error(err))
		return nil, utils.ErrInternal
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.UserID == uuid.Nil || claims.Username == "" || claims.Role == "" {
			js.zap.Warn("Token missing required claims", zap.Any("claims", claims))
			return nil, utils.ErrUnauthorized
		}
		return claims, nil
	}
	return nil, utils.ErrUnauthorized
}

func (js *JWTService) ValidateContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			js.zap.Warn("Missing Token")
			utils.JSONResponse(w, http.StatusUnauthorized, utils.ErrUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := js.ValidateToken(token)
		if err != nil {
			js.zap.Error("failed to validate token", zap.Error(err))
			utils.JSONResponse(w, http.StatusUnauthorized, utils.ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, utils.RoleKey, claims.Role)
		js.zap.Info("Success", zap.Any("user", ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
