module github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis

// +heroku goVersion go1.17.2
go 1.17

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/sys v0.0.0-20211106132015-ebca88c72f68 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/aadhaarapi v0.0.0-20211106145727-c50dbea70cfa
	github.com/beevik/etree v1.1.0
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/go-redis/redis/v8 v8.11.4
	github.com/otiai10/gosseract/v2 v2.3.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/russellhaering/goxmldsig v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/yeka/zip v0.0.0-20180914125537-d046722c6feb
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/ma314smith/signedxml v0.0.0-20210628192057-abc5b481ae1c // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// replace github.com/aaabhilash97/aadhaar-paperless-offline-ekyc-apis/pkg/aadhaarapi => ./pkg/aadhaarapi
