package build

import (
	"strings"

	"example.com/project/internal/configuration"
	"example.com/project/internal/feature/getstate"
	"example.com/project/internal/feature/upload"
	"example.com/project/internal/httpapi/repo"
	pkgmongo "example.com/project/pkg/mongo"
	"github.com/gin-gonic/gin"
)

type Set struct {
	GETs  map[string]gin.HandlerFunc
	POSTs map[string]gin.HandlerFunc
}

func Build(cfg configuration.Configuration) (*Set, error) {
	connString := strings.ReplaceAll(
		cfg.Mongo.ConnectionString,
		"{{PASSWORD}}",
		cfg.Mongo.Password,
	)

	mongoTimeout := cfg.GetMongoTimeout()

	client, err := pkgmongo.NewClient(connString, mongoTimeout)
	if err != nil {
		return nil, err
	}

	repo := repo.NewRepo(client)

	gethandler := getstate.NewHandler(repo)

	return &Set{
		GETs: map[string]gin.HandlerFunc{
			"endpoint1": gethandler.Handle,
		},
		POSTs: map[string]gin.HandlerFunc{
			"endpoint2": upload.NewHandler(cfg.Test.DDD).WithClient(cfg.Mongo.CCCC).Handle,
		},
	}, nil
}
