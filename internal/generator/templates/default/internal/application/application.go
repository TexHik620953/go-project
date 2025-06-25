package application

import (
	"context"

	"@APPNAME@/internal/appconfig"
@APPIMPORT@
)

type Application struct {
	ctx context.Context
	cfg *appconfig.AppConfig
@SERVICEDECL@
}

func New(ctx context.Context, cfg *appconfig.AppConfig) (*Application, error) {
	app := &Application{
		ctx: ctx,
		cfg: cfg,
	}
	// Building services here
	
@BUILDSERVICES@

	return app, nil
}

func (app *Application) Start() error {
@STARTAPP@
	return nil
}

@APPFUNCS@