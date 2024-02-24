package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/carloseduribeiro/golang-sqlc-use-example/internal/database"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

type CourseDB struct {
	dbConn *sql.DB
	*database.Queries
}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: database.New(dbConn),
	}
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*database.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := database.New(tx)
	if err = fn(queries); err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("error on rollback: %v, original error: %w", errRb, err)
		}
		return err
	}
	return tx.Commit()
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price       float64
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, argsCategory CategoryParams, argsCourse CourseParams) error {
	err := c.callTx(ctx, func(q *database.Queries) error {
		if err := q.CreateCategory(ctx, database.CreateCategoryParams{
			ID:          argsCategory.ID,
			Name:        argsCategory.Name,
			Description: argsCategory.Description,
		}); err != nil {
			return err
		}
		if err := q.CreateCourse(ctx, database.CreateCourseParams{
			ID:          argsCourse.ID,
			Name:        argsCourse.Name,
			Description: argsCourse.Description,
			CategoryID:  argsCategory.ID,
			Price:       argsCourse.Price,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	courseArgs := CourseParams{
		ID:          uuid.New().String(),
		Name:        "Go",
		Description: sql.NullString{String: "Go Course", Valid: true},
		Price:       10.95,
	}
	categoryArgs := CategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend Course", Valid: true},
	}
	courseDB := NewCourseDB(dbConn)
	if err := courseDB.CreateCourseAndCategory(ctx, categoryArgs, courseArgs); err != nil {
		panic(err)
	}

	queries := database.New(dbConn)
	courses, err := queries.ListCourses(ctx)
	if err != nil {
		panic(err)
	}
	for _, course := range courses {
		fmt.Printf(
			"Category: %s, Course ID: %s, Course Name: %s, Course Description: %s, Course Price: %f\n",
			course.CategoryName, course.ID, course.Name, course.Description.String, course.Price)
	}
}
