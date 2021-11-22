package reactable

import "github.com/genshinsim/gcsim/pkg/core"

type Reactable struct {
	Durability [core.ElementMaxCount]core.Durability
	DecayRate  [core.ElementMaxCount]core.Durability
}

func (r *Reactable) Decay(frames int) {
	//declay electro/pyro/hydro/cryo
	r.Durability[core.Electro] -= r.DecayRate[core.Electro] * core.Durability(frames)
	r.Durability[core.Pyro] -= r.DecayRate[core.Pyro] * core.Durability(frames)
	r.Durability[core.Hydro] -= r.DecayRate[core.Hydro] * core.Durability(frames)
	r.Durability[core.Cryo] -= r.DecayRate[core.Cryo] * core.Durability(frames)
	//frozen is special
}
