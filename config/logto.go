package config

type LogtoAppConfig struct {
	AppID    string `json:",optional"`
	AppName  string `json:",optional"`
	Issuer   string `json:",optional"`
	Audience string `json:",optional"`
	Endpoint string `json:",optional"`
	JwksUri  string `json:",optional"`
}

type LogtoMachineConfig struct {
	AppID    string `json:",optional"`
	AppName  string `json:",optional"`
	Issuer   string `json:",optional"`
	Audience string `json:",optional"`
	Endpoint string `json:",optional"`
	JwksUri  string `json:",optional"`
	
	AppSecret     string `json:",optional"`
	TokenEndpoint string `json:",optional"`
	APIResource   string `json:",optional"`
}

type LogtoWebhookConfig struct {
	WebhookSigningKey string `json:",optional"`
	WebhookID      string `json:",optional"`
}
