package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"to-do/handler"
)

type Api struct {
	Th *handler.TodoHandler
}

type dataResponse struct {
	Status  bool
	Message string
	Data    any
}

type baseResponses struct {
	Status  bool
	Message string
}

func (a *Api) Create(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Desc string `json:"description"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	res, err := a.Th.Create(params.Desc)

	if te, ok := errors.AsType[*handler.ToDoError](err); ok {
		respondWithError(w, te.Code, te.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("something went worng: ", err.Error()))
		return
	}

	type idStruct struct{ Id int }
	respondWithJSON(w, http.StatusCreated, dataResponse{true, "Todo Created successfully", idStruct{res}})
}

func (a *Api) Update(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	status := req.URL.Query().Get("status")

	err := a.Th.Update(status, id)
	if te, ok := errors.AsType[*handler.ToDoError](err); ok {
		respondWithError(w, te.Code, te.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("something went worng: ", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusNoContent, baseResponses{true, "Todo Update successfully"})
}

func (a *Api) Get(w http.ResponseWriter, req *http.Request) {
	var status = req.URL.Query().Get("status")
	data, err := a.Th.Get(status)

	if te, ok := errors.AsType[*handler.ToDoError](err); ok {
		respondWithError(w, te.Code, te.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("something went worng: ", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, dataResponse{true, "Data Fetched", data})
}

func (a *Api) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	err := a.Th.Delete(id)

	if te, ok := errors.AsType[*handler.ToDoError](err); ok {
		respondWithError(w, te.Code, te.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("something went worng: ", err.Error()))
		return
	}

	respondWithJSON(w, 201, baseResponses{
		true,
		fmt.Sprintf("Todo id: %s deletead successfully", id),
	})
}
