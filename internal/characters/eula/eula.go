package eula

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("eula", NewChar)
}

type char struct {
	*character.Tmpl
	grimheartSrc    int
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a4()
	c.onExitField()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if c.Core.Status.Duration("eulaq") == 0 {
			return false
		}
		if ds.ActorIndex != c.Index {
			return false
		}
		if c.burstCounterICD > c.Core.F {
			return false
		}
		switch ds.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		case core.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.Core.Log.Debugw("eula burst add stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.Core.Log.Debugw("eula c6 add additional stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Core.F + 6
		return false
	}, "eula-burst-counter")
	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 29
		case 1:
			f = 25
		case 2:
			f = 65
		case 3:
			f = 33
		case 4:
			f = 88
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 35, 35 //TODO: no idea
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 34, 34
		}
		if c.Base.Cons >= 2 {
			return 34, 34 //press and hold have same cd
		}
		return 80, 80
	case core.ActionBurst:
		return 116, 116 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0, 0
	}
}

func (c *char) a4() {
	c.Core.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.Core.Log.Debugw("eula a4 reset skill cd", "frame", c.Core.F, "event", core.LogCharacterEvent)
		c.ResetActionCooldown(core.ActionSkill)

		return false
	}, "eula-a4")
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit")
}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.Abil != "Glacial Illumination (Lightfall)" {
			return false
		}
		if !c.Core.Flags.DamageMode {
			return false
		}
		if t.HP()/t.MaxHP() < 0.5 {
			ds.Stats[core.DmgP] += 0.25
			c.Core.Log.Debugw("eula: c4 adding dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "final dmgp", ds.Stats[core.DmgP])
		}
		return false
	}, "eula-c4")
}

func (e *char) Tick() {
	e.Tmpl.Tick()
	e.grimheartReset--
	if e.grimheartReset == 0 {
		e.Tags["grimheart"] = 0
	}
}
