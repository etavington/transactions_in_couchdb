package main

import (
    "context"
    "log"
    "time"
    "sync"

    "google.golang.org/grpc"
    pb "goclient-and-goserver_bank/payment"
    "goclient-and-goserver_bank/GOCouchDBAPIs"
)

func CallTransferPayments(client pb.TransferServiceClient,payment *pb.Payment,wg *sync.WaitGroup){
  
  //,wg *sync.WaitGroup
  defer wg.Done()
  res,err :=client.TransferPayments(context.Background(),payment)
  if err!=nil{
    log.Fatalf("Fail to call TransferPayments %v", err)
  }
  if res.Response != "Succeed"{
    log.Fatal(err)
  }
   
}

func main() {
    conn, err := grpc.Dial(":50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial server: %v", err)
    }
    defer conn.Close()

    client := pb.NewTransferServiceClient(conn)
   
    //GOCouchDBAPIs.CreateDBs("bank6")
    //GOCouchDBAPIs.CreateDBs("bank7")
    //GOCouchDBAPIs.CreateView("bank7")
    //GOCouchDBAPIs.AddAccounts(1,"bank6")
    //GOCouchDBAPIs.AddAccounts(1,"bank7")
    receiverAccounts, err :=GOCouchDBAPIs.AllDocuments("bank6")
    giverAccounts, err :=GOCouchDBAPIs.AllDocuments("bank7")
    
    var count int
    start := time.Now()

    var wg sync.WaitGroup
    //var wg2 sync.WaitGroup
    //var mu sync.Mutex
    //mu.Lock()
    //defer mu.Unlock()
    var x int =100
    var y int =100
    //wg.Add(times)

    for i := 0; i < x; i++ {
      wg.Add(y)
      for j :=0; j < y; j++{
        go func(){
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
          //,&wg
          CallTransferPayments(client, payment,&wg)
        }()
      }
      wg.Wait()
      //time.Sleep(1*time.Second)
    }

    count = x*y
    end := time.Now()
    elapsed := end.Sub(start)
    seconds := elapsed/time.Second
    rate := float64(count)/float64(seconds)
    log.Println(elapsed)
    log.Println(seconds)
    log.Println(rate)
    log.Println("All payments sent") 
}
