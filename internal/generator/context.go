package generator

import (
	"fmt"
	"strings"

	"maps"

	"github.com/TexHik620953/go-project/internal/genconfig"
)

type TemplateReplaceFunc func(ctx *GenerationContext, template, fileName string) string

var defaultReplacements = map[string]TemplateReplaceFunc{
	"@APPNAME@": func(ctx *GenerationContext, template, fileName string) string {
		return ctx.config.PackageName
	},
	"@CONFIGFIELDS@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, field := range ctx.ConfigFields {
			envRequiredOrDefault := ""
			if field.Mandatory {
				envRequiredOrDefault = "env-required:\"\""
			} else {
				envRequiredOrDefault = fmt.Sprintf("env-default:\"%s\"", field.Default)
			}

			b.WriteString(fmt.Sprintf("\t%s %s `env:\"%s\" %s`\n", field.PascalName, field.Type, field.EnvName, envRequiredOrDefault))
		}
		return b.String()
	},
	"@LAUNCHENVS@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, field := range ctx.ConfigFields {
			b.WriteString(fmt.Sprintf("\"%s\": \"%s\",\n", field.EnvName, field.Default))
		}
		b.WriteString(fmt.Sprintf("\"APP_NAME\": \"dev-local-%s\"", ctx.config.PackageName))
		return b.String()
	},
	"@SERVICEDECL@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, v := range ctx.ServiceDeclarations {
			b.WriteString(fmt.Sprintf("\t%s\n", v))
		}
		return b.String()
	},
	"@BUILDSERVICES@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, v := range ctx.ServiceBuilders {
			b.WriteString(fmt.Sprintf("%s\n", v))
		}
		return b.String()
	},
	"@STARTAPP@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, v := range ctx.ServiceStartups {
			b.WriteString(v)
		}
		return b.String()
	},
	"@APPIMPORT@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, v := range ctx.AppImports {
			b.WriteString(fmt.Sprintf("\t\"%s\"\n", v))
		}
		return b.String()
	},
	"@APPFUNCS@": func(ctx *GenerationContext, template, fileName string) string {
		var b strings.Builder
		for _, v := range ctx.AppFuncs {
			b.WriteString(fmt.Sprintf("%s\n", v))
		}
		return b.String()
	},
}

type ConfigField struct {
	Type       string
	PascalName string

	Mandatory bool
	EnvName   string
	Default   string
}

type GenerationContext struct {
	config *genconfig.GenerationConfig

	templateFiles map[string][]byte // Path-content

	Installations        []string
	TemplateReplaceFuncs map[string]TemplateReplaceFunc
	ConfigFields         []ConfigField

	AppImports          []string
	AppFuncs            []string
	ServiceDeclarations []string
	ServiceBuilders     []string
	ServiceStartups     []string
}

func NewGenerationContext(config *genconfig.GenerationConfig) (*GenerationContext, error) {
	ctx := &GenerationContext{
		config: config,
		Installations: []string{
			"github.com/ilyakaznacheev/cleanenv",
			"github.com/google/uuid",
		},
		TemplateReplaceFuncs: map[string]TemplateReplaceFunc{},
		ConfigFields:         []ConfigField{},

		AppImports:          []string{},
		AppFuncs:            []string{},
		ServiceDeclarations: []string{},
		ServiceBuilders:     []string{},
		ServiceStartups:     []string{},
	}

	// Add defaull replacements
	maps.Copy(ctx.TemplateReplaceFuncs, defaultReplacements)

	var err error
	// Resolve template files
	ctx.templateFiles, err = getTemplate(config.TemplateName)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
