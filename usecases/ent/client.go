package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tsukiyoz/demos/usecases/ent/ent/user"

	"github.com/tsukiyoz/demos/usecases/ent/ent"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := ent.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ent?parseTime=true")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	ctx := context.Background()
	//_, _ = CreateUser(ctx, client)
	_, _ = QueryUser(ctx, client)
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(30).SetName("tsukiyo").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Query().Where(user.NameEQ("tsukiyo")).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user was found: ", u)
	return u, nil
}
