package main

import (
	"fmt"

	"environment/cfgargs"
	"environment/dump"
	"svrdemo/app"
	pb "svrdemo/proto"
	"svrdemo/server"

	uuid "github.com/satori/go.uuid"
)

var (
	BuildVersion = ""
)

func main() {
	// // 这段是测试代码为了demo rpc调用与mmapcache START
	// dump.InitDump(true, func(d string) { fmt.Println(d) })
	// go func() {
	// 	for i := 0; i < 5; i++ {
	// 		<-time.After(time.Second)
	// 		fmt.Printf("wait for start test rpg client(%v)...\n", i)
	// 	}
	// 	exampleGrpcClient()
	// }()
	// // 这段是测试代码为了demo rpc调用与mmapcache END

	srvCfg, err := cfgargs.InitSrvConfig(BuildVersion, func() {
		// user flag binding code
	})
	if nil != err {
		fmt.Println(err)
		return
	}
	app.GetApp().InitApp(srvCfg)

	srv := server.NewServer()
	srv.Run(srvCfg.Addr)
}

func exampleGrpcClient() {
	for index := 0; index < 100; index++ {
		go func() {
			for {
				dump.PacketRequestCounter()
				transid := uuid.NewV1()
				req := &pb.SimpleHello{
					Transid: transid.String(),
					Name:    "name",
				}

				srv := app.GetApp().SrvDemo
				sayHelloResponse, err := srv.SayHello(srv.GetCtx(), req)
				if nil != err {
					fmt.Printf("resp:%v err:%v\n", sayHelloResponse, err)
				}
				dump.PacketResponseCounter()
			}
		}()
	}
}
