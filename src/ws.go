package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kardianos/service"
)

var cmd *exec.Cmd

type program struct {
	Name      string   //应用程序名
	Args      []string //应用程序参数
	WorkSpace string
	Env       []string
}

func (p *program) Start(s service.Service) error {
	Info.Println("Start ...")
	go p.run()
	return nil
}
func (p *program) run() {
	cmd = exec.Command(p.Name, p.Args...)
	cmd.Dir = p.WorkSpace
	cmd.Env = append(os.Environ(), p.Env...)
	stdout, err := os.OpenFile(p.WorkSpace+"/out.log", os.O_CREATE|os.O_WRONLY, 0600)
	// Info.Fatal(cmd.Env)
	if err != nil {
		log.Fatalln(err)
	}
	defer stdout.Close()
	cmd.Stdout = stdout
	Info.Println("cd", p.WorkSpace)
	Info.Println(p.Name, strings.Join(p.Args, " "))
	err = cmd.Run()
	if err != nil {
		Error.Fatal("service application has error:", err)
	}
}
func (p *program) Stop(s service.Service) error {
	Info.Println("Stop ...", p.Name)
	cmd.Process.Kill()
	return nil
}
func usage() {
	fmt.Printf("Usage:\n%s [arguments] <command>\n", os.Args[0])
	fmt.Println("\nThe arguments are:")
	flag.PrintDefaults()
	fmt.Println("\nThe command should be:")
	fmt.Printf("  %s\t %s\n", "install", "install service")
	fmt.Printf("  %s\t %s\n", "remove", "remove service")
}
func main() {
	env := flag.String("e", "", "specify serivce application `env`, eg. 'key1=value1,key2=value2'")
	dir := flag.String("w", "", "specify service application  `workspace`")
	serviceName := flag.String("n", "", "service `name`")
	description := flag.String("d", "", "service `description`")
	executable := flag.String("p", "", "specify service application `path`")
	parameter := flag.String("a", "", "specify serivce application `args` using ','")
	help := flag.Bool("help", false, "for more help")
	flag.Parse()

	if *help || len(flag.Args()) == 0 {
		usage()
		return
	}
	workspace := *dir
	if workspace == "" {
		workspace, _ = os.Getwd()
	}
	cmd := flag.Args()[0]
	args := []string{
		"-n", *serviceName,
		"-d", *description,
		"-p", *executable,
		"-a", *parameter,
		"-e", *env,
		"-w", workspace,
		"run",
	}

	svcConfig := &service.Config{
		Name:             *serviceName, //服务名称
		DisplayName:      *serviceName, //显示名称
		Description:      *description, //服务描述
		Arguments:        args,
		WorkingDirectory: workspace,
	}
	prg := &program{
		Name:      *executable,
		Args:      strings.Split(*parameter, ","),
		WorkSpace: workspace,
		Env:       strings.Split(*env, ","),
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		usage()
		log.Fatal("create new service error: ", err.Error())
	}

	switch cmd {
	case "install":
		err := s.Install()
		if err != nil {
			usage()
			log.Fatal("install new service error: ", err.Error())
		}
		log.Println(*serviceName, " is installed")
		return
	case "remove":
		err := s.Uninstall()
		if err != nil {
			usage()
			log.Fatal("remove service ", serviceName, " error: ", err.Error())
		}
		log.Println(*serviceName + " is removed")
		return
	case "run":
		err := s.Run()
		if err != nil {
			usage()
			log.Fatal("run service error: ", err.Error())
		}
		return
	}
	usage()
}
