package main

import (
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

	/*
		// Target dir
		targetDir := "."

		fmt.Print("Target directory name(default \".\", leave empty to proceed): ")
		fmt.Scanf("%s", &targetDir)
		if targetDir == "" {
			generationConfig.TargetDir = "."
		} else {
			generationConfig.TargetDir = targetDir
		}
		// Package name
		fmt.Print("Package name: ")
		fmt.Scanf("%s", &generationConfig.PackageName)
	*/
	genCtx, err := generator.NewGenerationContext(generationConfig)
	if err != nil {
		log.Fatal(err)
	}

	genCtx.AddEcho(generator.EchoGenParams{
		UseRecover: true,
		UseCORS:    false,
		UseLogger:  true,
	})

	genCtx.AddNats(generator.NatsParams{
		UseNatsRPC: true,
	})

	err = genCtx.GenerateProject()
	if err != nil {
		log.Fatal(err)
	}

}
