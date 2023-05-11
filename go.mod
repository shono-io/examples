module github.com/shono-io/shono-examples

go 1.20

replace github.com/shono-io/go-shono v0.0.0 => ../go-shono

require (
	github.com/arangodb/go-driver v0.0.0-20200618111046-f3a9751e1cf5
	github.com/compose-spec/compose-go v1.12.0
	github.com/shono-io/go-shono v0.0.0
	github.com/sirupsen/logrus v1.9.0
	github.com/twmb/franz-go v1.13.2
	github.com/twmb/franz-go/pkg/sr v0.0.0-20230414014213-9e5db4dab85b
)

require (
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/invopop/jsonschema v0.7.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.4.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
