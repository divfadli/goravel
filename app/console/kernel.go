package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
)

type Kernel struct {
}

func (kernel Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Call(func() {
			// facades.Orm().Query().Where("1 = 1").Delete(&models.User{})
		}).Daily(),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{}
}
