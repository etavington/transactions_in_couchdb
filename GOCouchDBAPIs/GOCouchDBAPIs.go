package GOCouchDBAPIs

import (
	"github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"
	//"github.com/go-kivik/kivik/v4"
	"context"
	"math/rand"
	"fmt"
	//"log"
)

type CouchDBAccount struct {
	Id      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"` 
	Deposit int32    `json:"deposit,omitempty"`    
}

func CreateDBs(DBname string){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
    client.CreateDB(context.TODO(),DBname)
}

func AddAccounts(num int,DBname string){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
	db := client.DB(context.TODO(), DBname)
	num2 := num

	for i:= 0; i < num2; i++{
		Account := CouchDBAccount{Deposit: 100000000}
		id, rev, err := db.CreateDoc(context.TODO(), Account)
		if err != nil {
			panic(err)
		}
		Account.Rev = rev
		Account.Id = id
	 } 
}

func AllDocuments(DBname string)([]*CouchDBAccount, error){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
	db := client.DB(context.TODO(), DBname)

	rows, err := db.AllDocs(context.Background(), kivik.Options{"include_docs": true})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*CouchDBAccount, 0)
	for rows.Next() {
		var account CouchDBAccount
		if err := rows.ScanDoc(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}
func GetRandomCouchDBAccount(accounts[] *CouchDBAccount)(*CouchDBAccount, error){
	if len(accounts) == 0 {
		return nil, fmt.Errorf("沒有可用的帳戶")
	}
	randomIndex := rand.Intn(len(accounts))
	return accounts[randomIndex], nil

} 

/*func CreateIndex(DBname string){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
	db := client.DB(context.TODO(), DBname)
    //創建設計文檔
	designDoc := couchdb.DesignDoc{
		ID: "_design/mydesign",
		Views: map[string]couchdb.ViewDefinition{
			"byId": {
				Map: `function(doc) {
					emit(doc._id, null);
				}`,
			},
		},
	}

	// 保存設計文檔到數據庫
	_, err = db.Put(context.Background(), designDoc.ID, designDoc)
	if err != nil {
		log.Fatal(err)
	}
} */