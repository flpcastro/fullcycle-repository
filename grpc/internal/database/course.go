package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, description, categoryID string) (Course, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)", id, name, description, categoryID)
	if err != nil {
		return Course{}, err
	}

	return Course{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c *Course) FindAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []Course{}
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (c *Course) FindByID(id string) (Course, error) {
	var course Course
	err := c.db.QueryRow("SELECT id, name, description, category_id FROM courses WHERE id = $1", id).Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID)
	if err != nil {
		return Course{}, err
	}

	return course, nil
}

func (c *Course) Update(id, name, description, categoryID string) (Course, error) {
	_, err := c.db.Exec("UPDATE courses SET name = $1, description = $2, category_id = $3 WHERE id = $4", name, description, categoryID, id)
	if err != nil {
		return Course{}, err
	}

	return Course{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c *Course) Delete(id string) error {
	_, err := c.db.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
