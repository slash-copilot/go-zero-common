package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/slash-copilot/go-zero-common/config"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	LOGTO_SIGNATURE_HEADER = "logto-signature-sha-256"
	ErrSignatureMismatch   = errors.New("signature mismatch")
	ErrNoSignatureNotFound = errors.New("no signature found")
)

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
			unauthorized(w, r, ErrNoSignatureNotFound)
			return
		}

		isValid := Verify(m.Config.WebhookSigningKey, r, signature)

		if !isValid {
			unauthorized(w, r, ErrSignatureMismatch)
			return
		}

		next(w, r)
	}
}

// Verify checks the expected signature against the computed one using HMAC and SHA256
func Verify(signingKey string, r *http.Request, expectedSignature string) bool {
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
