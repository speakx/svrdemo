package app

import (
	"environment/cfgargs"
	"environment/dump"
	"environment/logger"
	"mmapcache/cache"
	"os"
	"svrdemo/client"
	"svrdemo/database"
	"svrdemo/proto/pbsvrdemo"
	"sync"

	"github.com/golang/protobuf/proto"
)

var once sync.Once
var app *App

// GetApp 获取当前服务的App实例
func GetApp() *App {
	once.Do(func() {
		app = &App{}
	})
	return app
}

// App 当前服务的App实例，用来存储一些运行时对象
type App struct {
	SrvCfg  *cfgargs.SrvConfig
	DB      *database.DB
	SrvDemo client.SimpleGrpcClient
	Single  client.SingleGrpcClient
}

// InitApp 加载配置、初始化日志、构建mmap缓存池
func (a *App) InitApp(srvCfg *cfgargs.SrvConfig) {
	a.SrvCfg = srvCfg

	// 初始化日志
	logger.Info("start init log")
	logger.InitLogger(srvCfg.Log.Path, srvCfg.Log.Console, srvCfg.Log.Level)
	logger.Info("end init log")

	// 初始化dump包，用来做服务端健康度检查&汇报
	logger.Info("start init dump, addr:", srvCfg.Dump.Addr)
	dump.InitDump(true, srvCfg.Dump.Interval, srvCfg.Dump.Addr,
		func(packetRecv, packetSend, packetRecvHandleRate, packetSendHandleRate int64) {
			logger.Info("dump rate recv:", packetRecvHandleRate, " send:", packetSendHandleRate)
		})
	logger.Info("end init dump")

	// 初始化mmap缓存
	logger.Info("start init mmap cache, dir:", srvCfg.Cache.Path)
	cache.InitMMapCachePool(
		srvCfg.Cache.Path, srvCfg.Cache.MMapSize,
		srvCfg.Cache.DataSize, srvCfg.Cache.PreAlloc,
		a.errorMMapCache, a.reloadMMapCache)
	logger.Info("end init mmap cache")

	// 初始化后端服务连接
	logger.Info("start init client")
	a.initClientSrv()
	logger.Info("end init client")

	// 初始化DB层
	logger.Info("start init db")
	if err := a.initDB(); nil != err {
		logger.Error("init db err:", err)
		os.Exit(0)
	}
	logger.Info("end init db")
}

func (a *App) errorMMapCache(err error) {
	logger.Error("mmapcache err:", err)
}

func (a *App) reloadMMapCache(mmapCaches []*cache.MMapCache) {
	logger.Info("reload mmapcache.count:", len(mmapCaches))
	for idx, mmapCache := range mmapCaches {
		logger.Info("reload mmapcache.idx:", idx, " data.count:%v", len(mmapCache.GetMMapDatas()))
		for _, mmapData := range mmapCache.GetMMapDatas() {
			var req pbsvrdemo.SimpleHello
			proto.Unmarshal(mmapData.GetData(), &req)
		}
		cache.DefPoolMMapCache.Collect(mmapCache)
	}
}

func (a *App) initClientSrv() {
	// a.SrvDemo.Connect(a.SrvCfg.Addr)
	a.SrvDemo.Connect("10.211.55.27:10000")
	a.Single.Connect("10.211.55.27:11000")
}

func (a *App) initDB() error {
	db, err := database.NewDB(a.SrvCfg)
	a.DB = db
	return err
}
