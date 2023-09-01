package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/slash-copilot/go-zero-common/config"
	xerrors "github.com/slash-copilot/go-zero-common/errors"
	xhttp "github.com/slash-copilot/go-zero-common/http"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	LOGTO_SIGNATURE_HEADER = "logto-signature-sha-256"
	ErrSignatureMismatch   = errors.New("signature mismatch")
	ErrNoSignatureNotFound = errors.New("no signature found")
)

type logtoWebhookReq struct {
	HookID      string                  `json:"hookId"`
	Application *logtoApplicationEntity `json:"application"`
}

type logtoApplicationEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WebhookAuthMiddleware struct {
	Config *config.LogtoWebhookConfig
}

func NewWebhookAuthMiddleware(c *config.LogtoWebhookConfig) *WebhookAuthMiddleware {
	return &WebhookAuthMiddleware{
		Config: c,
	}
}

func (m *WebhookAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//https://docs.logto.io/docs/recipes/webhooks/securing-your-webhooks/
		signature := r.Header.Get(LOGTO_SIGNATURE_HEADER)

		logx.Infof("WebhookAuthMiddleware signature: %s", signature)

		if signature == "" {
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessCodeUnAuthorized,
				Msg:  "no signature found",
			})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isValid := verify(m.Config.WebhookSigningKey, r, signature)

		if !isValid {
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessCodeUnAuthorized,
				Msg:  "signature mismatch",
			})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// ensure app id
		if m.Config.WebhookAppID != "" && !ensureAppId(r, m.Config.WebhookAppID) {
			logx.Infof("WebhookAuthMiddleware ensureAppId failed, expected app id: %s", m.Config.WebhookAppID)
			xhttp.JsonBaseResponseCtx(r.Context(), w, &xerrors.CodeMsg{
				Code: xhttp.BusinessMsgOk,
			})
			w.WriteHeader(http.StatusAccepted)
			return
		}

		next(w, r)
	}
}

// verify checks the expected signature against the computed one using HMAC and SHA256
func verify(signingKey string, r *http.Request, expectedSignature string) bool {
	var buffer bytes.Buffer

	// TeeReader returns a Reader that writes to buffer what it reads from r.Body.
	tee := io.TeeReader(r.Body, &buffer)

	bodyBytes, err := io.ReadAll(tee)

	if err != nil {
		logx.Errorf("Failed to read request body: %s", err.Error())
	}

	defer r.Body.Close()

	// Replace the original body with our buffer
	r.Body = io.NopCloser(&buffer)

	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(bodyBytes)

	computedSignature := hex.EncodeToString(mac.Sum(nil))
	return computedSignature == expectedSignature
}

// ensureAppId checks the expected app id against the one in the request body
func ensureAppId(r *http.Request, expectedAppID string) bool {
	var buffer bytes.Buffer

	// TeeReader returns a Reader that writes to buffer what it reads from r.Body.
	tee := io.TeeReader(r.Body, &buffer)

	bodyBytes, err := io.ReadAll(tee)

	if err != nil {
		logx.Errorf("Failed to read request body: %s", err.Error())
	}
	defer r.Body.Close()

	// Replace the original body with our buffer
	r.Body = io.NopCloser(&buffer)

	// decode body
	var req logtoWebhookReq

	err = json.Unmarshal(bodyBytes, &req)

	if err != nil {
		logx.Errorf("Failed to unmarshal request body: %s", err.Error())
		return false
	}

	if req.Application == nil {
		logx.Errorf("Failed to get application from request body")
		return false
	}

	return req.Application.ID == expectedAppID
}
