package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "goclient-and-goserver_bank/payment"
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

//    row, err :=

    payment := &pb.Payment{
        GiverId:    "05525c5942ba5aeabd923ad0e7931b89",
        ReceiverId: "05525c5942ba5aeabd923ad0e7932338",
        Amount:     1,
    }
    var count int

    start := time.Now()
    for i := 0; i < 3000; i++ {
        if err := stream.Send(payment); err != nil {
            log.Fatalf("Failed to send payment: %v", err)
        }
        //time.Sleep(time.Second * 1)
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
