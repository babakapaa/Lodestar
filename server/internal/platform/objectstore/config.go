package objectstore

type Config struct {
	Endpoint       string
	PublicEndpoint string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
}

func (c Config) Enabled() bool {
	return c.Endpoint != ""
}
