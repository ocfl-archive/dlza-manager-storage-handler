package certs

import "embed"

//go:embed ub-log.ub.unibas.ch.key.pem
//go:embed ub-log.ub.unibas.ch.cert.pem
var CertFS embed.FS
