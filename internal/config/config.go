package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/k1v4/organization_management_service/pkg/database/postgres"
)

type Config struct {
	postgres.DBConfig

	RestServerPort              int    `env:"REST_SERVER_PORT" env-description:"rest server port" env-default:"8080"`
	OrgMembershipServiceAddress string `env:"ORG_MEMBERSHIP_SERVICE_ADDRESS" env-description:"address of management users service"`
	KeyCloakIssuer              string `env:"KEYCLOAK_ISSUER" env-description:"issuer of keycloak"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
