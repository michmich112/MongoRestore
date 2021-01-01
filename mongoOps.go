package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mOptions "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

// MongoCtxDb is the type defining a Db connection instance and the context that is used by that instance
type MongoCtxDb struct {
	ctx context.Context
	db  *mongo.Database
}

// ConnectDbs is used to start connections
func ConnectDbs(options RestoreOptions) (base MongoCtxDb, dest MongoCtxDb) {
	base = connect(options.BaseDbSrv, options.BaseDbName)
	dest = connect(options.DestDbSrv, options.DestDbName)
	return base, dest
}

func connect(srv string, db string) MongoCtxDb {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, mOptions.Client().ApplyURI(srv))
	check(err)
	clientDB := client.Database(db)
	return MongoCtxDb{ctx, clientDB}
}

// GetCollections for a specified database
func (m MongoCtxDb) GetCollections() []string {
	cols, err := m.db.ListCollectionNames(m.ctx, bson.D{})
	check(err)
	return cols
}

// DropCollections is the method do drop all the collections indicated
func (m MongoCtxDb) DropCollections(collections []string, ignore []string, ll LiveLogger) {
	cols := StringsRemoveElements(collections, ignore...)
	for _, collection := range cols {
		err := m.db.Collection(collection).Drop(m.ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[ DROP ] - Dropped %s \n", collection)
		//ll.msgChannel <- LiveLoggerMessage{
		//	source:    collection,
		//	eventType: "DROP",
		//	message:   "dropped",
		//}
	}
}

/**
 * end mongo ctx db connection
 * defer this after connecting to it
 */
func (m MongoCtxDb) end() {
	err := m.db.Client().Disconnect(m.ctx)
	check(err)
	m.ctx.Done()
}

// CopyCollection is to copy all the documents that match a filter of a
// Collection onto another
func CopyCollection(base MongoCtxDb, dest MongoCtxDb, collection string, filters interface{}, ll LiveLogger, cb chan string) {
	var wg sync.WaitGroup // create wait group for concurrency handling
	startTime := time.Now()
	bufferSize := 1000 // buffer before inserting new documents
	opts := mOptions.Find().SetBatchSize(1000)
	cursor, err := base.db.Collection(collection).Find(base.ctx, filters, opts)

	defer cursor.Close(base.ctx)

	if err != nil {
		fmt.Println(err)
		return
	}

	var count int = 0
	var results []interface{} // results slice
	for cursor.Next(base.ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
		count++
		if count%bufferSize == 0 {
			wg.Add(1)
			// create a new go routine to insert all the buffered documents
			go func(res []interface{}) {
				defer wg.Done()
				opts := mOptions.InsertMany().SetOrdered(false)
				dest.db.Collection(collection).InsertMany(dest.ctx, res, opts)
			}(results[:])
			results = results[len(results):] // reset the results
		}
	}
	if len(results) > 0 {
		dest.db.Collection(collection).InsertMany(dest.ctx, results)
	}
	wg.Wait()
	fmt.Printf("%s => DONE %.2f\n", collection, time.Since(startTime).Seconds())
	cb <- collection
}
