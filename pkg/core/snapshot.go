package core

import "log"

type SnapshotHandler interface {
	Release(*Snapshot) *Snapshot
	Get() *Snapshot
	Clone(*Snapshot) *Snapshot
}

type Snapshot struct {
	CharLvl  int
	ActorEle EleType
	// Actor      string //name of the character triggering the damage
	ActorIndex int
	ExtraIndex int  //this is currently purely for Kaeya icicle ICD
	Cancelled  bool //set to true if this snap should be ignored

	DamageSrc int //this is used for the purpose of calculating self harm
	SelfHarm  bool
	Targets   int //if TargetAll then not single target, resolve hitbox; other = target index

	SourceFrame     int
	AnimationFrames int //really only for amos bow...

	Abil        string      //name of ability triggering the damage
	WeaponClass WeaponClass //b.c. Gladiators...
	AttackTag   AttackTag
	ICDTag      ICDTag
	ICDGroup    ICDGroup
	ImpulseLvl  int

	HitWeakPoint bool
	Mult         float64 //ability multiplier. could set to 0 from initial Mona dmg
	StrikeType   StrikeType
	Element      EleType    //element of ability
	Durability   Durability //durability of aura, 0 if nothing applied

	UseDef  bool    //default false
	FlatDmg float64 //flat dmg; so far only zhongli

	Stats []float64 //total character stats including from artifact, bonuses, etc...

	BaseAtk float64 //base attack used in calc
	BaseDef float64 //base def used in calc
	//DmgBonus float64   //total damage bonus, including appropriate ele%, etc..
	RaidenDefAdj float64 //attack specific def shred (raiden c2)

	//reaction flags
	IsReactionDamage bool
	IsReaction       bool
	ReactionType     ReactionType
	IsMeltVape       bool
	ReactMult        float64 //reaction multiplier for melt/vape
	ReactBonus       float64 //reaction bonus %+ such as witch; should be 0 and only affected by hooks

	//callbacks
	OnHitCallback func(t Target)
}

var blankSnapshot = Snapshot{
	Stats: make([]float64, EndStatType),
}

func (s *Snapshot) Clear() {
	s.Copy(&blankSnapshot)
}

func (s *Snapshot) Copy(next *Snapshot) {
	s.CharLvl = next.CharLvl
	s.ActorEle = next.ActorEle
	s.ActorIndex = next.ActorIndex
	s.ExtraIndex = next.ExtraIndex
	s.Cancelled = next.Cancelled
	s.DamageSrc = next.DamageSrc
	s.SelfHarm = next.SelfHarm
	s.Targets = next.Targets
	s.SourceFrame = next.SourceFrame
	s.AnimationFrames = next.AnimationFrames
	s.Abil = next.Abil
	s.WeaponClass = next.WeaponClass
	s.AttackTag = next.AttackTag
	s.ICDTag = next.ICDTag
	s.ICDGroup = next.ICDGroup
	s.ImpulseLvl = next.ImpulseLvl
	s.HitWeakPoint = next.HitWeakPoint
	s.Mult = next.Mult
	s.StrikeType = next.StrikeType
	s.Element = next.Element
	s.Durability = next.Durability
	s.UseDef = next.UseDef
	s.FlatDmg = next.FlatDmg
	s.BaseAtk = next.BaseAtk
	s.BaseDef = next.BaseDef
	s.RaidenDefAdj = next.RaidenDefAdj
	s.IsReactionDamage = next.IsReactionDamage
	s.IsReaction = next.IsReaction
	s.ReactionType = next.ReactionType
	s.IsMeltVape = next.IsMeltVape
	s.ReactMult = next.ReactMult
	s.ReactBonus = next.ReactBonus
	s.OnHitCallback = next.OnHitCallback

	for i, v := range next.Stats {
		s.Stats[i] = v
	}
}

type Durability float64

// func (s *Snapshot) Clone() Snapshot {
// 	c := *s
// 	c.Stats = make([]float64, len(s.Stats))
// 	copy(c.Stats, s.Stats)
// 	return c
// }

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

const (
	TargetPlayer int = -2
	TargetAll    int = -1
)

type MemCtrl struct {
	free []*Snapshot
}

func NewMemCtrl(size int) *MemCtrl {
	c := &MemCtrl{}
	c.free = make([]*Snapshot, size, size+10)
	for i := 0; i < size; i++ {
		s := &Snapshot{
			Stats: make([]float64, EndStatType),
		}
		c.free[i] = s
	}
	return c
}

func (m *MemCtrl) Release(s *Snapshot) *Snapshot {
	//can be called safely on an already released snapshot
	s.Clear()
	m.free = append(m.free, s)
	return nil
}

func (m *MemCtrl) Get() *Snapshot {
	if len(m.free) == 0 {
		//append 10 more. capacity is + 10 so we should be ok
		for i := 0; i < 10; i++ {
			s := &Snapshot{
				Stats: make([]float64, EndStatType),
			}
			m.free = append(m.free, s)
		}
	}
	last := len(m.free) - 1
	s := m.free[last]
	m.free = m.free[:last]
	return s
}

func (m *MemCtrl) Clone(from *Snapshot) *Snapshot {
	if from == nil {
		log.Println("cloning from nil snapshot??")
	}
	next := m.Get()
	next.Copy(from)
	return next
}
