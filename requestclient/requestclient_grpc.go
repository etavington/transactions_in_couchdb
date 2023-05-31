package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "goclient-and-goserver_bank/payment"
    "goclient-and-goserver_bank/GOCouchDBAPIs"
)

func main() {
    conn, err := grpc.Dial(":50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial server: %v", err)
    }
    defer conn.Close()

    client := pb.NewTransferServiceClient(conn)

    stream, err := client.TransferPayments(context.Background())
    if err != nil {
        log.Fatalf("Failed to call SendPayments: %v", err)
    }
    defer stream.CloseSend()
    

    
    GOCouchDBAPIs.CreateDBs("bank6")
    //GOCouchDBAPIs.CreateDBs("bank7")
    //GOCouchDBAPIs.CreateIndex("bank7")
    //GOCouchDBAPIs.AddAccounts(1,"bank6")
    //GOCouchDBAPIs.AddAccounts(1,"bank7")
    receiverAccounts, err :=GOCouchDBAPIs.AllDocuments("bank6")
    giverAccounts, err :=GOCouchDBAPIs.AllDocuments("bank7")
    
    var count int
    start := time.Now()
    for i := 0; i < 10000; i++ {
        giverAccount, err :=GOCouchDBAPIs.GetRandomCouchDBAccount(receiverAccounts)
        if err != nil {
		  log.Fatal(err)
	    }

        receiverAccount, err :=GOCouchDBAPIs.GetRandomCouchDBAccount(giverAccounts)
        if err != nil {
		   log.Fatal(err)
	    } 
       
        payment := &pb.Payment{
            GiverId:    giverAccount.Id,
            ReceiverId: receiverAccount.Id,
            Amount:     1,
        }
        if err := stream.Send(payment); err != nil {
            log.Fatalf("Failed to send payment: %v", err)
        }
        count = i
    }
    end := time.Now()

    elapsed := end.Sub(start)
    seconds := elapsed/time.Second
    rate := float64(count)/float64(seconds)
    log.Println(elapsed)
    log.Println(seconds)
    log.Println(rate)
    log.Println("All payments sent")
    time.Sleep(time.Second * 20) 
}
