package generator

import "strings"

type NatsParams struct {
	UseNatsRPC    bool
	UseNatsEvents bool
}

func (c *GenerationContext) AddNats(params NatsParams) {
	c.Installations = append(c.Installations, "github.com/nats-io/nats.go")
	c.AppImports = append(c.AppImports, "github.com/nats-io/nats.go", "fmt")
	c.ServiceDeclarations = append(c.ServiceDeclarations, "natsClient *nats.Conn")
	c.ConfigFields = append(c.ConfigFields, ConfigField{
		Type:       "string",
		PascalName: "NatsAddr",
		EnvName:    "NATS_ADDR",
		Default:    "nats://127.0.0.1:4222",
	})

	c.ServiceBuilders = append(c.ServiceBuilders, `	{
		// Nats client declaration
		var err error
		app.natsClient, err = nats.Connect(cfg.NatsAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to nats: %v", err)
		}	
	}`)

	if params.UseNatsRPC {
		c.Installations = append(c.Installations, "github.com/TexHik620953/natsrpc-go")
		c.AppImports = append(c.AppImports, "github.com/TexHik620953/natsrpc-go")
		c.ServiceDeclarations = append(c.ServiceDeclarations, "natsRpc natsrpc.NatsRPC")

		c.ServiceBuilders = append(c.ServiceBuilders, `	{
		// Nats rpc handler declaration
		app.natsRpc = natsrpc.New(app.natsClient, natsrpc.WithBaseName(app.cfg.AppName))
		app.natsRpcRouter()
	}`)
		c.ServiceStartups = append(c.ServiceStartups, `	{
		err := app.natsRpc.StartWithContext(context.Background())
		if err != nil {
			return err
		}
	}`)
		{
			ab := strings.Builder{}
			ab.WriteString("func (app *Application ) natsRpcRouter() {\n")
			ab.WriteString("\tapp.natsRpc.AddRPC(\"health\", func(c natsrpc.NatsRPCContext) error {\n")
			ab.WriteString("\t\treturn c.RespondJSON(map[string]string{\"status\": \"OK\"})\n")
			ab.WriteString("\t})\n")
			ab.WriteString("}")
			c.AppFuncs = append(c.AppFuncs, ab.String())
		}
	}

	if params.UseNatsEvents {
		c.Installations = append(c.Installations, "github.com/TexHik620953/natsevent-go")
		c.AppImports = append(c.AppImports, "github.com/TexHik620953/natsevent-go", "github.com/nats-io/nats.go/jetstream")

		c.ServiceDeclarations = append(c.ServiceDeclarations,
			"natsEvent natsevent.NatsEventService",
			"jetStream jetstream.JetStream",
		)

		c.ServiceBuilders = append(c.ServiceBuilders, `	{
		// Nats JetStream declaration
		var err error
		app.jetStream, err = jetstream.New(nc)
		if err != nil {
			return nil, fmt.Errorf("failed to connect jetstream: %v", err)
		}

		// Nats event handler declaration
		app.natsEvent = natsevent.New(app.jetStream, natsevent.WithSubjectRoot("events"))
		app.natsEventRouter()
	}`)

		{
			ab := strings.Builder{}
			ab.WriteString("func (app *Application ) natsEventRouter() {\n")

			ab.WriteString("}")
			c.AppFuncs = append(c.AppFuncs, ab.String())
		}

		c.ServiceStartups = append(c.ServiceStartups, `	{
		err := app.natsEvent.StartWithContext(app.ctx)
		if err != nil {
			return err
		}
	}`)
	}
}
