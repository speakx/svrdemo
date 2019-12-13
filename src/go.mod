module svrdemo

go 1.13

replace environment => ../../environment/src

replace idgenerator => ../../idgenerator/src

replace mmapcache => ../../mmapcache/src

replace single => ../../single/src

replace singledb => ../../singledb/src

replace github.com/panjf2000/gnet => ../../pkg/github.com/panjf2000/gnet

require (
	environment v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/satori/go.uuid v1.2.0
	github.com/smallnest/goframe v0.0.0-20191101094441-1fbd8e51db18
	google.golang.org/grpc v1.25.1
	mmapcache v0.0.0-00010101000000-000000000000
	single v0.0.0-00010101000000-000000000000
)
