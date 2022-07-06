package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"mts/task-service/internal/db"
	"mts/task-service/internal/models"

	"golang.org/x/net/context"
)

type FileDatabase struct {
	File string
}

func NewFileRepository(ctx context.Context) (db.IDatabase, error) {
	file := "./task.json"
	_, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &FileDatabase{
		File: file,
	}, nil
}

func (db *FileDatabase) ListTask() (models.Tasks, error) {

	tasks := models.Tasks{}

	return tasks, nil
}

func (db *FileDatabase) GetTask(id string) (*models.Task, error) {

	data, err := ioutil.ReadFile(db.File)

	if err != nil {
		return nil, err
	}

	var tasks models.Tasks
	err = json.Unmarshal(data, &tasks.Tasks)

	if err != nil {
		return nil, err
	}

	if id != "" {
		for _, v := range tasks.Tasks {
			return &v, nil
		}
	}

	return nil, errors.New("Not found")
}

func (db *FileDatabase) WriteTask(tsk models.Task) error {
	data, err := ioutil.ReadFile(db.File)
	if err != nil {
		return err
	}

	var tasks []models.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}

	tasks = append(tasks, tsk)

	rankingsJSONgo, _ := json.Marshal(tasks)
	err = ioutil.WriteFile("task.json", rankingsJSONgo, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *FileDatabase) DeleteTask(id, login string) error {
	data, err := ioutil.ReadFile(db.File)

	if err != nil {
		return err
	}

	var tasks []models.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}

	var filteredTasks models.Tasks
	for _, v := range tasks {
		if v.Id != id {
			filteredTasks.Tasks = append(filteredTasks.Tasks, v)
		} else if v.AuthorId != login {
			filteredTasks.Tasks = append(filteredTasks.Tasks, v)
		}
	}

	rankingsJSONgo, _ := json.Marshal(filteredTasks.Tasks)
	err = ioutil.WriteFile(db.File, rankingsJSONgo, 0644)
	return err
}

func (db *FileDatabase) UpdateTask(tsk models.Task) error {
	data, err := ioutil.ReadFile(db.File)
	if err != nil {
		return err
	}

	var tasks []models.Task
	var newTasks []models.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}

	var newtask models.Task
	for _, v := range tasks {
		if v.Id == tsk.Id && v.AuthorId == tsk.AuthorId {
			newtask = models.Task{
				Id:          v.Id,
				Title:       tsk.Title,
				Description: tsk.Description,
				MailList:    tsk.MailList,
				AuthorId:    tsk.AuthorId,
				Status:      tsk.Status,
				Response:    models.Response{},
			}
			newTasks = append(newTasks, newtask)
		} else {
			newTasks = append(newTasks, v)
		}
	}

	rankingsJSONgo, _ := json.Marshal(newTasks)
	err = ioutil.WriteFile(db.File, rankingsJSONgo, 0644)
	return err
}

func (db *FileDatabase) Close(context.Context) {

}
