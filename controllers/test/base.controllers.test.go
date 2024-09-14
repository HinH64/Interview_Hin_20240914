package controllers

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB() {
	err := mgm.SetDefaultConfig(nil, "testdb", options.Client().ApplyURI("mongodb+srv://Admin:DjHfIiOYzX5WGV2Q@cluster0.gvmek.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		panic(err)
	}
}