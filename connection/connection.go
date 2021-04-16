package connection

import (
	"google.golang.org/grpc"
	"log"
	pb "rest/protos"
)

func Connect() pb.VehicleClient {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	return c
}
