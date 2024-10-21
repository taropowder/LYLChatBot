package main

import (
	"LYLChatBot/conf"
	"LYLChatBot/handlers"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/database/postgresql"
	"LYLChatBot/pkg/redis_conn"
	"LYLChatBot/web"
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "LYL chat bot"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Required:    false,
			Value:       "config.yaml",
			Destination: &conf.ConfigFilePath,
		},

		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Required:    false,
			Destination: &conf.DebugMode,
		},
	}

	app.Before = func(c *cli.Context) error {

		_, err := os.Stat(conf.ConfigFilePath)
		if err != nil {
			out, err := yaml.Marshal(conf.NewDefaultConfig())
			if err == nil {
				err := ioutil.WriteFile(conf.ConfigFilePath, out, 0o644)
				if err != nil {
					log.Warnf("failed to write default config file to %s, err: %v", conf.ConfigFilePath, err)
				} else {
					log.Infof("write default config file to %s", conf.ConfigFilePath)
				}
			}
		} else {
			content, err := ioutil.ReadFile(conf.ConfigFilePath)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(content, &conf.ConfigureInstance)
			if err != nil {
				return err
			}
		}

		log.Infof("log level: %s", conf.ConfigureInstance.LogLevel)
		log.SetLevel(conf.ConfigureInstance.LogLevel)
		//conf.ConfigFilePath = configPath

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(c *cli.Context) error {

				database.SetDB(postgresql.InitDatabase(conf.ConfigureInstance.Database.Host,
					conf.ConfigureInstance.Database.Username,
					conf.ConfigureInstance.Database.Password,
					conf.ConfigureInstance.Database.Database,
					conf.ConfigureInstance.Database.Port))

				conn, err := redis_conn.NewRedisConnPool(conf.ConfigureInstance.Redis.Host, conf.ConfigureInstance.Redis.Password)
				if err != nil {
					log.Fatalf("error when init redis:  %v", err)
				}
				redis_conn.RedisConnPool = conn

				conf.BotInstance = openwechat.DefaultBot(openwechat.Desktop)
				dispatcher := openwechat.NewMessageMatchDispatcher()

				dispatcher.SetAsync(false)

				for _, h := range handlers.MessagesHandlers {
					dispatcher.RegisterHandler(h.Match, h.Handle)
				}

				// 注册消息回调函数
				conf.BotInstance.MessageHandler = dispatcher.AsMessageHandler()

				// 注册登陆二维码回调
				conf.BotInstance.UUIDCallback = openwechat.PrintlnQrcodeUrl

				// 执行热登录
				reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
				defer reloadStorage.Close()

				err = conf.BotInstance.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
				if err != nil {
					return err
				}

				go web.RunServer(conf.ConfigureInstance.APIAddress)

				// 阻塞主goroutine, 直到发生异常或者用户主动退出
				err = conf.BotInstance.Block()
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "debug",
			Action: func(c *cli.Context) error {

				database.SetDB(postgresql.InitDatabase(conf.ConfigureInstance.Database.Host,
					conf.ConfigureInstance.Database.Username,
					conf.ConfigureInstance.Database.Password,
					conf.ConfigureInstance.Database.Database,
					conf.ConfigureInstance.Database.Port))

				conn, err := redis_conn.NewRedisConnPool(conf.ConfigureInstance.Redis.Host, conf.ConfigureInstance.Redis.Password)
				if err != nil {
					log.Fatalf("error when init redis:  %v", err)
				}
				redis_conn.RedisConnPool = conn

				conf.BotInstance = openwechat.DefaultBot(openwechat.Desktop)
				dispatcher := openwechat.NewMessageMatchDispatcher()

				dispatcher.SetAsync(false)

				for _, h := range handlers.MessagesHandlers {
					dispatcher.RegisterHandler(h.Match, h.Handle)
				}

				// 注册消息回调函数
				conf.BotInstance.MessageHandler = dispatcher.AsMessageHandler()

				// 注册登陆二维码回调
				conf.BotInstance.UUIDCallback = openwechat.PrintlnQrcodeUrl

				// 执行热登录
				reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
				defer reloadStorage.Close()

				err = conf.BotInstance.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
				if err != nil {
					return err
				}

				go web.RunServer(conf.ConfigureInstance.APIAddress)

				// 阻塞主goroutine, 直到发生异常或者用户主动退出
				err = conf.BotInstance.Block()
				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
