package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		//提供run后面 -name 指定容器的名字
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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
		detach := context.Bool("d")
		if tty && detach {
			logrus.Errorf("tty and detach can not both provided")
		}

		// data binding
		volume := context.String("v")

		//将取到的容器名称传递下去，如果没有指定设置为空
		containerName := context.String("name")

		Run(tty, cmdArray, resConf, volume, containerName)
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

// mydocker ps
var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the container",
	Action: func(context *cli.Context) error {
		ListContainer()
		return nil
	},
}

// mydocker log
var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			logrus.Errorf("Missing container name")
		}
		containerName := ctx.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

// mydocker exec 容器 command
var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(ctx *cli.Context) error {
		//this is for callback
		if os.Getenv(ENV_EXEC_PID) != "" {
			logrus.Infof("pid callback pid %s", os.Getpid())
			return nil
		}

		if len(ctx.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}
		containerName := ctx.Args().Get(0)
		var commandArr []string
		for _, arg := range ctx.Args().Tail() {
			commandArr = append(commandArr, arg)
		}
		//执行命令
		ExecContainer(containerName, commandArr)
		return nil
	},
}
