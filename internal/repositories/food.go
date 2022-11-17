package repositories

import (
	"context"
	"fmt"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/pkg/logger"
	"sharefood/pkg/postgres"
	"sharefood/pkg/tracer"
	"time"

	"github.com/google/uuid"
)

type Food interface {
	List(context.Context) ([]entity.Food, error)
	GetDetailByID(context.Context, uuid.UUID) (entity.Food, error)
	DeleteByID(context.Context, uuid.UUID) error
	Create(context.Context, *entity.Food) error
	Update(context.Context, *entity.Food) error
	ListMy(context.Context, uuid.UUID) ([]entity.Food, error)
}

type foodImplementation struct {
	conn postgres.Adapter
}

func NewFoodRepository(conn postgres.Adapter) Food {
	return &foodImplementation{conn}
}

// Get all foods in the repository
func (r foodImplementation) List(ctx context.Context) (foods []entity.Food, err error) {
	query := `
		SELECT 
			id_food, 
			id_user, 
			name, 
			description, 
			category, 
			quantity, 
			image_url,
			expired_at,
			latitude,
			longitude
		FROM foods`
	rows, err := r.conn.QueryRows(ctx, query)
	fmt.Println(rows)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var food entity.Food
		err := rows.Scan(
			&food.ID,
			&food.IDUser,
			&food.Name,
			&food.Description,
			&food.Category,
			&food.Quantity,
			&food.ImageUrl,
			// &food.Location,
			&food.ExpiredAt,
			&food.Latitude,
			&food.Longitude,
		)

		if err != nil {
			return nil, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

// Get single user by ID
func (r foodImplementation) GetDetailByID(ctx context.Context, id uuid.UUID) (food entity.Food, err error) {
	query := `
		SELECT 
			id_food, 
			id_user, 
			name, 
			description, 
			category, 
			quantity, 
			image_url,
			expired_at,
			latitude,
			longitude
		FROM foods
		WHERE id_food = $1 AND deleted_at IS NULL;
	`

	row := r.conn.QueryRow(ctx, query, id)
	err = row.Scan(
		&food.ID,
		&food.IDUser,
		&food.Name,
		&food.Description,
		&food.Category,
		&food.Quantity,
		&food.ImageUrl,
		// &food.Location,
		&food.ExpiredAt,
		&food.Latitude,
		&food.Longitude,
	)
	if err != nil {
		err = fmt.Errorf("scanning food %w", err)
		return entity.Food{}, err
	}

	return food, nil
}

// Create Food
func (r foodImplementation) Create(ctx context.Context, food *entity.Food) (err error) {
	query := `
	INSERT INTO foods(
		id_food,
		id_user,
		name,
		description,
		category,
		quantity,
		image_url,
		is_active,
		expired_at,
		latitude,
		longitude
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.conn.Exec(
		ctx,
		query,
		food.ID,
		food.IDUser,
		food.Name,
		food.Description,
		food.Category,
		food.Quantity,
		food.ImageUrl,
		true,
		food.ExpiredAt,
		food.Latitude,
		food.Longitude,
	)

	if err != nil {
		fmt.Println(fmt.Errorf("executing query: %w", err))
		return err
	}

	return nil
}

// Update My food, only owner can update
func (r foodImplementation) Update(ctx context.Context, food *entity.Food) (err error) {
	errorEvent := consts.ErrorEvent("update_my_foods")
	ctx = tracer.SpanStart(ctx, "update_my_foods")
	defer tracer.SpanFinish(ctx)

	query := `
		UPDATE foods SET 
			name = $1, 
			description = $2, 
			category = $3, 
			quantity = $4, 
			image_url = $5,
			expired_at = $6,
			latitude = $7,
			longitude= $8,
			updated_at = $9
		WHERE id_food=$10;

		`
	updatedTime := time.Now().Local()

	_, err = r.conn.Exec(ctx, query, food.Name, food.Description, food.Category, food.Quantity, food.ImageUrl, food.ExpiredAt, food.Latitude, food.Longitude, updatedTime, food.ID)
	if err != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

// List my food
func (r foodImplementation) ListMy(ctx context.Context, id_user uuid.UUID) (foods []entity.Food, err error) {
	errorEvent := consts.ErrorEvent("list_my_foods")
	ctx = tracer.SpanStart(ctx, "list_my_foods")
	defer tracer.SpanFinish(ctx)

	query := `
		SELECT 
			id_food, 
			id_user, 
			name, 
			description, 
			category, 
			quantity, 
			image_url,
			expired_at,
			latitude,
			longitude
		FROM foods
		WHERE id_user=$1 AND deleted_at IS NULL;

		`
	rows, errQueryRows := r.conn.QueryRows(ctx, query, id_user)
	fmt.Println(rows)
	if errQueryRows != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(errQueryRows)
		tracer.SpanError(ctx, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var food entity.Food
		errStructScan := rows.Scan(
			&food.ID,
			&food.IDUser,
			&food.Name,
			&food.Description,
			&food.Category,
			&food.Quantity,
			&food.ImageUrl,
			// &food.Location,
			&food.ExpiredAt,
			&food.Latitude,
			&food.Longitude,
		)

		if errStructScan != nil {
			err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(errStructScan)
			tracer.SpanError(ctx, err)
			return nil, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

// Delete food by idFood
func (r foodImplementation) DeleteByID(ctx context.Context, idFood uuid.UUID) (err error) {
	errorEvent := consts.ErrorEvent("delete_my_food")
	ctx = tracer.SpanStart(ctx, "delete_my_food")
	defer tracer.SpanFinish(ctx)

	query := `
		UPDATE foods SET 
			is_active = $1, 
			deleted_at = $2
		WHERE id_food=$3 AND deleted_at != NULL;;

		`

	deletedTime := time.Now().Local()

	_, err = r.conn.Exec(ctx, query, false, deletedTime, idFood)
	if err != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}
