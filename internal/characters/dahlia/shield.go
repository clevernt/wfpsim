package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) genShield() {
	dur := 12 * 60
	if c.Base.Cons >= 4 {
		dur = 15 * 60
	}
	c.Core.Player.Shields.Add(&shield.Tmpl{
		ActorIndex: c.Index,
		Target:     -1,
		Src:        c.Core.F,
		ShieldType: shield.DahliaShield,
		Name:       "Dahlia Shield",
		HP:         c.shieldHP(),
		Ele:        attributes.Hydro,
		Expires:    c.Core.F + dur,
	})
}

func (c *char) hasShield() bool {
	return c.Core.Player.Shields.Get(shield.DahliaShield) != nil
}

func (c *char) shieldHP() float64 {
	return shieldHp[c.TalentLvlSkill()]*c.MaxHP() + shieldFlat[c.TalentLvlSkill()]
}
