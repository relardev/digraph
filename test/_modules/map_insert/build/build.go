package build

import (
	"example.com/project/internal/configuration"
	"example.com/project/internal/httpapi/another"
	"example.com/project/internal/httpapi/base"
	"example.com/project/internal/httpapi/common"
)

func Build(cfg configuration.Configuration) base.Processor {
	mappings := make(map[int]int)
	mappings[1] = 9
	mappings[2] = 8
	mappings[3] = 7

	translators := map[string]string{}
	translators["555"] = "123"
	translators["666"] = "12345"

	endpooints := map[string]common.Endpoint{
		"/123": common.NewCommon(),
	}

	endpooints["/"] = common.New()
	endpooints["/health"] = another.New()

	return base.New(endpooints, translators, mappings)
}
