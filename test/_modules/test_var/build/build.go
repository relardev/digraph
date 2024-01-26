package build

import (
	"example.com/project/internal/configuration"
	"example.com/project/internal/httpapi/another"
	"example.com/project/internal/httpapi/base"
	"example.com/project/internal/httpapi/common"
)

type Set struct {
	Comp1 *common.Component
	Comp2 *common.Component
	Comp3 *another.Component
}

func (s *Set) Build(cfg configuration.Configuration) {
	timeout := cfg.Timeout

	var1 := base.New(timeout)

	s.Comp2 = common.NewComponent(var1, "two")

	s.Comp3 = another.New(timeout, "two")

	s.fillCommonComponent(var1)
}

func (s *Set) fillCommonComponent(baseComponent base.Base) {
	s.Comp1 = common.NewComponent(baseComponent, "two")
}
