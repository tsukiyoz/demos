package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoUpsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
		},
	}

	opts := options.Client().ApplyURI("mongodb://root:for.nothing@localhost:27017/?connect=direct").SetMonitor(monitor)

	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	collection := mdb.Collection("articles")

	defer func() {
		_, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Logf("delete error: %v\n", err)
			return
		}
	}()

	var atcl Article

	//if res, err := collection.InsertOne(ctx, &Article{
	//	Id:      1,
	//	Title:   "my title",
	//	Content: "my content",
	//}); err != nil {
	//	t.Error(err)
	//} else {
	//	t.Logf("[INSERT] id: %v\n", res.InsertedID)
	//}

	if res, err := collection.UpdateOne(ctx, bson.M{"id": 1}, bson.M{
		"$set": bson.M{
			"title":   "my new title",
			"content": "my new content",
		},
	}, options.Update().SetUpsert(true)); err != nil {
		t.Error(err)
	} else {
		t.Logf("[Upsert] result: %v\n", res)
	}

	if err := collection.FindOne(ctx, bson.M{"id": 1}).Decode(&atcl); err != nil {
		t.Error(err)
	} else {
		t.Logf("[RESULT] article: %v", atcl)
	}
}

func TestMongoTransaction(t *testing.T) {
	ctx := context.Background()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
		},
	}

	opts := options.Client().ApplyURI("mongodb://root:for.nothing@localhost:27017/?connect=direct").SetMonitor(monitor)

	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	collection := mdb.Collection("articles")

	defer func() {
		_, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Logf("delete error: %v\n", err)
			return
		}
	}()

	_, err = collection.InsertOne(ctx, &Article{
		Id:      1,
		Title:   "title",
		Content: "content",
	})
	if err != nil {
		t.Error(err)
		return
	}
	session, err := client.StartSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(tx mongo.SessionContext) (interface{}, error) {
		_, err := collection.UpdateOne(tx, bson.M{"id": 1}, bson.M{"$set": bson.M{
			"ctime": 123,
		}})
		if err != nil {
			return nil, err
		}

		// return nil, errors.New("mock error")

		_, err = collection.UpdateOne(tx, bson.M{"id": 1}, bson.M{"$set": bson.M{
			"title":   "new title",
			"content": "new content",
			"utime":   456,
		}})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		t.Log(err)
	}

	var atcl Article
	err = collection.FindOne(ctx, bson.M{"id": 1}).Decode(&atcl)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("[RESULT] %v\n", atcl)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
