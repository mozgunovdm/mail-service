package router

import (
	"encoding/json"
	"math/rand"
	"mts/task-service/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (rs *Router) taskHandlers() http.Handler {
	h := chi.NewMux()

	h.Route("/", func(r chi.Router) {
		h.Post("/task", rs.CreteTask)
		h.Get("/task", rs.GetTask)
		h.Get("/tasks", rs.ListTask)
		h.Delete("/task/{id}", rs.DeleteTask)
		h.Patch("/task", rs.ChangeTask)
	})

	return h
}

func (rs *Router) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (rs *Router) CreteTask(w http.ResponseWriter, r *http.Request) {
	var body *models.Body
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&body)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(999999-100000) + 100000

	newTask := models.Task{
		Id:          strconv.Itoa(id),
		Title:       body.Title,
		Description: body.Description,
		MailList:    body.MailList,
		AuthorId:    body.AuthorId,
		Status:      "send",
		Response:    models.Response{},
	}

	err = rs.Db.WriteTask(newTask)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	render.Render(w, r, &newTask)
}

func (rs *Router) GetTask(w http.ResponseWriter, r *http.Request) {
	var body *models.Body
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&body)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	var tasks models.Tasks
	if body.Id != "" {
		task, err := rs.Db.GetTask(body.Id)
		if err != nil {
			rs.Log.Error().Err(err)
			return
		}
		render.Render(w, r, task)
		return
	}

	render.Render(w, r, &tasks)

}

func (rs *Router) ListTask(w http.ResponseWriter, r *http.Request) {
	var tasks models.Tasks
	tasks, err := rs.Db.ListTask()
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}
	render.Render(w, r, &tasks)
}

func (rs *Router) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	login, err := rs.Validate(w, r)
	if err != nil {
		return
	}

	err = rs.Db.DeleteTask(id, login)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	task := models.Task{
		Id: id,
	}

	w.WriteHeader(http.StatusAccepted)
	render.Render(w, r, &task)
}

func (rs *Router) ChangeTask(w http.ResponseWriter, r *http.Request) {
	var body *models.Body
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&body)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	login, err := rs.Validate(w, r)
	if err != nil {
		return
	}

	newtask := models.Task{
		Id:          body.Id,
		Title:       body.Title,
		Description: body.Description,
		MailList:    body.MailList,
		AuthorId:    login,
		Status:      "Send",
		Response:    models.Response{},
	}

	err = rs.Db.UpdateTask(newtask)
	if err != nil {
		rs.Log.Error().Err(err)
		return
	}

	render.Render(w, r, &newtask)
}
