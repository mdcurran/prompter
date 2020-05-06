package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/mdcurran/prompter/internal/pkg/redis"
)

// Server exposes HTTP endpoints and a structured logger.
type Server struct {
	Router *httprouter.Router
	Logger *zap.SugaredLogger
}

// New instantiates an HTTP server + structured logger and builds a route table.
func New() (*Server, error) {
	l, err := logger()
	if err != nil {
		return nil, err
	}

	s := Server{
		Router: httprouter.New(),
		Logger: l,
	}
	s.buildRouteTable()

	return &s, nil
}

func logger() (*zap.SugaredLogger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}

func (s *Server) buildRouteTable() {
	s.Router.GET("/liveness", s.liveness)
	s.Router.POST("/save", s.save)
	s.Router.GET("/get/:pos", s.get)
	s.Router.POST("/user/add", s.addUser)
	s.Router.POST("/user/verify", s.verifyUser)
}

func encode(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}

func (s *Server) liveness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type response struct {
		Message string `json:"message"`
	}
	encode(w, http.StatusOK, response{Message: "Application healthy!"})
}

// vocabulary contains an array of tokens for a part of speech.
type vocabulary struct {
	PartOfSpeech string   `json:"pos"`
	Tokens       []string `json:"tokens"`
}

func (s *Server) save(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var v vocabulary
	err = json.Unmarshal(b, &v)
	if err != nil {
		Error(w, err, http.StatusUnprocessableEntity)
		return
	}

	if len(v.Tokens) == 0 {
		Error(w, errors.New("no tokens provided"), http.StatusInternalServerError)
		return
	}

	s.Logger.Infof("saving tokens %v with part of speech %q", v.Tokens, v.PartOfSpeech)
	err = redis.Save(v.PartOfSpeech, v.Tokens)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	encode(w, http.StatusCreated, v)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pos := ps.ByName("pos")
	if pos == "" {
		Error(w, errors.New("no part of speech provided"), http.StatusBadRequest)
		return
	}

	params := r.URL.Query()
	n, err := strconv.Atoi(params.Get("number"))
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	if n < 1 {
		Error(w, errors.New("provide a number of tokens to retrieve more than 0"), http.StatusBadRequest)
		return
	}

	// Returns n random tokens for the given part of speech.
	tokens, err := redis.Get(pos, int64(n))
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	v := vocabulary{PartOfSpeech: pos, Tokens: tokens}

	encode(w, http.StatusOK, v)
}

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var u user
	err = json.Unmarshal(b, &u)
	if err != nil {
		Error(w, err, http.StatusUnprocessableEntity)
		return
	}

	err = redis.AddUser(u.Email, u.Password)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	encode(w, http.StatusCreated, map[string]string{"email": u.Email})
}

func (s *Server) verifyUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type response struct {
		Verified bool `json:"verified"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode(w, http.StatusForbidden, response{Verified: false})
		return
	}

	var u user
	err = json.Unmarshal(b, &u)
	if err != nil {
		encode(w, http.StatusForbidden, response{Verified: false})
		return
	}

	ok := redis.VerifyUser(u.Email, u.Password)
	if !ok {
		encode(w, http.StatusForbidden, response{Verified: false})
		return
	}

	encode(w, http.StatusOK, response{Verified: true})
}
