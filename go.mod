module github.com/shono-io/shono-examples

go 1.20

replace github.com/shono-io/go-shono v0.0.0 => ../go-shono

replace github.com/shono-io/shono-ccloud v1.0.0 => ../shono-ccloud
replace github.com/shono-io/shono v0.0.0 => ../shono

require (
	github.com/compose-spec/compose-go v1.12.0
	github.com/shono-io/shono v0.0.0
	github.com/sirupsen/logrus v1.9.0
	github.com/twmb/franz-go v1.13.2
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/shono-io/shono-ccloud v1.0.0 // indirect
	github.com/twmb/franz-go/pkg/kadm v1.8.1 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.4.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783 // indirect
	golang.org/x/sys v0.8.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
