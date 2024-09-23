package model

import (
	"context"
	"log"
	"slices"
	"strings"
)

type Category struct {
	Id     int32    `json:"id"`
	Title  string   `json:"title"`
	Points []string `json:"points",omitifempty`
}

func (d *Dao) ListCategories() ([]Category, error) {
	categories := make([]Category, 0)

	rows, err := d.pg.Query(
		context.Background(),
		`SELECT id, title FROM categories`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		id := values[0].(int32)
		title := values[1].(string)
		points, err := d.GetCategoryPoints(id)
		if err != nil {
			return nil, err
		}

		categories = append(categories, Category{
			Id:     id,
			Title:  title,
			Points: points,
		})
	}

	return categories, nil
}

func (d *Dao) GetCategoryPoints(id int32) ([]string, error) {
	points := []string{}

	rows, err := d.pg.Query(
		context.Background(),
		`
SELECT points.text
FROM category_points
JOIN
	points ON category_points.point_id = points.id
WHERE
	category_points.category_id = $1
		`,
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		text := values[0].(string)
		points = append(points, text)
	}

	slices.SortFunc(points, strings.Compare)

	return points, nil
}

func (d *Dao) CreateCategory(createReq Category) (category Category, err error) {
	var id int32
	err = d.pg.QueryRow(
		context.Background(),
		"INSERT INTO categories (title) VALUES($1) RETURNING id",
		createReq.Title,
	).Scan(&id)

	if err != nil {
		return
	}

	if err = d.BindCategoryPoints(id, createReq.Points); err != nil {
		return
	}

	category = createReq
	category.Id = id

	return category, nil
}

func (d *Dao) BindCategoryPoints(category_id int32, points []string) error {
	for _, point := range points {
		pointId, err := d.CreateOrGetPoint(point)
		if err != nil {
			return err
		}

		if err := d.AddCategoryPoint(category_id, pointId); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dao) CreateOrGetPoint(text string) (id int32, err error) {
	err = d.pg.QueryRow(
		context.Background(),
		"SELECT id FROM points WHERE text = $1",
		text,
	).Scan(&id)

	if err == nil {
		return id, err
	}

	err = d.pg.QueryRow(
		context.Background(), "INSERT INTO points (text) VALUES ($1) RETURNING id;", text).Scan(&id)

	if err != nil {
		return
	}

	return
}

func (d *Dao) AddCategoryPoint(categoryId int32, pointId int32) (err error) {
	_, err = d.pg.Exec(
		context.Background(),
		"INSERT INTO category_points (category_id, point_id) VALUES ($1, $2)",
		categoryId, pointId,
	)
	return
}

func (d *Dao) UpdateCategory(newCategoryValue Category) (category *Category, err error) {
	if len(newCategoryValue.Title) > 0 {
		cmd, err := d.pg.Exec(
			context.Background(),
			"UPDATE categories SET title = $2 WHERE id = $1",
			newCategoryValue.Id,
			newCategoryValue.Title,
		)

		if err != nil {
			return nil, err
		}

		if cmd.RowsAffected() == 0 {
			log.Printf("UpdateCategory %d category not found.", category.Id)
			return nil, nil
		}
	} else {
		err := d.pg.QueryRow(
			context.Background(),
			"SELECT title FROM categories WHERE id = $1",
			newCategoryValue.Id,
		).Scan(&newCategoryValue.Title)

		if err != nil {
			return nil, err
		}
	}

	_, err = d.pg.Exec(
		context.Background(),
		"DELETE FROM category_points WHERE category_id = $1",
		newCategoryValue.Id,
	)
	if err != nil {
		return
	}

	if len(newCategoryValue.Points) > 0 {
		err = d.BindCategoryPoints(newCategoryValue.Id, newCategoryValue.Points)
		if err != nil {
			return
		}
	} else {
		newCategoryValue.Points = make([]string, 0)
	}

	slices.SortFunc(newCategoryValue.Points, strings.Compare)
	return &newCategoryValue, nil
}

func (d *Dao) DeleteCategory(categoryId int32) (deleted bool, err error) {
	cmd, err := d.pg.Exec(
		context.Background(),
		"DELETE FROM categories WHERE id = $1",
		categoryId,
	)

	if err == nil {
		deleted = cmd.RowsAffected() > 0
	}
	return
}
