package build

import (
	"example.com/project/internal/components"
	"example.com/project/internal/configuration"
)

type Base struct {
	comp2 components.Comp2
	comp3 components.Comp3
	comp4 components.Comp4
}

type Worker struct {
	Comp1 components.Comp1
	Base
}

func (b *Base) setComp2() {
	b.comp2 = components.NewComp2("two")
}

func (b *Base) setComp4() {
	b.comp4 = components.NewComp4(b.comp2)
}

func (b *Base) getComp3() components.Comp3 {
	return components.NewComp3("three")
}

func (w *Worker) Fill(cfg configuration.Configuration) {
	w.setComp2()
	w.setComp4()
	comp3 := w.getComp3()
	w.fillComp1(comp3)
}

func (w *Worker) fillComp1(comp3 components.Comp3) {
	w.Comp1 = components.NewComp1(w.comp2, comp3, w.comp4)
}
