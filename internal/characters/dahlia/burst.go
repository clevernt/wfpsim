package dahlia

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 54
	burstBuffKey = "dahlia-favonian-favor"
)

func init() {
	burstFrames = frames.InitAbilSlice(61) // Q -> N1/Dash/Walk
	burstFrames[action.ActionSkill] = 51
	burstFrames[action.ActionJump] = 51
	burstFrames[action.ActionSwap] = 51
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.bensionStacks = 0
	if c.Base.Ascension >= 4 {
		for _, char := range c.Core.Player.Chars() {
			c.applyBurstBuff(char)
		}
	}
	c.genShield()
	dur := 12 * 60
	if c.Base.Cons >= 4 {
		dur += 3 * 60
	}
	c.AddStatus("dahlia-burst", dur, true)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Radiant Psalter",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), 0, 0)
	}, burstHitmark)
	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(6)

	naCaCount := 0
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if !(atk.Info.AttackTag == attacks.AttackTagNormal || atk.Info.AttackTag == attacks.AttackTagExtra) {
			return false
		}
		if c.bensionStacks == c.maxBensionStacks {
			return false
		}
		naCaCount++
		if naCaCount == 4 {
			c.addBensionStack(1)
			naCaCount = 0
		}
		return false
	}, fmt.Sprintf("bension-stack-%v", c.Base.Key.String()))
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) applyBurstBuff(char *character.CharWrapper) {
	favonianFavorDur := 12 * 60
	if c.Base.Cons >= 4 {
		favonianFavorDur += 3 * 60
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = min(c.MaxHP()/1000*0.005, 0.2)
	if c.Base.Cons >= 6 {
		m[attributes.AtkSpd] += 0.1
	}

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstBuffKey, favonianFavorDur),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
