package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	UserQuery struct {
		ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		Username string             `json:"username"`
		Password string             `json:"password"`
		Email    string             `json:"email"`
		Clients  []Client           `json:"clients"`
	}
	User struct {
		Username string             `json:"username"`
		Password string             `json:"password"`
		Email    string             `json:"email"`
		Clients  []Client           `json:"clients"`

	}
	Client struct {
		ClientName        string   `json:"client_name"`
		ClientDescription string   `json:"client_description"`
		ClientContacts    []string `json:"client_contacts"`
		// daily weekly monthly yearly
		ClientContactRate string `json:"client_contact_rate"`
		ClientAddress     string `json:"client_adress"`
		LastContacted primitive.DateTime `json:"last_contacted"`
		PreviousPurchases []Sales `json:"previous_purchases"`
	}
	Sales struct {
		DateSold primitive.DateTime `json:"date_sold"`
		ProductsBought []string `json:"products_bought"`
	}
	Status struct {
		OK         bool
		FixedCount int64
	}
)

func db() (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client.Database("email-agenda").Collection("users"), nil
}

func DeleteUser(ID string, Users *mongo.Collection) Status {
	objectId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		fmt.Println(err)
		return Status{OK: false, FixedCount: 0}
	}
	opts := options.Delete().SetCollation(&options.Collation{})
	res, err := Users.DeleteOne(context.TODO(), bson.M{"_id": objectId}, opts)
	if err != nil {
		log.Fatal(err)
		return Status{OK: false, FixedCount: 0}
	}
	return Status{OK: true, FixedCount: res.DeletedCount}
}

func CreateUser(Users *mongo.Collection, UserData User) interface{} {
	insertResult, err := Users.InsertOne(context.TODO(), UserData)
	if err != nil {
		log.Fatal(err)
	}
	ID := insertResult.InsertedID
	return ID
}

func FindUserByID(ID string, Users *mongo.Collection) *UserQuery {
	objectId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		fmt.Println(err)
	}
	result := &UserQuery{}
	err = Users.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func CheckUserExistenceByQuery(query bson.M, Users *mongo.Collection) (*UserQuery, error) {
	result := &UserQuery{}
	err := Users.FindOne(context.TODO(), query).Decode(result)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	return result, nil
}

func UpdateUserByID(query bson.M, updateFilter bson.M, Users *mongo.Collection) *UserQuery {
	// filter := bson.M{"_id": insertResult.InsertedID}
	// update := bson.M{"$set": bson.M{"username": "WASAAAAP"}}
	result := &UserQuery{}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := Users.FindOneAndUpdate(context.TODO(), query, updateFilter, &returnOpt).Decode(result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}
