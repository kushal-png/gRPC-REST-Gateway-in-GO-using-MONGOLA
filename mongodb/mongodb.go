package mongodb

import (
	"context"
	app "project/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error, context.Context, context.CancelFunc) {
	cs := app.MongoConnectionString

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)

	mg, err := mongo.Connect(ctx, options.Client().ApplyURI(cs))
	return mg, err, ctx, cancelFunc
}

func Close(client *mongo.Client, ctx context.Context, cancelFunc context.CancelFunc) {
	// CancelFunc to cancel to context
	defer cancelFunc()

	// client provides a method to close a mongoDB connection.
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func InsertOne(ctx context.Context, client *mongo.Client, db string, collection string, doc interface{}) (*mongo.InsertOneResult, error) {
	Collection := client.Database(db).Collection(collection)
	res, err := Collection.InsertOne(ctx, doc)
	return res, err
}

func FindOne(ctx context.Context, client *mongo.Client, db string, collection string, filter interface{}) *mongo.SingleResult {
	Collection := client.Database(db).Collection(collection)
	res := Collection.FindOne(ctx, filter)
	return res
}

func FindMany(ctx context.Context, client *mongo.Client, db string, collection string, filter interface{}) (*mongo.Cursor, error){
    Collection := client.Database(db).Collection(collection)
	cur, err:=Collection.Find(ctx, filter)
	return cur, err
}

func UpdateOne(ctx context.Context, client *mongo.Client, db string, collection string, filter interface{}, update interface{})(*mongo.UpdateResult, error){
	Collection := client.Database(db).Collection(collection)
	res, err:=Collection.UpdateOne(ctx, filter, update)
	return res, err
}

func DeleteOne(ctx context.Context, client *mongo.Client, db string, collection string, filter interface{})(*mongo.DeleteResult, error){
	Collection := client.Database(db).Collection(collection)
	res, err:=Collection.DeleteOne(ctx, filter)
	return res, err
}


