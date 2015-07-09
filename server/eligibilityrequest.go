package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func EligibilityRequestIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.EligibilityRequest
	c := Database.C("eligibilityrequests")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var eligibilityrequestEntryList []models.EligibilityRequestBundleEntry
	for _, eligibilityrequest := range result {
		var entry models.EligibilityRequestBundleEntry
		entry.Id = eligibilityrequest.Id
		entry.Resource = eligibilityrequest
		eligibilityrequestEntryList = append(eligibilityrequestEntryList, entry)
	}

	var bundle models.EligibilityRequestBundle
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Type = "searchset"
	bundle.Total = len(result)
	bundle.Entry = eligibilityrequestEntryList

	log.Println("Setting eligibilityrequest search context")
	context.Set(r, "EligibilityRequest", result)
	context.Set(r, "Resource", "EligibilityRequest")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadEligibilityRequest(r *http.Request) (*models.EligibilityRequest, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("eligibilityrequests")
	result := models.EligibilityRequest{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting eligibilityrequest read context")
	context.Set(r, "EligibilityRequest", result)
	context.Set(r, "Resource", "EligibilityRequest")
	return &result, nil
}

func EligibilityRequestShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadEligibilityRequest(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "EligibilityRequest"))
}

func EligibilityRequestCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	eligibilityrequest := &models.EligibilityRequest{}
	err := decoder.Decode(eligibilityrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("eligibilityrequests")
	i := bson.NewObjectId()
	eligibilityrequest.Id = i.Hex()
	err = c.Insert(eligibilityrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting eligibilityrequest create context")
	context.Set(r, "EligibilityRequest", eligibilityrequest)
	context.Set(r, "Resource", "EligibilityRequest")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/EligibilityRequest/"+i.Hex())
}

func EligibilityRequestUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	eligibilityrequest := &models.EligibilityRequest{}
	err := decoder.Decode(eligibilityrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("eligibilityrequests")
	eligibilityrequest.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, eligibilityrequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting eligibilityrequest update context")
	context.Set(r, "EligibilityRequest", eligibilityrequest)
	context.Set(r, "Resource", "EligibilityRequest")
	context.Set(r, "Action", "update")
}

func EligibilityRequestDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("eligibilityrequests")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting eligibilityrequest delete context")
	context.Set(r, "EligibilityRequest", id.Hex())
	context.Set(r, "Resource", "EligibilityRequest")
	context.Set(r, "Action", "delete")
}