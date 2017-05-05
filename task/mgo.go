package task

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoQueueEngine struct {
	Sess	*mgo.Session
}

func (this *MgoQueueEngine) FetchTasks(n uint32) ([]Task, error) {
	var ta []Task

	c := this.Sess.DB(MDbName).C(MDColl)
	err := c.Find(bson.M{}).Limit(n).All(ta)

	if err != nil {
		return ta, err
	}

	return ta, nil
}

func (this *MgoQueueEngine) DelTask(id interface{}) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.RemoveId(id)
}
