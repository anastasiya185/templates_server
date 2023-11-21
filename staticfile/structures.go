package staticfile

type Credentials struct {
	Host                  string `json:"host"`
	AuthenticationEnabled bool   `json:"authentication_enabled"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Status                string `json:"status"`
}
type EnvironmentConfig struct {
	Prod struct {
		Elasticsearch Credentials `json:"elasticsearch"`
		Kibana        Credentials `json:"kibana"`
	} `json:"prod"`
	Mon struct {
		Elasticsearch Credentials `json:"elasticsearch"`
	} `json:"mon"`
}

type Status struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type StatusData struct {
	Prod struct {
		Elasticsearch Status `json:"elasticsearch"`
		Kibana        Status `json:"kibana"`
	} `json:"prod"`
	Mon struct {
		Elasticsearch Status `json:"elasticsearch"`
	} `json:"mon"`
}
