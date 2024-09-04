package models

type ServerConfig struct {
	Addr    string `toml:"addr"`           // server will start at this "host:port"
	ExtAddr string `toml:"extaddr"`        // server will assume running at this base url
	TLSCert string `toml:"certificate"`    // TLS Certificate
	TLSKey  string `toml:"certificatekey"` // TLS Certificate Private Key
}
