module svrdemo

go 1.13

replace environment => ../../environment/src

replace mmapcache => ../../mmapcache/src

replace single => ../../single/src

require (
	environment v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/satori/go.uuid v1.2.0
	google.golang.org/grpc v1.25.1
	mmapcache v0.0.0-00010101000000-000000000000
)
