package conf

import (
	"controllers"

	"github.com/mrkt/cellgo"
)

func SetEvent() {
	cellgo.RegisterEvent("event1", 1)
	cellgo.RegisterEvent("event2", 1)
	//Add Event's happen
	cellgo.Events["event2"].EventAdd("EventCreate", &controllers.EventCreateController{}, []string{"Run"}, int64(0), int64(0), 1)
}
