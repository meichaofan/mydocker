package main

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/Sirupsen/logrus"
	"mydocker/container"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
           mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	/**
	这里是run命令执行的真正函数
	1. 判断参数是否包含command
	2. 获取用户指定的command
	3. 调用 Run function 去准备启动容器
	 */
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("it")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container, Do not call it outside",
	/**
	1.获取传递过来的command参数
	2.执行容器初始化操作
	 */
	Action: func(context *cli.Context) error {
		logrus.Infof("init come on")
		cmd := context.Args().Get(0)
		logrus.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
