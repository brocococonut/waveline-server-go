package router

import "go.mongodb.org/mongo-driver/mongo"

type (
	Router struct {
		Client *mongo.Client
	}
)
