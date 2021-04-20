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
	router.HandleFunc("/vehicle/insert", InsertVehicle).Methods("POST")
	router.HandleFunc("/vehicle/update", UpdateVehicle).Methods("PATCH")
	router.HandleFunc("/vehicle/delete/{id}", DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/geospatial/insert", InsertGeoDatas).Methods("POST")

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

	var result *pb.ListVehiclesResponse
	result, err = c.ListVehicles(ctx, &pb.ListVehiclesRequest{})

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

	var result *pb.GetVehicleResponse
	result, err = c.GetVehicle(ctx, &pb.GetVehicleRequest{VehicleId: id})

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

	var result *pb.StatusMessage
	result, err = c.InsertVehicle(ctx, &pb.VehicleMessage{
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

	var result *pb.StatusMessage
	result, err = c.UpdateVehicle(ctx, &pb.VehicleMessage{
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

	var result *pb.ListVehiclesResponse
	result, err = c.DeleteVehicle(ctx, &pb.GetVehicleRequest{VehicleId: id})

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

func InsertGeoDatas(w http.ResponseWriter, r *http.Request) {

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewVehicleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.InsertGeoDatas(ctx, &pb.GeoDatas{
		VehicleId:   data["vehicle_id"],
		Distance:    data["distance"],
		XCoordinate: data["x_coordinate"],
		YCoordinate: data["y_coordinate"],
	})

	if err != nil {
		log.Panic(err)
	}
}
