package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"

	"github.com/carloseduribeiro/golang-sqlc-use-example/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	queries := database.New(dbConn)

	fmt.Println("Creating category...")
	if err = queries.CreateCategory(ctx, database.CreateCategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend description", Valid: true},
	}); err != nil {
		panic(err)
	}

	fmt.Println("\nListing created categories:")
	categories, err := queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		println(category.ID, category.Name, category.Description.String)
	}

	categoryId := categories[0].ID
	fmt.Printf("\nUpdating category with ID %s ...\n", categoryId)
	err = queries.UpdateCategory(ctx, database.UpdateCategoryParams{
		ID:          categoryId,
		Name:        "Backend updated",
		Description: sql.NullString{String: "Backend description updated", Valid: true},
	})

	fmt.Println("\nListing updated category:")
	categories, err = queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		println(category.ID, category.Name, category.Description.String)
	}

	fmt.Printf("Deleting category with ID %s ...\n", categoryId)
	err = queries.DeleteCategory(ctx, categoryId)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nListing categories...")
	categories, err = queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		println(category.ID, category.Name, category.Description.String)
	}

	fmt.Printf("\nEnd\n")
}
