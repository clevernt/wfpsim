package dahlia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressCDStart = 16
	skillPressHitmark = 17
	skillPressTravel  = 1

	skillHoldCDStart = 11
	skillHoldHitmark = 12
	skillHoldTravel  = 3
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(39) // E -> N1/Q
	skillPressFrames[action.ActionDash] = 34
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionWalk] = 19
	skillPressFrames[action.ActionSwap] = 37

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(46) // E -> Swap
	skillHoldFrames[action.ActionAttack] = 38
	skillHoldFrames[action.ActionBurst] = 37
	skillHoldFrames[action.ActionDash] = 30
	skillHoldFrames[action.ActionJump] = 30
	skillHoldFrames[action.ActionWalk] = 30
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Immersive Ordinance",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Hydro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.4}, 2.5, 10),
		skillPressHitmark,
		skillPressHitmark+skillPressTravel,
		c.makeParticleCB(),
	)

	c.SetCDWithDelay(action.ActionSkill, 15*60, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionWalk], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Hydro, c.ParticleDelay)
	}
}
