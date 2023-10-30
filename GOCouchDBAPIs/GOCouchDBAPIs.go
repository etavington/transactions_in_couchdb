package GOCouchDBAPIs

import (
	"github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"
	//"github.com/go-kivik/kivik/v4/driver/couchdb"
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
	//randomIndex := rand.Intn(200)
	return accounts[randomIndex], nil

} 

/*func CreateView(DBname string){
	client ,err :=kivik.New("couch","http://timo:t102260424@localhost:5984")
	if err != nil{
		panic(err)
	}
	defer client.Close(context.Background())
	db := client.DB(context.TODO(), DBname)
	// 定義 View 的 Map 函數
	mapFunction := `
	function(doc) {
	   emit(doc._id,{_id:doc._id,deposit:doc.deposit});
	}
	`

	// 定義 Design Document 和 View 名稱
	ddocName := DBname+"design"
	viewName := DBname+"view"

	// 創建 Design Document，並設定 View 的 Map 函數
	ddoc := &couchdb.DesignDoc{
		ID:   "_design/" + ddocName,
		Views: map[string]interface{}{
			viewName: map[string]interface{}{
				"map": mapFunction,
			},
		},
	}

	// 上傳 Design Document
	_, err = db.Put(context.Background(), ddoc.ID, ddoc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("View '%s' in design document '%s' created successfully.\n", viewName, ddocName)
}*/

