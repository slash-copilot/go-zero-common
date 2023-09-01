package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/logto-io/go/core"
	"github.com/slash-copilot/go-zero-common/config"
	"github.com/slash-copilot/go-zero-common/ctxdata"
	xerrors "github.com/slash-copilot/go-zero-common/errors"
	xhttp "github.com/slash-copilot/go-zero-common/http"
	"github.com/slash-copilot/go-zero-common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type LogtoAuthMiddleware struct {
	Config     *config.LogtoAppConfig
	jwksKeySet *jose.JSONWebKeySet
}

func NewLogtoAuthMiddleware(c *config.LogtoAppConfig) *LogtoAuthMiddleware {
	m := &LogtoAuthMiddleware{
		Config: c,
	}

	k, err := m.CreateRemoteJwks()

	if err != nil {
		logx.Errorf("failed to fetch remote jwks: %v", err)
	}

	m.jwksKeySet = k
	return m
}

func (m *LogtoAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracts a bearer token from Authorization header
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			logx.Errorf("no authorization found")
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessCodeUnAuthorized,
				Msg:  "no authorization found",
			})
			return
		}

		// strips 'Bearer ' prefix from bearer token string
		token, err := utils.StripBearerPrefixFromTokenString(authorization)

		if err != nil {
			logx.Errorf("prefix from bearer token string, no authorization found")
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessCodeUnAuthorized,
				Msg:  "no authorization found",
			})
			return
		}

		claims, err := m.VerifyToken(token)

		if err != nil {
			logx.Errorf("verify token failed: %s,  error: %v", token, err)
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessCodeUnAuthorized,
				Msg:  "verify token failed",
			})
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, ctxdata.CtxKeyJwtUserId{}, claims.Sub)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *LogtoAuthMiddleware) CreateRemoteJwks() (*jose.JSONWebKeySet, error) {
	httpClient := &http.Client{}
	jwksResponse, err := core.FetchJwks(httpClient, m.Config.JwksUri)

	if err != nil {
		return nil, err
	}

	jwks := jose.JSONWebKeySet{}
	for _, rawJsonWebKeyData := range jwksResponse.Keys {
		// Note: convert rawJsonWebKeyData to JSON string for we need to unmarshal it to JSONWebKey
		rawJsonWebKeyJsonString, parseToJsonWebKeyJsonErr := json.Marshal(rawJsonWebKeyData)
		if parseToJsonWebKeyJsonErr != nil {
			return nil, parseToJsonWebKeyJsonErr
		}

		jwk := jose.JSONWebKey{}
		// Note: Use rawJsonWebKeyJsonString to construct the JsonWebKey
		parseToJsonWebKeyErr := jwk.UnmarshalJSON(rawJsonWebKeyJsonString)
		if parseToJsonWebKeyErr != nil {
			return nil, parseToJsonWebKeyErr
		}

		jwks.Keys = append(jwks.Keys, jwk)
	}
	return &jwks, nil
}

func (m *LogtoAuthMiddleware) VerifyToken(token string) (*core.IdTokenClaims, error) {
	if m.jwksKeySet == nil {
		return nil, ErrJwksSetNotFound
	}
	claim, err := verifyIdToken(token, m.Config.Audience, m.Config.Issuer, m.jwksKeySet)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

var ISSUED_AT_RESTRICTIONS int64 = 60 // in seconds

var (
	ErrTokenIssuerNotMatch            = errors.New("token issuer not match")
	ErrTokenAudienceNotMatch          = errors.New("token audience not match")
	ErrTokenExpired                   = errors.New("token expired")
	ErrTokenIssuedInTheFuture         = errors.New("token issued in the future")
	ErrTokenIssuedInThePast           = errors.New("token issued in the past")
	ErrCallbackUriNotMatchRedirectUri = errors.New("callback uri not match redirect uri")
	ErrStateNotMatch                  = errors.New("state not match")
	ErrCodeNotFoundInCallbackUri      = errors.New("code not found in callback uri")
	ErrJwksSetNotFound                = errors.New("jwks set not found")
)

func verifyIdToken(idToken, aud, issuer string, jwks *jose.JSONWebKeySet) (*core.IdTokenClaims, error) {
	jws, err := jwt.ParseSigned(idToken)
	if err != nil {
		return nil, err
	}

	// Verify signature
	idTokenClaims := core.IdTokenClaims{}
	verifySignatureError := jws.Claims(jwks, &idTokenClaims)

	if verifySignatureError != nil {
		return nil, verifySignatureError
	}

	// Verify claims
	if idTokenClaims.Iss != issuer {
		return nil, ErrTokenIssuerNotMatch
	}

	if idTokenClaims.Aud != aud {
		return nil, ErrTokenAudienceNotMatch
	}

	now := time.Now().Unix()

	if idTokenClaims.Exp < now {
		return nil, ErrTokenExpired
	}

	if idTokenClaims.Iat > now+ISSUED_AT_RESTRICTIONS {
		return nil, ErrTokenIssuedInTheFuture
	}

	// if idTokenClaims.Iat < now-ISSUED_AT_RESTRICTIONS {
	// 	return nil, ErrTokenIssuedInThePast
	// }

	return &idTokenClaims, nil
}
