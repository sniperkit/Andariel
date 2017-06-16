package task

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoQueueEngine struct {
	Sess	*mgo.Session
}

func (this *MgoQueueEngine) FetchTasks(n int) ([]Task, error) {
	var ta []Task

	c := this.Sess.DB(MDbName).C(MDColl)
	err := c.Find(bson.M{"status": 1}).Limit(n).All(&ta)

	if err != nil {
		return ta, err
	}

	return ta, nil
}

func (this *MgoQueueEngine) DeleteTask(id interface{}) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.RemoveId(id)
}

func (this *MgoQueueEngine) Activate(id interface{}, status int16) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.UpdateId(id, bson.M{"$set": bson.M{"status": status}})
}
