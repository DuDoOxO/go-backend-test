package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-backend-test/infra"
	"github.com/go-backend-test/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoVar struct {
	Connection        string `env:"MONGO_CONNECTION,required" mapstructure:"MONGO_CONNECTION"`
	Database          string `env:"MONGO_DATABASE,required"  mapstructure:"MONGO_DATABASE"`
	CollectionUser    string `env:"MONGO_COLLECTION_USER,required" mapstructure:"MONGO_COLLECTION_USER"`
	CollectionMessage string `env:"MONGO_COLLECTION_MESSAGE,required" mapstructure:"MONGO_COLLECTION_MESSAGE"`
}

type UserRepo struct {
	vars MongoVar
	db   *mongo.Client
}

func NewUserRepo() IUserRepo {
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

func (r *UserRepo) AddUser(req model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	user := r.db.Database(r.vars.Database).Collection(r.vars.CollectionUser)
	_, err := user.InsertOne(ctx, req)

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	if err != nil {
		return fmt.Errorf("AddUser error :%s", err.Error())
	}

	return nil
}

func (r *UserRepo) CheckUserExist(req model.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	user := r.db.Database(r.vars.Database).Collection(r.vars.CollectionUser)
	filter := bson.M{"_id": req.LineUserId}

	counts, err := user.CountDocuments(ctx, filter)
	if err != nil {
		er := fmt.Errorf("CheckUserExist error :%s", err.Error())
		return false, er
	}

	if counts > 0 {
		return true, nil
	}

	return false, nil

}

func (r *UserRepo) AddUserMessage(req model.LineMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	msg := r.db.Database(r.vars.Database).Collection(r.vars.CollectionMessage)
	_, err := msg.InsertOne(ctx, req)

	req.CreatedAt = time.Now()
	if err != nil {
		return fmt.Errorf("AddUserMessage error :%s", err.Error())
	}

	return nil
}

func (r *UserRepo) ListUserMessageByUserId(userLineId string) ([]model.LineMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	msg := r.db.Database(r.vars.Database).Collection(r.vars.CollectionMessage)

	result := make([]model.LineMessage, 0)

	filter := bson.M{"user_id": userLineId}

	cursor, err := msg.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ListUserMessageByUserId find error :%s", err.Error())
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("ListUserMessageByUserId decodeerror :%s", err.Error())
	}

	return result, nil

}
