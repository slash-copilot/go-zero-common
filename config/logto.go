package config

type LogtoAppConfig struct {
	AppID    string
	AppName  string
	Issuer   string
	Audience string
	Endpoint string
	JwksUri  string
}

type LogtoMachineConfig struct {
	AppID    string
	AppName  string
	Issuer   string
	Audience string
	Endpoint string
	JwksUri  string

	AppSecret     string
	TokenEndpoint string
	APIResource   string
}

type LogtoWebhookConfig struct {
	WebhookSigningKey string
}
