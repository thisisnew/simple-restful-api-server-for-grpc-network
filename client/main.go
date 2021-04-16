package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	pb "rest/protos"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Vehicle struct {
	VehicleId           string `json:"vehicle_id"`
	VehicleName         string `json:"vehicle_name"`
	VehicleNumber       string `json:"vehicle_number"`
	VehicleVinNumber    string `json:"vehicle_vin_number"`
	VehicleSerialNumber string `json:"vehicle_serial_number"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", GetVehicleList).Methods("GET")
	router.HandleFunc("/vehicle/{id}", GetVehicle).Methods("GET")
	router.HandleFunc("/insert", InsertVehicle).Methods("POST")
	router.HandleFunc("/update", UpdateVehicle).Methods("PATCH")
	router.HandleFunc("/delete/{id}", DeleteVehicle).Methods("DELETE")

	http.ListenAndServe(":8080", httpHandler(router))
}

func GetVehicleList(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := c.ListVehicles(ctx, &pb.ListVehiclesRequest{})

	if err != nil {
		log.Fatalf("could not request: %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func GetVehicle(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p := mux.Vars(r)
	id := p["id"]

	result, err := c.GetVehicle(ctx, &pb.GetVehicleRequest{VehicleId: id})

	if err != nil {
		log.Fatalf("could not request: %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func InsertVehicle(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newVehicle := Vehicle{}
	err = json.NewDecoder(r.Body).Decode(&newVehicle)

	if err != nil {
		log.Panic(err)
	}

	result, err := c.InsertVehicle(ctx, &pb.VehicleMessage{
		VehicleId:           newVehicle.VehicleId,
		VehicleName:         newVehicle.VehicleName,
		VehicleNumber:       newVehicle.VehicleNumber,
		VehicleVinNumber:    newVehicle.VehicleVinNumber,
		VehicleSerialNumber: newVehicle.VehicleSerialNumber,
	})

	if err != nil {
		log.Panic(err)
	}

	json.NewEncoder(w).Encode(result)
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newVehicle := Vehicle{}
	err = json.NewDecoder(r.Body).Decode(&newVehicle)

	if err != nil {
		log.Panic(err)
	}

	result, err := c.UpdateVehicle(ctx, &pb.VehicleMessage{
		VehicleId:           newVehicle.VehicleId,
		VehicleName:         newVehicle.VehicleName,
		VehicleNumber:       newVehicle.VehicleNumber,
		VehicleVinNumber:    newVehicle.VehicleVinNumber,
		VehicleSerialNumber: newVehicle.VehicleSerialNumber,
	})

	if err != nil {
		log.Panic(err)
	}

	json.NewEncoder(w).Encode(result)
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p := mux.Vars(r)
	id := p["id"]

	result, err := c.DeleteVehicle(ctx, &pb.GetVehicleRequest{VehicleId: id})

	if err != nil {
		log.Fatalf("could not request: %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func httpHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.RemoteAddr, " ", r.Proto, " ", r.Method, " ", r.URL)
		handler.ServeHTTP(w, r)
	})
}
