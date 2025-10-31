module github.com/ocfl-archive/dlza-manager-storage-handler

go 1.25.3

require (
	emperror.dev/emperror v0.33.0
	emperror.dev/errors v0.8.1
	github.com/BurntSushi/toml v1.5.0
	github.com/aws/aws-sdk-go-v2/config v1.31.16
	github.com/aws/aws-sdk-go-v2/credentials v1.18.20
	github.com/aws/aws-sdk-go-v2/service/s3 v1.89.1
	github.com/aws/smithy-go v1.23.1
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.11.0
	github.com/je4/filesystem/v3 v3.0.34
	github.com/je4/trustutil/v2 v2.0.31
	github.com/je4/utils/v2 v2.0.60
	github.com/ocfl-archive/dlza-manager v1.0.3-beta2
	github.com/ocfl-archive/dlza-manager-handler v1.0.3-beta2
	github.com/ocfl-archive/error v1.0.5
	github.com/ocfl-archive/gocfl/v2 v2.0.6-beta12
	github.com/ocfl-archive/indexer/v3 v3.0.20
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.34.0
	github.com/tus/tusd/v2 v2.8.0
	gitlab.switch.ch/ub-unibas/go-ublogger/v2 v2.0.2-0.20250331093945-4f0f0ce8c72d
	go.ub.unibas.ch/cloud/certloader/v2 v2.0.22
	go.ub.unibas.ch/cloud/genericproto/v2 v2.0.4
	go.ub.unibas.ch/cloud/miniresolver/v2 v2.0.51
	golang.org/x/net v0.46.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/atsushinee/go-markdown-generator v0.0.0-20231027094725-92d26ffbe778 // indirect
	github.com/aws/aws-sdk-go-v2 v1.39.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.39.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.2 // indirect
	github.com/bytedance/sonic/loader v0.4.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davidbyttow/govips/v2 v2.16.0 // indirect
	github.com/dgraph-io/badger/v4 v4.8.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.3.0 // indirect
	github.com/digitorus/pkcs7 v0.0.0-20250730155240-ffadbf3f398c // indirect
	github.com/digitorus/timestamp v0.0.0-20250524132541-c45532741eea // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/gin-contrib/multitemplate v1.1.1 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.2 // indirect
	github.com/go-openapi/spec v0.22.0 // indirect
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/jsonname v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/loading v0.25.1 // indirect
	github.com/go-openapi/swag/stringutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.28.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gomiran/volmgmt v0.0.0-20221201020756-5e535b6f4941 // indirect
	github.com/google/certificate-transparency-go v1.3.2 // indirect
	github.com/google/flatbuffers v25.9.23+incompatible // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gosimple/slug v1.15.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/je4/goffmpeg v0.0.0-20220114092308-33ab9986404d // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.95 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/sftp v1.13.10 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.55.0 // indirect
	github.com/richardlehane/characterize v1.0.0 // indirect
	github.com/richardlehane/match v1.0.5 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/richardlehane/siegfried v1.11.2 // indirect
	github.com/richardlehane/xmldetect v1.0.2 // indirect
	github.com/ross-spencer/spargo v0.4.1 // indirect
	github.com/ross-spencer/wikiprov v1.0.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/smallstep/certinfo v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	github.com/swaggo/gin-swagger v1.6.1 // indirect
	github.com/swaggo/swag v1.16.6 // indirect
	github.com/tamerh/xml-stream-parser v1.5.0 // indirect
	github.com/tamerh/xpath v1.0.0 // indirect
	github.com/telkomdev/go-stash v1.0.6 // indirect
	github.com/tink-crypto/tink-go/v2 v2.5.0 // indirect
	github.com/tinylib/msgp v1.5.0 // indirect
	github.com/tobischo/argon2 v0.1.0 // indirect
	github.com/tobischo/gokeepasslib/v3 v3.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.step.sm/crypto v0.73.0 // indirect
	go.ub.unibas.ch/cloud/minikvstore v1.0.2 // indirect
	go.ub.unibas.ch/cloud/minivault/v2 v2.0.27 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.22.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/exp v0.0.0-20251023183803-a4bb9ffd2546 // indirect
	golang.org/x/image v0.32.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	gopkg.in/gographics/imagick.v3 v3.7.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
