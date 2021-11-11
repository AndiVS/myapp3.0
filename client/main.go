/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"myapp3.0/client/client"
	"myapp3.0/protocol"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

const (
	username        = "admin"
	password        = "123"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/proto.AuthService/"

	return map[string]bool{
		laptopServicePath + "SingIn": true,
		laptopServicePath + "SingUp": true,
	}
}

func main() {

	cc1, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	c := protocol.NewUserServiceClient(cc2)

	/*	// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		//c := protocol.NewRecordServiceClient(conn)
		c := protocol.NewUserServiceClient(conn)
	*/

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serch, err := c.SearchUser(ctx, &protocol.SearchUserRequest{Username: "admin"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", serch.GetUser())

	/*
		serch, err := c.SearchRecord(ctx, &protocol.SearchRecordRequest{Id: "d30a2bcd-c296-41bf-af9d-2b72eccbb0d0"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", serch.GetRecord())

		update, err := c.UpdateRecord(ctx, &protocol.UpdateRecordRequest{ Record: &protocol.Record{Id:"d30a2bcd-c296-41bf-af9d-2b72eccbb0d0", Name: "update", Type: "type"}})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", update.GetErr())
		/*
		serch, err = c.SearchRecord(ctx, &protocol.SearchRecordRequest{Id: "d30a2bcd-c296-41bf-af9d-2b72eccbb0d0"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", serch.GetRecord())

		del, err := c.DeleteRecord(ctx, &protocol.DeleteRecordRequest{Id: "d30a2bcd-c296-41bf-af9d-2b72eccbb0d0"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", del.GetErr())


		getall, err := c.GetAllRecord(ctx, &protocol.GetAllRecordRequest{ Id: ""})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", getall.GetRecords())

		create, err := c.CreateRecord(ctx, &protocol.CreateRecordRequest{ Record: &protocol.Record{Id:"id", Name: "name", Type: "type"}})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", create.GetId())

	*/

}
