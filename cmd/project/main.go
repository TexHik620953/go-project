package main

import (
	"fmt"
	"log"

	"github.com/TexHik620953/go-project/internal/genconfig"
	"github.com/TexHik620953/go-project/internal/generator"
)

/*
1.

		Package Name:

	 2. What interfaces do you want to use
	    HTTP: Echo (github.com/labstack/echo/v4)
	    GRPC Server
	    NATS

3. Use database connector?
*/

func main() {
	generationConfig := &genconfig.GenerationConfig{
		TemplateName: "default",
		TargetDir:    ".",
		PackageName:  "sraka",
	}

	// Target directory
	fmt.Print("Target directory name(default '.', leave empty to proceed): ")
	var targetDir string
	fmt.Scanln(&targetDir) // Убрали форматную строку %s
	if targetDir == "" {
		generationConfig.TargetDir = "."
	} else {
		generationConfig.TargetDir = targetDir
	}

	// Package name
	fmt.Print("Package name: ")
	fmt.Scanln(&generationConfig.PackageName) // Убрали форматную строку %s

	// Use echo
	fmt.Print("Use echo?(y/n): ")
	var useEchoInput string
	fmt.Scanln(&useEchoInput)
	useEcho := useEchoInput == "y" || useEchoInput == "Y"

	// Use nats rpc
	fmt.Print("Use nats rpc?(y/n): ")
	var useNatsRpcInput string
	fmt.Scanln(&useNatsRpcInput)
	useNatsRpc := useNatsRpcInput == "y" || useNatsRpcInput == "Y"

	genCtx, err := generator.NewGenerationContext(generationConfig)
	if err != nil {
		log.Fatal(err)
	}

	if useEcho {
		genCtx.AddEcho(generator.EchoGenParams{
			UseRecover: true,
			UseCORS:    false,
			UseLogger:  true,
		})
	}

	if useNatsRpc {
		genCtx.AddNats(generator.NatsParams{
			UseNatsRPC:    true,
			UseNatsEvents: true,
		})
	}

	err = genCtx.GenerateProject()
	if err != nil {
		log.Fatal(err)
	}

}
