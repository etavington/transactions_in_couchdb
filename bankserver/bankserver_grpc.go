package main

import (
    "context"
    "log"
    "net"
    "sync"
    "github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb"goclient-and-goserver_bank/payment"
)


type CouchDBServer struct{
    pb.UnimplementedTransferServiceServer
    client *kivik.Client 
}

type CouchDBAccount struct {
	Id      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"` 
	Deposit int32    `json:"deposit,omitempty"`    
}

var m sync.Mutex
var m1 sync.Mutex
var m2 sync.Mutex

func Update(owner CouchDBAccount,newDeposit int32,db *kivik.DB)(){
 //   var mu sync.Mutex
 //   grpccouchdbserver.mu.Lock()
 //   defer grpccouchdbserver.mu.Unlock()
 //   ,mu *sync.Mutex

    owner.Rev = owner.Rev // Must be set
	owner.Deposit = newDeposit
    newRev, err := db.Put(context.TODO(), owner.Id, owner)
	if err != nil {
		panic(err)
	}
	owner.Rev = newRev 
    
    return
}

func TransferCouchDB(giver_id string,receiver_id string,amount int32,db1 *kivik.DB,db2 *kivik.DB)(){
   // grpccouchdbserver.mu.Lock()
   // defer grpccouchdbserver.mu.Unlock()
   //,mu *sync.Mutex

    var giverAccount CouchDBAccount
    var receiverAccount CouchDBAccount

    m1.Lock()
    err := db1.Get(context.Background(), giver_id).ScanDoc(&giverAccount)
    if err != nil {
        panic(err)
    } 
    /*viewName1 := "_view/id1"
    designDoc1 := "_design/bank1_view"
    options1 := kivik.Options{
        "key": giver_id,
        "include_docs": true,
    }
    rows1, err := db1.Query(context.Background(), designDoc1, viewName1, options1)
    if err != nil {
        panic(err)
    }
        if rows1.Next() { 
            if err := rows1.ScanDoc(&giverAccount); err != nil {
               panic(err)
            }
            if(giverAccount.Id=="dfsf"){
                println(giverAccount.Deposit)
            }
        }else{
          println("No rows found")
        }*/
    if giverAccount.Deposit < amount{
        log.Fatalf("Giver doesn't have enough deposit")
        return
    }
    
    giverNewdeposit := giverAccount.Deposit - amount
    /*if(giverNewdeposit==94935){}*/
    Update(giverAccount, giverNewdeposit,db1) 
    m1.Unlock()

    m2.Lock()
    err = db2.Get(context.Background(), receiver_id).ScanDoc(&receiverAccount)
    if err != nil {
        panic(err)
    } 
    /*viewName2 := "_view/id2"
    designDoc2 := "_design/bank2_view"
    options2 := kivik.Options{
        "key": receiver_id,
        "include_docs": true,
    }
    rows2, err := db2.Query(context.Background(), designDoc2, viewName2, options2)
    if err != nil {
        panic(err)
    }
    if rows2.Next() { 
       if err := rows2.ScanDoc(&receiverAccount); err != nil {
            panic(err)
       }
       if(receiverAccount.Id=="ytryr"){
           println(receiverAccount.Deposit)
       }
    }else{
        println("No rows found")
    }*/
    receiverNewdeposit := receiverAccount.Deposit + amount  
    //if(receiverNewdeposit==94935){}
    Update(receiverAccount, receiverNewdeposit,db2)
    m2.Unlock()

    return
}

func (grpccouchdbserver *CouchDBServer) TransferPayments(ctx context.Context,payment *pb.Payment) (*pb.Response,error) {
    //var mut sync.Mutex = grpccouchdbserver.mu
    //,&grpccouchdbserver.mu
    //grpccouchdbserver.mu.Lock()
    //defer grpccouchdbserver.mu.Unlock()
    //var giverAccount CouchDBAccount

    db1 := grpccouchdbserver.client.DB(context.TODO(), "bank6")
    db2 := grpccouchdbserver.client.DB(context.TODO(), "bank7") 
    giver_id :=payment.GetGiverId()
    receiver_id :=payment.GetReceiverId()
    amount :=payment.GetAmount() 
    TransferCouchDB(giver_id ,receiver_id ,amount,db1,db2) 

    return &pb.Response{
           Response: "Succeed",
         },nil   
}

func NewCouchDBServer()(*CouchDBServer,error){
   client, err := kivik.New("couch", "http://timo:t102260424@localhost:5984")
   if err != nil {
      panic(err)
   }
   defer client.Close(context.Background())

   couchdbserver := &CouchDBServer{
      client : client,
   }
   
   return couchdbserver,nil
}

func main() {
    couchdbserver, err:= NewCouchDBServer()
    if err != nil{
       log.Fatal(err) 
    }

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    grpcserver := grpc.NewServer()
    pb.RegisterTransferServiceServer(grpcserver, couchdbserver)
    reflection.Register(grpcserver)
    if err := grpcserver.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
