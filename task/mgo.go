package task

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoQueueEngine struct {
	Sess	*mgo.Session
}

func (this *MgoQueueEngine) FetchTask() (*Task, error) {
	var ta Task

	c := this.Sess.DB(MDbName).C(MDColl)
	err := c.Find(bson.M{}).One(&ta)

	if err != nil {
		return &ta, err
	}

	return &ta, nil
}

func (this *MgoQueueEngine) DelTask(id interface{}) error {
	c := this.Sess.DB(MDbName).C(MDColl)

	return c.RemoveId(id)
}
