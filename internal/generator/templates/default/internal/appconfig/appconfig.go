package appconfig

import "github.com/ilyakaznacheev/cleanenv"

type AppConfig struct {
	AppName string `env:"APP_NAME" env-default:"@APPNAME@"`
@CONFIGFIELDS@
}


func LoadAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
