//package GOCouchDBAPIs
package main

import (
	"github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"
	"context"
)

type CouchDBAccount struct {
	Id      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"` 
	Name    string `json:"name,omitempty"` 
	Deposit int32    `json:"deposit,omitempty"`    
}

func CreateDBs(client *kivik.Client,DBname string){
   client.CreateDB(context.TODO(),DBname)
}

func AddAccounts(num int,db *kivik.DB){
	for i:= 0; i < num; i++{
		Account := CouchDBAccount{Deposit: 100000000}
		id, rev, err := db.CreateDoc(context.TODO(), Account)
		if err != nil {
			panic(err)
		}
		Account.Rev = rev
		Account.Id = id
	 } 
}

func main(){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
	
	db := client.DB(context.TODO(), "bank3")
	num := 1000
	AddAccounts(num,db)
}