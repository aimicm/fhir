package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"github.com/intervention-engine/fhir/search"
	"gopkg.in/mgo.v2/bson"
)

func GroupIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if r := recover(); r != nil {
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			switch x := r.(type) {
			case search.SearchError:
				rw.WriteHeader(x.HTTPStatus())
				json.NewEncoder(rw).Encode(x.OperationOutcome())
				return
			default:
				e := search.InternalServerError(fmt.Sprintf("%s", x))
				rw.WriteHeader(e.HTTPStatus())
				json.NewEncoder(rw).Encode(e.OperationOutcome())
			}
		}
	}()

	var result []models.Group
	c := Database.C("groups")

	r.ParseForm()
	if len(r.Form) == 0 {
		iter := c.Find(nil).Limit(100).Iter()
		err := iter.All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else {
		searcher := search.NewMongoSearcher(Database)
		query := search.Query{Resource: "Group", Query: r.URL.RawQuery}
		err := searcher.CreateQuery(query).All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}

	var groupEntryList []models.BundleEntryComponent
	for i := range result {
		var entry models.BundleEntryComponent
		entry.Resource = &result[i]
		groupEntryList = append(groupEntryList, entry)
	}

	var bundle models.Bundle
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Type = "searchset"
	var total = uint32(len(result))
	bundle.Total = &total
	bundle.Entry = groupEntryList

	log.Println("Setting group search context")
	context.Set(r, "Group", result)
	context.Set(r, "Resource", "Group")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(&bundle)
}

func LoadGroup(r *http.Request) (*models.Group, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("groups")
	result := models.Group{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting group read context")
	context.Set(r, "Group", result)
	context.Set(r, "Resource", "Group")
	return &result, nil
}

func GroupShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadGroup(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Group"))
}

func GroupCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	group := &models.Group{}
	err := decoder.Decode(group)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("groups")
	i := bson.NewObjectId()
	group.Id = i.Hex()
	err = c.Insert(group)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting group create context")
	context.Set(r, "Group", group)
	context.Set(r, "Resource", "Group")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Group/"+i.Hex())
	rw.WriteHeader(http.StatusCreated)
}

func GroupUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	group := &models.Group{}
	err := decoder.Decode(group)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("groups")
	group.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, group)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting group update context")
	context.Set(r, "Group", group)
	context.Set(r, "Resource", "Group")
	context.Set(r, "Action", "update")
}

func GroupDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("groups")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting group delete context")
	context.Set(r, "Group", id.Hex())
	context.Set(r, "Resource", "Group")
	context.Set(r, "Action", "delete")
}
