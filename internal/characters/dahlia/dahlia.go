package dahlia

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	maxBensionStacks int
	bensionStacks    int
}

func init() {
	core.RegisterCharFunc(keys.Dahlia, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.maxBensionStacks = 4
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.EnergyMax = 60
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1()
	return nil
}

func (c *char) addBensionStack(amt int) {
	if c.bensionStacks < c.maxBensionStacks {
		c.bensionStacks += amt
		c.Core.Log.NewEvent("add bension stacks", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.bensionStacks).
			Write("maxstacks", c.maxBensionStacks)
	}
	c.SetTag("bension-stacks", c.bensionStacks)

	if c.Base.Cons >= 1 {
		c.AddEnergy("c1", 2.5)
	}
}
