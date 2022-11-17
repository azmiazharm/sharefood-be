package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sharefood/internal/consts"
	"sharefood/internal/entity"
	"sharefood/pkg/logger"
	"sharefood/pkg/postgres"
	"sharefood/pkg/tracer"
	"time"

	"github.com/google/uuid"
)

type Request interface {
	ListbyFoodUser(ctx context.Context, idFood uuid.UUID, idUser uuid.UUID) ([]entity.Request, error)
	ListbyFood(ctx context.Context, idFood uuid.UUID) ([]entity.Request, error)
	ListbyUser(ctx context.Context, idUser uuid.UUID) ([]entity.Request, error)
	Create(context.Context, *entity.Request) error
	GetRequestFoodByIDRequest(ctx context.Context, idRequest uuid.UUID) (entity.RequestWithFood, error)
	AcceptRequest(ctx context.Context, idRequest uuid.UUID) error
	RejectRequest(ctx context.Context, idRequest uuid.UUID) error
	// GetByEmail(context.Context, string) (entity.User, error)
	// IsRegistered(context.Context, string) bool
}

type requestImplementation struct {
	conn postgres.Adapter
}

func NewRequestRepository(conn postgres.Adapter) Request {
	return &requestImplementation{conn}
}

// Get all requests in the by id user and id food
func (r requestImplementation) ListbyFoodUser(ctx context.Context, idFood uuid.UUID, idUser uuid.UUID) (requests []entity.Request, err error) {
	errorEvent := consts.ErrorEvent("update_my_foods")
	ctx = tracer.SpanStart(ctx, "update_my_foods")
	defer tracer.SpanFinish(ctx)

	query := `SELECT id_request, id_user, id_food, status, quantity, created_at, updated_at from requests 
		WHERE id_food = $1
		ORDER BY updated_at DESC`
	rows, err := r.conn.QueryRows(ctx, query, idFood)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request entity.Request
		err := rows.Scan(
			&request.ID,
			&request.IDUser,
			&request.IDFood,
			&request.Status,
			&request.Quantity,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
			tracer.SpanError(ctx, err)
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

// Get all requests in the by id food
func (r requestImplementation) ListbyFood(ctx context.Context, idFood uuid.UUID) (requests []entity.Request, err error) {
	errorEvent := consts.ErrorEvent("list_requests_food")
	ctx = tracer.SpanStart(ctx, "list_requests_food")
	defer tracer.SpanFinish(ctx)

	query := `SELECT id_request, id_user, id_food, status, quantity, created_at, updated_at from requests 
		WHERE id_food = $1
		ORDER BY updated_at DESC`
	rows, err := r.conn.QueryRows(ctx, query, idFood)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request entity.Request
		err := rows.Scan(
			&request.ID,
			&request.IDUser,
			&request.IDFood,
			&request.Status,
			&request.Quantity,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
			tracer.SpanError(ctx, err)
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

// Get all requests in the by id user
func (r requestImplementation) ListbyUser(ctx context.Context, idUser uuid.UUID) (requests []entity.Request, err error) {
	errorEvent := consts.ErrorEvent("update_my_foods")
	ctx = tracer.SpanStart(ctx, "update_my_foods")
	defer tracer.SpanFinish(ctx)

	query := `SELECT id_request, id_user, id_food, status, quantity, created_at, updated_at from requests 
		WHERE id_user = $1
		ORDER BY updated_at DESC`
	rows, err := r.conn.QueryRows(ctx, query, idUser)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request entity.Request
		err := rows.Scan(
			&request.ID,
			&request.IDUser,
			&request.IDFood,
			&request.Status,
			&request.Quantity,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
			tracer.SpanError(ctx, err)
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

// Create new request
func (r requestImplementation) Create(ctx context.Context, request *entity.Request) (err error) {
	errorEvent := consts.ErrorEvent("update_my_foods")
	ctx = tracer.SpanStart(ctx, "update_my_foods")
	defer tracer.SpanFinish(ctx)

	query := `
	INSERT INTO requests(id_request, id_user, id_food, quantity)
	VALUES ($1, $2, $3, $4)
	`

	_, err = r.conn.Exec(
		ctx,
		query,
		request.ID,
		request.IDUser,
		request.IDFood,
		request.Quantity,
	)

	if err != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

func (r requestImplementation) GetRequestFoodByIDRequest(ctx context.Context, idRequest uuid.UUID) (reqFood entity.RequestWithFood, err error) {
	errorEvent := consts.ErrorEvent("get_request_food_by_id_request")
	ctx = tracer.SpanStart(ctx, "get_request_food_by_id_request")
	defer tracer.SpanFinish(ctx)

	query := `
	SELECT 
		requests.id_request,
		requests.id_user,
		requests.id_food,
		requests.status,
		requests.quantity,
		foods.id_user AS giver,
		foods.quantity AS stock
	FROM requests
	INNER JOIN foods
	ON requests.id_food = foods.id_food
	WHERE id_request = $1 AND status = 0;
	`
	row := r.conn.QueryRow(ctx, query, idRequest)
	fmt.Println(row)
	err = row.Scan(
		&reqFood.ID,
		&reqFood.IDUser,
		&reqFood.IDFood,
		&reqFood.Status,
		&reqFood.Quantity,
		&reqFood.IDUserFood,
		&reqFood.Stock,
	)
	if err == sql.ErrNoRows {
		err := errorEvent.WithCode(consts.CodeUnprocessableEntity).WrapError(consts.Error(consts.ActionAlreadyDone))
		tracer.SpanError(ctx, err)
		return reqFood, err
	}
	if err != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(consts.Error(err.Error()))
		tracer.SpanError(ctx, err)
		return reqFood, err
	}

	return reqFood, nil
}

// Accept Request, only owner can update
func (r requestImplementation) AcceptRequest(ctx context.Context, idRequest uuid.UUID) (err error) {
	errorEvent := consts.ErrorEvent("accept_request")
	ctx = tracer.SpanStart(ctx, "accept_request")
	defer tracer.SpanFinish(ctx)

	updatedTime := time.Now().Local()

	// DB Transaction with BeginTx
	// ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Begin query 1, not applied yet to db, only to trx
	query1 := `
	UPDATE requests
    SET
        status = 1,
    	updated_at = $1
	WHERE
		id_request = $2;
	`
	_, err = tx.ExecContext(ctx, query1, updatedTime, idRequest)
	if err != nil {
		err = errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()

		return err
	}

	// Begin query 2, not applied yet to db, only to trx
	query2 := `	
	UPDATE foods
	SET 
    	quantity = foods.quantity - requests.quantity,
    	updated_at = $1
	FROM requests
	WHERE 
		requests.id_food = foods.id_food
	AND 
		requests.id_request = $2;
	`
	_, err = tx.ExecContext(ctx, query2, updatedTime, idRequest)
	if err != nil {
		err = errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()

		return err
	}

	// close the transaction with a Commit() or Rollback() method on the resulting Tx variable.
	// this applies the above changes to our database
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Reject Request, only owner can update
func (r requestImplementation) RejectRequest(ctx context.Context, idRequest uuid.UUID) (err error) {
	errorEvent := consts.ErrorEvent("reject_request")
	ctx = tracer.SpanStart(ctx, "reject_request")
	defer tracer.SpanFinish(ctx)

	query := `
		UPDATE requests SET 
			status = 2, 
			updated_at = $1
		WHERE id_request=$2;

		`
	updatedTime := time.Now().Local()

	_, err = r.conn.Exec(ctx, query, updatedTime, idRequest)
	if err != nil {
		err := errorEvent.WithCode(consts.CodeInternalServerError).WrapError(err)
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

// // Get single user by ID
// func (r userImplementation) GetByID(ctx context.Context, id int64) (user entity.User, err error) {
// 	query := `
// 		SELECT name
// 		FROM users
// 		WHERE id = $1
// 	`

// 	row := r.conn.QueryRow(ctx, query, id)
// 	err = row.Scan(
// 		&user.Name,
// 	)
// 	if err != nil {
// 		err = fmt.Errorf("scanning user %w", err)
// 		return entity.User{}, err
// 	}

// 	return user, nil
// }

// // Get single user by email
// func (r userImplementation) GetByEmail(ctx context.Context, email string) (user entity.User, err error) {
// 	query := `
// 		SELECT id_user, name, email, phone_number, password, image_url
// 		FROM users
// 		WHERE (email = $1) AND (deleted_at IS NULL)
// 	`

// 	row := r.conn.QueryRow(ctx, query, email)

// 	err = row.Scan(
// 		&user.ID,
// 		&user.Name,
// 		&user.Email,
// 		&user.PhoneNumber,
// 		&user.Password,
// 		&user.ImageUrl,
// 	)

// 	if err != nil {
// 		err = fmt.Errorf("scanning user %w", err)
// 		return entity.User{}, err
// 	}

// 	return user, nil
// }

// // Check is user registered by querying 1 matching email
// func (r userImplementation) IsRegistered(ctx context.Context, email string) bool {
// 	query := `
// 		SELECT name
// 		FROM users
// 		WHERE email = $1
// 		LIMIT 1
// 	`
// 	var name string
// 	err := r.conn.QueryRow(ctx, query, email).Scan(&name)
// 	fmt.Println(err)
// 	if err == sql.ErrNoRows {
// 		return false
// 	}
// 	return true
// }
