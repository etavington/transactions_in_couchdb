package main

import (
    "context"
//    "fmt"
    "log"
    "net"
//    "time"
    "github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb"goclient-and-goserver_bank/payment"
)


type server struct {
    pb.UnimplementedTransferServiceServer
}

type CouchDBAccount struct {
	Id      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"` 
	Deposit int32    `json:"deposit,omitempty"`    
}


func Update(owner CouchDBAccount,newDeposit int32,db *kivik.DB)(){
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
    
    var giverAccount CouchDBAccount
    var receiverAccount CouchDBAccount

    err := db1.Get(context.Background(), giver_id).ScanDoc(&giverAccount)
    if err != nil {
        panic(err)
    }

    err = db2.Get(context.Background(), receiver_id).ScanDoc(&receiverAccount)
    if err != nil {
        panic(err)
    }

    if giverAccount.Deposit < amount{
        log.Fatalf("Giver doesn't have enough deposit")
        return
    }
    
    giverNewdeposit := giverAccount.Deposit - amount
    receiverNewdeposit := receiverAccount.Deposit + amount
    
    Update(giverAccount, giverNewdeposit,db1)
    Update(receiverAccount, receiverNewdeposit,db2)
   
    return
}

var count int =1
var str string ="string"
func (s *server) TransferPayments(stream pb.TransferService_TransferPaymentsServer) error {
    client, err := kivik.New("couch", "http://timo:t102260424@localhost:5984")
    if err != nil {
        panic(err)
    }
    defer client.Close(context.Background()) 
    db1 := client.DB(context.TODO(), "bank6")
    db2 := client.DB(context.TODO(), "bank7") 
    for {
        payment, err := stream.Recv()
        if err != nil {
            return err
        }
        giver_id :=payment.GetGiverId()
        receiver_id :=payment.GetReceiverId()
        amount :=payment.GetAmount()
        TransferCouchDB(giver_id ,receiver_id ,amount,db1,db2)
    }    
    
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterTransferServiceServer(s, &server{})
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
