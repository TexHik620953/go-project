package generator

import "strings"

type EchoGenParams struct {
	UseRecover bool
	UseCORS    bool
	UseLogger  bool
}

func (c *GenerationContext) AddEcho(params EchoGenParams) {
	c.Installations = append(c.Installations,
		"github.com/labstack/echo/v4",
		"github.com/labstack/echo/v4/middleware",
	)
	c.AppImports = append(c.AppImports,
		"github.com/labstack/echo/v4",
		"github.com/labstack/echo/v4/middleware",
		"net/http",
		"time",
	)
	c.ServiceDeclarations = append(c.ServiceDeclarations, "echoSrv *echo.Echo")

	{
		// Service builder
		sb := strings.Builder{}
		sb.WriteString("\t{\n\t\t// Echo server declaration\n")
		sb.WriteString("\t\tapp.echoSrv = echo.New()\n")
		if params.UseRecover {
			sb.WriteString("\t\tapp.echoSrv.Use(middleware.Recover())\n")
		}

		if params.UseCORS {
			sb.WriteString("\t\tapp.echoSrv.Use(middleware.CORS())\n")
		}
		if params.UseLogger {
			sb.WriteString("\t\tapp.echoSrv.Use(middleware.Logger())\n\n")
		}

		sb.WriteString("\t\tapp.echoRouter()\n")

		sb.WriteString("\t}")
		c.ServiceBuilders = append(c.ServiceBuilders, sb.String())
	}

	// Add env config
	c.ConfigFields = append(c.ConfigFields, ConfigField{
		Type:       "string",
		PascalName: "HttpHostname",
		EnvName:    "HTTP_HOSTNAME",
		Default:    ":8080",
	})

	// Add startup
	c.ServiceStartups = append(c.ServiceStartups, `	go app.echoSrv.Start(app.cfg.HttpHostname)
	go func() {
		<-app.ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.echoSrv.Shutdown(shutdownCtx)
	}()`)

	// Add routing method
	{
		rb := strings.Builder{}
		rb.WriteString("func (app *Application) echoRouter() {\n")

		rb.WriteString("\tapp.echoSrv.Any(\"/health\", func(c echo.Context) error {\n")
		rb.WriteString("\t\treturn c.String(http.StatusOK, \"OK\")\n")
		rb.WriteString("\t})\n")
		rb.WriteString("}")

		c.AppFuncs = append(c.AppFuncs, rb.String())
	}
}
