package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"os"
)

// RestoreOptions is the collection of all the restore options as defined by the user
type RestoreOptions struct {
	BaseDbSrv               string            `json:"baseDbSrv"`
	BaseDbName              string            `json:"baseDbName"`
	DestDbSrv               string            `json:"destDbSrv"`
	DestDbName              string            `json:"destDbName"`
	CollectionsIgnore       []string          `json:"collectionsIgnore"`
	Accounts                []string          `json:"accounts"`
	UserCenteredCollections []CollectionField `json:"userCenteredCollections"`
}

// CollectionField is the type to define a field for documents in a specific collection
type CollectionField struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getConfigPath() (string, error) {
	path := os.Getenv("CONF_FILE_PATH")
	if path == "" {
		return path, errors.New("CONF_FILE_PATH environment variable not set")
	}
	return path, nil
}

func readConfig() RestoreOptions {
	path, err := getConfigPath()
	check(err)

	data, err := ioutil.ReadFile(path)
	check(err)
	var config RestoreOptions
	err = json.Unmarshal(data, &config)
	check(err)
	fmt.Println("config", config)
	return config
}

type copyManagerOpts struct {
	max      int
	current  *[]string
	rest     *[]string
	options  RestoreOptions
	base     MongoCtxDb
	dest     MongoCtxDb
	ll       LiveLogger
	copyDone chan string
}

func copyManager(opts copyManagerOpts) {
	for i := 0; i < opts.max-len(*opts.current); i++ {
		if len(*opts.rest) > 0 {
			*opts.current = append(*opts.current, (*opts.rest)[0])
			go CopyCollection(opts.base,
				opts.dest,
				(*opts.rest)[0],
				bson.D{},
				opts.ll,
				opts.copyDone)
			*opts.rest = (*opts.rest)[1:]
		}
	}
}

func main() {
	options := readConfig()
	fmt.Println("Connecting to DBs")
	base, dest := ConnectDbs(options)
	defer base.end()
	defer dest.end()
	fmt.Println("Connected")

	ll := NewLiveLogger()
	defer ll.End()

	fmt.Println("Getting Collections")
	collections := StringsRemoveElements(base.GetCollections(), options.CollectionsIgnore...)
	destCollections := dest.GetCollections()
	fmt.Println("Dropping Collections")
	dest.DropCollections(destCollections, options.CollectionsIgnore, ll)
	fmt.Println("Collections Dropped")
	current := []string{}
	rest := make([]string, len(collections))
	copy(rest, collections)

	copyDone := make(chan string, len(collections))
	// since we are using pointers for the dynamic parts, the options dont change from invocatioon
	copyManOpts := copyManagerOpts{
		max:      15,
		current:  &current,
		rest:     &rest,
		options:  options,
		base:     base,
		dest:     dest,
		ll:       ll,
		copyDone: copyDone,
	}
	copyManager(copyManOpts)
	for done := range copyDone {
		current = StringsRemoveElements(current, done)
		if len(rest) == 0 && len(current) == 0 {
			close(copyDone) // TODO make better flow control, receiver should not be the one closing the channel
			break
		}
		copyManager(copyManOpts)
	}
}
