package repo

import (
	"context"
	"log"
	"time"

	"github.com/go-backend-test/infra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoVar struct {
	Connection     string `env:"MONGO_CONNECTION,required" mapstructure:"MONGO_CONNECTION"`
	Database       string `env:"MONGO_DATABASE,required"  mapstructure:"MONGO_DATABASE"`
	CollectionUser string `env:"MONGO_COLLECTION_USER,required" mapstructure:"MONGO_COLLECTION_USER"`
}

type UserRepo struct {
	vars MongoVar
	db   *mongo.Client
}

func NewUserRepo() IUser {
	user := UserRepo{}
	initEnv := func() {
		if err := infra.GetEnv(&user.vars); err != nil {
			log.Fatal(err)
		}
	}

	initEnv()

	initDB := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(user.vars.Connection))
		if err != nil {
			log.Fatal("connect error :", err)
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal("ping error :", err)
		}

		user.db = client
	}
	initDB()

	return &user
}
