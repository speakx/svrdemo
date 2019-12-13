package main

import (
	"environment/cfgargs"
	"environment/dump"
	"fmt"
	"os"
	"single/proto/pbsingle"
	"svrdemo/app"
	"svrdemo/genclient"
	"svrdemo/proto/pbsvrdemo"
	"svrdemo/server"
	"time"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

var (
	BuildVersion = ""
)

func main() {
	// 这段是测试代码为了demo rpc调用与mmapcache START
	go func() {
		for i := 0; i < 3; i++ {
			<-time.After(time.Second)
			fmt.Printf("wait for start test rpg client(%v)...\n", i)
		}
		// exampleGrpcClient()
		// exampleSingleClient()
		exampleGenSingleClient()
	}()
	// 这段是测试代码为了demo rpc调用与mmapcache END

	srvCfg, err := cfgargs.InitSrvConfig(BuildVersion, func() {
		// user flag binding code
	})
	if nil != err {
		fmt.Println(err)
		return
	}
	app.GetApp().InitApp(srvCfg)

	srv := server.NewServer()
	srv.Run(srvCfg.Info.Addr)
}

func exampleGrpcClient() {
	for index := 0; index < 100; index++ {
		go func() {
			for {
				dump.NetEventSendIncr(0)
				transid := uuid.NewV1()
				req := &pbsvrdemo.SimpleHello{
					Transid: transid.String(),
					Name:    "name",
				}

				srv := app.GetApp().SrvDemo
				replay, err := srv.SayHello(srv.GetCtx(), req)
				if nil != err {
					fmt.Printf("resp:%v err:%v\n", replay, err)
				}
				dump.NetEventSendDecr(0)
			}
		}()
	}
}

func exampleSingleClient() {
	for index := 0; index < 100; index++ {
		instanceid := uint32(index) << 20 & 0xFFF00000
		go func(id uint64) {
			seqid := 0
			for {
				dump.NetEventSendIncr(0)
				transid := uuid.NewV1()
				req := &pbsingle.Message{
					Transid:  transid.String(),
					ClientId: uint64(time.Now().Unix())<<32 | uint64(id) | uint64(seqid)&0xFFFFF,
					FromUid:  1,
					ToUid:    2,
					Msg:      "hello world.",
					Type:     0,
				}

				srv := app.GetApp().Single
				replay, err := srv.SendMessage(srv.GetCtx(), req)
				if nil != err {
					fmt.Printf("resp:%v err:%v\n", replay, err)
				}
				dump.NetEventSendDecr(0)
			}
		}(uint64(instanceid))
	}
}

func exampleGenSingleClient() {
	for index := 0; index < 100; index++ {
		instanceid := uint32(index) << 20 & 0xFFF00000
		go func(id uint64) {
			seqid := 0
			gsc := &genclient.GenSingleClient{}
			gsc.Connect()

			for {
				dump.NetEventSendIncr(0)
				transid := uuid.NewV1()
				req := &pbsingle.Message{
					Transid:  transid.String(),
					ClientId: uint64(time.Now().Unix())<<32 | uint64(id) | uint64(seqid)&0xFFFFF,
					FromUid:  1,
					ToUid:    2,
					Msg:      "hello world.",
					Type:     0,
				}
				seqid++
				data, _ := proto.Marshal(req)
				// fmt.Printf("id:%v seq:%v start %v\n", id, seqid, data)
				err := gsc.FC.WriteFrame(data)
				// fmt.Printf("id:%v seq:%v end\n", id, seqid)
				if nil != err {
					fmt.Println(err)
					os.Exit(0)
				}
				_, err = gsc.FC.ReadFrame()
				if nil != err {
					fmt.Println(err)
					os.Exit(0)
				}
				// fmt.Println(buf)
				dump.NetEventSendDecr(0)
			}
		}(uint64(instanceid))
	}
}
