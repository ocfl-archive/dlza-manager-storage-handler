module github.com/ocfl-archive/dlza-manager-storage-handler

go 1.23.5

require (
	emperror.dev/emperror v0.33.0
	emperror.dev/errors v0.8.1
	github.com/BurntSushi/toml v1.4.0
	github.com/aws/aws-sdk-go-v2/config v1.29.1
	github.com/aws/aws-sdk-go-v2/credentials v1.17.54
	github.com/aws/aws-sdk-go-v2/service/s3 v1.73.2
	github.com/aws/smithy-go v1.22.1
	github.com/gin-contrib/cors v1.7.3
	github.com/gin-gonic/gin v1.10.0
	github.com/je4/filesystem/v3 v3.0.22
	github.com/je4/indexer/v3 v3.0.7
	github.com/je4/trustutil/v2 v2.0.28
	github.com/je4/utils/v2 v2.0.56
	github.com/ocfl-archive/dlza-manager v1.0.0-beta1
	github.com/ocfl-archive/dlza-manager-handler v1.0.0-beta1
	github.com/ocfl-archive/gocfl/v2 v2.0.6-beta6
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/tus/tusd/v2 v2.6.0
	gitlab.switch.ch/ub-unibas/go-ublogger/v2 v2.0.1
	go.ub.unibas.ch/cloud/certloader/v2 v2.0.18
	go.ub.unibas.ch/cloud/genericproto/v2 v2.0.4
	go.ub.unibas.ch/cloud/miniresolver/v2 v2.0.31
	golang.org/x/net v0.34.0
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.3
)

require (
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/atsushinee/go-markdown-generator v0.0.0-20231027094725-92d26ffbe778 // indirect
	github.com/aws/aws-sdk-go-v2 v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.9 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/bytedance/sonic v1.12.7 // indirect
	github.com/bytedance/sonic/loader v0.2.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davidbyttow/govips/v2 v2.16.0 // indirect
	github.com/dgraph-io/badger/v4 v4.5.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/multitemplate v1.0.2 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gomiran/volmgmt v0.0.0-20221201020756-5e535b6f4941 // indirect
	github.com/google/certificate-transparency-go v1.3.1 // indirect
	github.com/google/flatbuffers v24.12.23+incompatible // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/tink/go v1.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gosimple/slug v1.15.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/je4/goffmpeg v0.0.0-20220114092308-33ab9986404d // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.84 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/ocfl-archive/error v1.0.2 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/sftp v1.13.7 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/richardlehane/characterize v1.0.0 // indirect
	github.com/richardlehane/match v1.0.5 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/richardlehane/siegfried v1.11.1 // indirect
	github.com/richardlehane/xmldetect v1.0.2 // indirect
	github.com/ross-spencer/spargo v0.4.1 // indirect
	github.com/ross-spencer/wikiprov v0.2.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/smallstep/certinfo v1.13.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	github.com/swaggo/gin-swagger v1.6.0 // indirect
	github.com/swaggo/swag v1.16.4 // indirect
	github.com/tamerh/xml-stream-parser v1.5.0 // indirect
	github.com/tamerh/xpath v1.0.0 // indirect
	github.com/telkomdev/go-stash v1.0.6 // indirect
	github.com/tobischo/argon2 v0.1.0 // indirect
	github.com/tobischo/gokeepasslib/v3 v3.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.step.sm/crypto v0.57.0 // indirect
	go.ub.unibas.ch/cloud/minivault/v2 v2.0.15 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.13.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	golang.org/x/image v0.23.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.29.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	gopkg.in/gographics/imagick.v3 v3.7.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
)
