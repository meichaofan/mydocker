package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		}, cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		}, cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		}, cli.StringFlag{
			Name:  "v",
			Usage: "volume",
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

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		//fmt.Printf("%v", *resConf)

		tty := context.Bool("it")
		volume := context.String("v")

		Run(tty, cmdArray, resConf,volume)
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
		err := container.RunContainerInitProcess()
		if err != nil {
			logrus.Errorf("%v", err)
		}
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			logrus.Errorf("Missing image name")
		}
		imageName := context.Args().Get(0)
		commitContainer(imageName)
		return nil
	},
}