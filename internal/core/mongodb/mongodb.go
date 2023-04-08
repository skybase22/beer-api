package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoInstance
type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var MI MongoInstance

// Options mongo option
type Options struct {
	URL              string
	Port             int
	DatabaseName     string
	Username         string
	Password         string
	Debug            bool
	HandleNullValues []interface{}
}

var defaultNullValues = []interface{}{
	"",
	int(0),
}

// InitDatabase new database
func InitDatabase(o *Options) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%d", o.URL, o.Port)
	if o.Username != "" && o.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?connect=direct", o.Username, o.Password, o.URL, o.Port, o.DatabaseName)
	}
	clientOptions := options.Client().ApplyURI(uri).SetRegistry(buildNullValueDecoder(append(defaultNullValues, o.HandleNullValues)...))
	if o.Debug {
		clientOptions.Monitor = &event.CommandMonitor{
			Started: func(c context.Context, e *event.CommandStartedEvent) {
				fmt.Printf("\033[0;36mMongoDB command exec\033[0m: \033[1;95m%s\033[0m\n", e.Command.String())
			},
		}
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	MI = MongoInstance{
		Client: client,
		DB:     client.Database(o.DatabaseName),
	}

	return nil
}

func buildNullValueDecoder(val ...interface{}) *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()
	for _, v := range val {
		t := reflect.TypeOf(v)
		defDecoder, err := bson.DefaultRegistry.LookupDecoder(t)
		if err != nil {
			panic(err)
		}
		rb.RegisterTypeDecoder(t, &nullValueDecoder{defDecoder, reflect.Zero(t)})
	}
	return rb.Build()
}

type nullValueDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (d *nullValueDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bsontype.Null {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}
	if !val.CanSet() {
		return errors.New("value not settable")
	}
	if err := vr.ReadNull(); err != nil {
		return err
	}
	val.Set(d.zeroValue)
	return nil
}
