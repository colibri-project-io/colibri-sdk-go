package restclient

type RestClientConfig struct {
	Name                string
	BaseURL             string
	Timeout             uint
	Retries             uint8
	RetrySleepInSeconds uint
}
