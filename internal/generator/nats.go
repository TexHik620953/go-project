package generator

type NatsParams struct {
	UseNatsRPC bool
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
		}`)
	}

}
