package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"service1/internal/collector"
	"service1/internal/database"
	"service1/proto"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := mydatabase.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	collector := collector.NewCollector(db)

	// Create a gRPC server
	server := grpc.NewServer()

	// Register the collector server with the gRPC server
	proto.RegisterCollectorServer(server, collector)

	// Start the gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Println("gRPC server started on port 50051")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for termination signal to gracefully shutdown the server
	waitForTerminationSignal()
	server.GracefulStop()
}

func waitForTerminationSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop
}
