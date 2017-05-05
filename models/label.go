/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 05/05/2017        Jia Chenhui
 */

package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"Andariel/mongo"
)

type LabelServiceProvider struct {
}

var LabelService *LabelServiceProvider
var LabelCollection *mgo.Collection

func PrepareGitLabel() {
	LabelCollection = mongo.GithubSession.DB(mongo.MDGitName).C("Label")

	LabelIndex := mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	if err := LabelCollection.EnsureIndex(LabelIndex); err != nil {
		panic(err)
	}

	LabelService = &LabelServiceProvider{}
}

type Label struct {
	LabelID bson.ObjectId `bson:"LabelID,omitempty" json:"id"`
	Name    string        `bson:"Name" json:"name"`
	Desc    string        `bson:"Desc" json:"desc"`
	Active  bool          `bson:"Active" json:"active"`
	Total   uint64        `bson:"Total" json:"total"`
}

type Activate struct {
	Name   string
	Active bool
}

func (tsp *LabelServiceProvider) Create(Label *Label) error {
	l := Label{
		Name:   Label.Name,
		Active: Label.Active,
		Desc:   Label.Desc,
	}

	err := LabelCollection.Insert(&l)
	if err != nil {
		return err
	}

	return nil
}

func (tsp *LabelServiceProvider) ListAll() ([]Label, error) {
	var l []Label

	err := LabelCollection.Find(nil).All(&l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (tsp *LabelServiceProvider) Activate(activate Activate) error {
	update := bson.M{"$set": bson.M{
		"Active": activate.Active,
	}}

	err := LabelCollection.Update(bson.M{"Name": activate.Name}, &update)
	if err != nil {
		return err
	}

	return nil
}
