package storage

import (
	"fmt"

	"github.com/guemidiborhane/factorydigitale.tech/internal/config"
	"github.com/guemidiborhane/factorydigitale.tech/internal/utils"
	meilisearch "github.com/meilisearch/meilisearch-go"
)

type meilisearchConfig struct {
	Host      string `mapstructure:"MEILI_HOST"`
	Port      int    `mapstructure:"MEILI_PORT"`
	MasterKey string `mapstructure:"MEILI_MASTER_KEY"`
}

var (
	MeilisearchConfig = &meilisearchConfig{
		Host:      "localhost",
		Port:      7700,
		MasterKey: "",
	}

	MeiliClient meilisearch.ServiceManager
)

func SetupMeilisearch() {
	if err := config.EnvFile.LoadConfig(&MeilisearchConfig); err != nil {
		utils.WriteToStderr(err)
	}

	MeiliClient = NewMeiliClient()
}

func NewMeiliClient() meilisearch.ServiceManager {
	return meilisearch.New(fmt.Sprintf("http://%s:%d", MeilisearchConfig.Host, MeilisearchConfig.Port), meilisearch.WithAPIKey(MeilisearchConfig.MasterKey))
}
