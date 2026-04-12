module github.com/neurochar/backend

go 1.25.7

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250425153114-8976f5be98c1.1
	buf.build/go/protovalidate v0.12.0
	github.com/CloudyKit/jet/v6 v6.3.1
	github.com/LastPossum/kamino v0.0.2
	github.com/Masterminds/squirrel v1.5.4
	github.com/aws/aws-sdk-go-v2 v1.36.6
	github.com/aws/aws-sdk-go-v2/config v1.29.18
	github.com/aws/aws-sdk-go-v2/credentials v1.17.71
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.85
	github.com/aws/aws-sdk-go-v2/service/s3 v1.84.1
	github.com/disintegration/imaging v1.6.2
	github.com/gabriel-vasile/mimetype v1.4.9
	github.com/georgysavva/scany/v2 v2.1.4
	github.com/go-playground/validator/v10 v10.26.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/govalues/decimal v0.1.36
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/invopop/jsonschema v0.13.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/k3a/html2text v1.2.1
	github.com/openai/openai-go/v3 v3.31.0
	github.com/pashagolub/pgxmock/v4 v4.8.0
	github.com/pressly/goose/v3 v3.24.3
	github.com/samber/lo v1.51.0
	github.com/sethvargo/go-password v0.3.1
	github.com/stretchr/testify v1.10.0
	github.com/ulule/limiter/v3 v3.11.2
	go.uber.org/fx v1.24.0
	golang.org/x/crypto v0.47.0
	golang.org/x/net v0.49.0
	google.golang.org/genproto v0.0.0-20260226221140-a57be14db171
	google.golang.org/genproto/googleapis/api v0.0.0-20260217215200-42d3e9bedb6d
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260217215200-42d3e9bedb6d
	google.golang.org/grpc v1.79.1
	google.golang.org/protobuf v1.36.11
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	cel.dev/expr v0.25.1 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.11 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.34.1 // indirect
	github.com/aws/smithy-go v1.22.4 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
