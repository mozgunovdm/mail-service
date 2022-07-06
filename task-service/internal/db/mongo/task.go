package mongo

import (
	"context"
	"fmt"
	"mts/task-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mongo) ListTask() (models.Tasks, error) {

	tasks := models.Tasks{}

	filter := bson.D{{}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	for cursor.Next(ctx) {
		var tsk models.Task
		if err := cursor.Decode(&tsk); err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, tsk)
	}

	return tasks, nil
}

func (m *Mongo) GetTask(id string) (*models.Task, error) {

	task := &models.Task{}
	filter := bson.M{"_id": id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := m.Collection.FindOne(ctx, filter)
	if err := res.Decode(&task); err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Mongo) WriteTask(tsk models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Collection.InsertOne(ctx, tsk)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) DeleteTask(id, login string) error {
	filter := bson.D{
		{"$and", bson.A{
			bson.M{"authorID": login},
			bson.M{"_id": id},
		}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) UpdateTask(tsk models.Task) error {
	filter := bson.M{
		"_id": tsk.Id,
	}
	update := bson.D{
		{"$set", bson.D{{"title", tsk.Title}}},
		{"$set", bson.D{{"description", tsk.Description}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) Close(ctx context.Context) {
	fmt.Println("Closing MongoDB Connection")
	if err := m.Client.Disconnect(ctx); err != nil {
		//log.Fatal().Err(errors.Wrapf(err, "Error on disconnection with MongoDB"))
		fmt.Println("Error on disconnection with MongoDB")
	}
}
