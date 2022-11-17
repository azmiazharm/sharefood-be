// +heroku goVersion go1.18
// +heroku install ./cmd/...
module sharefood

go 1.18

require (
	cloud.google.com/go/storage v1.24.0
	github.com/Shopify/sarama v1.34.1
	github.com/aws/aws-sdk-go v1.31.9
	github.com/brainlabs/snowflake v0.3.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pressly/goose/v3 v3.6.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/thedevsaddam/govalidator v1.9.10
	github.com/xdg/scram v1.0.3
	github.com/ziutek/mymysql v1.5.4
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	google.golang.org/api v0.85.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.32.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go v0.102.1 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/DataDog/datadog-go v4.8.1+incompatible // indirect
	github.com/DataDog/sketches-go v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/gookit/validate => github.com/brainlabs/validate v1.4.3-0.20220818053512-60a9b27a9b74
