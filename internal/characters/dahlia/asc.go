package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
)

const (
	a1IcdKey = "a1-icd-key"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnFrozen, func(args ...interface{}) bool {
		if !c.StatusIsActive("dahlia-burst") {
			return false
		}
		if c.StatusIsActive(a1IcdKey) {
			return false
		}
		c.AddStatus(a1IcdKey, 8*60, true)
		c.addBensionStack(2)
		return false
	}, "add-bension-stack-from-freeze")
}
