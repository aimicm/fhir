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

func LocationIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

	var result []models.Location
	c := Database.C("locations")

	r.ParseForm()
	if len(r.Form) == 0 {
		iter := c.Find(nil).Limit(100).Iter()
		err := iter.All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else {
		searcher := search.NewMongoSearcher(Database)
		query := search.Query{Resource: "Location", Query: r.URL.RawQuery}
		err := searcher.CreateQuery(query).All(&result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}

	var locationEntryList []models.BundleEntryComponent
	for i := range result {
		var entry models.BundleEntryComponent
		entry.Resource = &result[i]
		locationEntryList = append(locationEntryList, entry)
	}

	var bundle models.Bundle
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Type = "searchset"
	var total = uint32(len(result))
	bundle.Total = &total
	bundle.Entry = locationEntryList

	log.Println("Setting location search context")
	context.Set(r, "Location", result)
	context.Set(r, "Resource", "Location")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(&bundle)
}

func LoadLocation(r *http.Request) (*models.Location, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("locations")
	result := models.Location{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting location read context")
	context.Set(r, "Location", result)
	context.Set(r, "Resource", "Location")
	return &result, nil
}

func LocationShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadLocation(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Location"))
}

func LocationCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	location := &models.Location{}
	err := decoder.Decode(location)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("locations")
	i := bson.NewObjectId()
	location.Id = i.Hex()
	err = c.Insert(location)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting location create context")
	context.Set(r, "Location", location)
	context.Set(r, "Resource", "Location")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Location/"+i.Hex())
	rw.WriteHeader(http.StatusCreated)
}

func LocationUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	location := &models.Location{}
	err := decoder.Decode(location)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("locations")
	location.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, location)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting location update context")
	context.Set(r, "Location", location)
	context.Set(r, "Resource", "Location")
	context.Set(r, "Action", "update")
}

func LocationDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("locations")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting location delete context")
	context.Set(r, "Location", id.Hex())
	context.Set(r, "Resource", "Location")
	context.Set(r, "Action", "delete")
}
