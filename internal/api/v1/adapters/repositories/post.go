package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"example.com/m/internal/api/v1/core/application/dto"
	object_storage "example.com/m/internal/api/v1/infrastructure/s3"
	"example.com/m/internal/config"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type PostRepository struct {
	db     *sql.DB
	logger *zap.Logger
	s3     *object_storage.ClientS3
}

func NewPostRepository(db *sql.DB, logger *zap.Logger, s3 *object_storage.ClientS3) *PostRepository {
	return &PostRepository{
		db:     db,
		logger: logger,
		s3:     s3,
	}
}

func (r *PostRepository) Create(ctx context.Context, p *dto.PostDtoWithoutId) (*int64, error) {
	var id int64

	// Manually construct the query to handle the images field
	query := `
        INSERT INTO posts (images, user_email, place_id, title, description, genre, author, publication_year, publisher, condition, status, created_at, cover, pages_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id
    `

	// Execute the query with the images field converted to a PostgreSQL array
	err := r.db.QueryRowContext(ctx, query, pq.Array(p.Images), p.UserEmail, p.PlaceID, p.Title, p.Description, p.Genre, p.Author, p.PublicationYear, p.Publisher, p.Condition, p.Status, p.CreatedAt, &p.Cover, &p.PagesCount).Scan(&id)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "Create"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &id, nil
}

func (r *PostRepository) Get(ctx context.Context, id int64, userEmail string) (*dto.PostDto, error) {
	var post dto.PostDto
	var images pq.StringArray
	query, _, _ := goqu.
		Select(
			"posts.*",
			goqu.L("places.name AS place_name"),
			goqu.L("places.address AS place_address"),
			goqu.L("users.username AS owner_username"),
			goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
		).
		From("posts").
		Join(
			goqu.T("places"),
			goqu.On(goqu.I("posts.place_id").Eq(goqu.I("places.id"))),
		).
		Join(
			goqu.T("users"),
			goqu.On(goqu.I("posts.user_email").Eq(goqu.I("users.email"))),
		).
		Where(
			goqu.I("posts.id").Eq(id),
		).
		GroupBy("posts.id", "places.name", "places.address", "users.username").
		ToSQL()
	err := r.db.QueryRow(query).
		Scan(&post.ID, &images, &post.UserEmail,
			&post.PlaceID, &post.Title, &post.Description, &post.Genre,
			&post.Author, &post.PublicationYear, &post.Publisher, &post.Condition,
			&post.Status, &post.CreatedAt, &post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.PlaceName, &post.PlaceAddress, &post.OwnerUsername, &post.IsFavorite)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "Get"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	if images == nil {
		post.Images = []string{}
	} else {
		post.Images = images
	}
	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, id int64, p *dto.UpdatePostDto) error {
	var uMap map[string]interface{}
	inrec, _ := json.Marshal(*p)
	json.Unmarshal(inrec, &uMap)
	var rec goqu.Record = uMap
	query, _, _ := goqu.From("posts").Where(goqu.C("id").Eq(id)).Update().Set(
		rec,
	).ToSQL()

	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "Update"),
			zap.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id int64) error {
	query, _, _ := goqu.From("posts").Where(goqu.C("id").Eq(id)).Delete().ToSQL()
	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "Delete"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

type PostFilterOptions struct {
	Genre           string
	Condition       string
	PublicationYear string
	PlaceID         int64
}

// TODO таргетинг
func (r *PostRepository) GetAllAvailable(ctx context.Context, userEmail string, options PostFilterOptions, limit, offset uint) (*[]dto.PostDto, error) {
	var posts []dto.PostDto
	query := goqu.
		Select(
			"posts.*",
			goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
		).
		From("posts").
		Where(goqu.C("status").Eq("available")).
		Where(goqu.I("posts.user_email").Neq(userEmail)).
		GroupBy("posts.id")
	if options.Genre != "" {
		query = query.Where(goqu.C("genre").Eq(options.Genre))
	}
	if options.Condition != "" {
		query = query.Where(goqu.C("condition").Eq(options.Condition))
	}
	if options.PublicationYear != "" {
		query = query.Where(goqu.C("publication_year").Eq(options.PublicationYear))
	}
	if options.PlaceID != 0 {
		query = query.Where(goqu.C("place_id").Eq(options.PlaceID))
	}
	query = query.Order(goqu.I("posts.created_at").Desc())
	sql, args, err := query.Limit(limit).Offset(offset).ToSQL()
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllAvailable"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllAvailable"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostDto
		var images pq.StringArray
		if err := rows.Scan(
			&post.ID, &images, &post.UserEmail,
			&post.PlaceID, &post.Title, &post.Description, &post.Genre,
			&post.Author, &post.PublicationYear, &post.Publisher,
			&post.Condition, &post.Status, &post.CreatedAt, &post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.IsFavorite); err != nil {
			return nil, err
		}
		if images == nil {
			post.Images = []string{}
		} else {
			post.Images = images
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (r *PostRepository) SearchByTitleOrAuthorOrGenre(ctx context.Context, titleOrAuthorOrGenre string, limit, offset uint, userEmail string) (*[]dto.PostDto, error) {
	var posts []dto.PostDto
	query, _, _ := goqu.
		Select(
			"posts.*",
			goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
		).
		From("posts").Where(
		goqu.Or(
			goqu.C("title").ILike("%"+titleOrAuthorOrGenre+"%"),
			goqu.C("author").ILike("%"+titleOrAuthorOrGenre+"%"),
			goqu.C("genre").ILike("%"+titleOrAuthorOrGenre+"%"),
		),
	).
		Where(goqu.And(
			goqu.C("status").Eq("available"),
			goqu.I("posts.user_email").Neq(userEmail),
		)).
		GroupBy("posts.id").
		Order(goqu.I("posts.created_at").Desc()).
		Limit(limit).Offset(offset).ToSQL()

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "SearchByTitleOrAuthor"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostDto
		var images pq.StringArray
		if err := rows.Scan(
			&post.ID, &images, &post.UserEmail,
			&post.PlaceID, &post.Title, &post.Description, &post.Genre,
			&post.Author, &post.PublicationYear, &post.Publisher,
			&post.Condition, &post.Status, &post.CreatedAt,
			&post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.IsFavorite); err != nil {
			return nil, err
		}
		if images == nil {
			post.Images = []string{}
		} else {
			post.Images = images
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (r *PostRepository) GetAllMyPosted(ctx context.Context, userEmail string, status string, limit, offset uint) (*[]dto.PostDto, error) {
	var posts []dto.PostDto
	// query, _, _ := goqu.
	// 	Select(
	// 		"posts.*",
	// 		goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
	// 	).
	// 	From("posts").Where(goqu.Ex{
	// 	"posts.user_email": userEmail,
	// 	"status":           status,
	// }).
	// 	GroupBy("posts.id").
	// 	Order(goqu.I("posts.created_at").Desc()).
	// 	Limit(uint(limit)).Offset(uint(offset)).ToSQL()

	var query string
	if status != "all" {
		query, _, _ = goqu.
			Select(
				"posts.*",
				goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
			).
			From("posts").Where(goqu.Ex{
			"posts.user_email": userEmail,
			"status":           status,
		}).
			GroupBy("posts.id").
			Order(goqu.I("posts.created_at").Desc()).
			Limit(uint(limit)).Offset(uint(offset)).ToSQL()
	} else {
		query, _, _ = goqu.
			Select(
				"posts.*",
				goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
			).
			From("posts").Where(goqu.Ex{
			"posts.user_email": userEmail,
		}).
			GroupBy("posts.id").
			Order(goqu.I("posts.created_at").Desc()).
			Limit(uint(limit)).Offset(uint(offset)).ToSQL()
	}

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllMyPosted"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostDto
		var images pq.StringArray
		if err := rows.Scan(
			&post.ID, &images, &post.UserEmail,
			&post.PlaceID, &post.Title, &post.Description, &post.Genre,
			&post.Author, &post.PublicationYear, &post.Publisher,
			&post.Condition, &post.Status, &post.CreatedAt,
			&post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.IsFavorite); err != nil {
			r.logger.Error(
				"Post Repository Error",
				zap.String("method", "GetAllMyPosted"),
				zap.String("error", err.Error()),
			)
			return nil, err
		}
		if images == nil {
			post.Images = []string{}
		} else {
			post.Images = images
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllMyPosted"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &posts, nil
}

func (r *PostRepository) AddFavorite(ctx context.Context, postID int64, userEmail string) error {
	query, _, _ := goqu.
		Insert("favorites").
		Rows(goqu.Record{
			"user_email": userEmail,
			"post_id":    postID,
		}).
		OnConflict(goqu.DoNothing()).
		ToSQL()

	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "AddFavorite"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PostRepository) DeleteFavorite(ctx context.Context, postID int64, userEmail string) error {
	query, _, _ := goqu.
		Delete("favorites").
		Where(goqu.Ex{
			"user_email": userEmail,
			"post_id":    postID,
		}).
		ToSQL()

	_, err := r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "DeleteFavorite"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PostRepository) GetAllMyBooked(ctx context.Context, userEmail string, limit, offset uint) (*[]dto.PostDto, error) {
	var posts []dto.PostDto
	query, _, err := goqu.
		Select(
			"posts.*",
			goqu.L("(SELECT COUNT(*) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?) AS is_favourite", userEmail),
		).
		From("posts").
		Join(
			goqu.T("bookings"),
			goqu.On(
				goqu.I("posts.id").Eq(goqu.I("bookings.post_id")),
			),
		).
		Where(goqu.Ex{"bookings.user_email": userEmail}).
		GroupBy("posts.id").
		Limit(limit).
		Offset(offset).
		ToSQL()
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllMyBooked"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllMyBooked"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostDto
		var images pq.StringArray
		if err := rows.Scan(
			&post.ID, &images, &post.UserEmail, &post.PlaceID,
			&post.Title, &post.Description, &post.Genre, &post.Author, &post.PublicationYear,
			&post.Publisher, &post.Condition, &post.Status, &post.CreatedAt,
			&post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.IsFavorite); err != nil {
			r.logger.Error(
				"Post Repository Error",
				zap.String("method", "GetAllMyBooked"),
				zap.String("error", err.Error()),
			)
			return nil, err
		}
		if images == nil {
			post.Images = []string{}
		} else {
			post.Images = images
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetAllMyBooked"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &posts, nil
}

func (r *PostRepository) GetAllMyTaken(ctx context.Context, userEmail string) {

}

func (r *PostRepository) GetAllFavourite(ctx context.Context, userEmail string, limit, offset uint) (*[]dto.PostDto, error) {
	var posts []dto.PostDto
	isFavoriteQuery := "(SELECT COUNT(id) > 0 FROM favorites WHERE favorites.post_id = posts.id AND favorites.user_email = ?)"
	query, _, err := goqu.
		Select(
			"posts.*",
			goqu.L(isFavoriteQuery+" AS is_favourite", userEmail),
		).
		Where(
			goqu.L(isFavoriteQuery, userEmail).Eq(true),
		).
		Where(goqu.I("posts.status").Eq("available")).
		From("posts").
		Limit(limit).
		Offset(offset).
		ToSQL()
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetFavoritePosts"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetFavoritePosts"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post dto.PostDto
		var images pq.StringArray
		if err := rows.Scan(
			&post.ID, &images, &post.UserEmail, &post.PlaceID,
			&post.Title, &post.Description, &post.Genre, &post.Author, &post.PublicationYear,
			&post.Publisher, &post.Condition, &post.Status, &post.CreatedAt,
			&post.Cover, &post.PagesCount, &post.Summary, &post.Quote, &post.IsFavorite); err != nil {
			r.logger.Error(
				"Post Repository Error",
				zap.String("method", "GetFavoritePosts"),
				zap.String("error", err.Error()),
			)
			return nil, err
		}
		if images == nil {
			post.Images = []string{}
		} else {
			post.Images = images
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "GetFavoritePosts"),
			zap.String("error", err.Error()),
		)
		return nil, err
	}

	return &posts, nil
}

func (r *PostRepository) AddImage(
	ctx context.Context, image multipart.File, imageName string, imageType string,
) (string, error) {
	fileName := fmt.Sprintf("%s.%s", uuid.NewString(), imageType)
	fmt.Println(fileName)
	err := r.s3.UploadFile(ctx, "posts", fileName, image)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "AddImage"),
			zap.String("error", err.Error()),
		)
		return "", err
	}
	return fmt.Sprintf("%s/%s", config.Config.S3Endpoint, fileName), nil
}

func (r *PostRepository) UpdateImageURL(ctx context.Context, postID int64, imageURL string) error {
	query, _, err := goqu.
		Update("posts").
		Set(goqu.Record{"images": goqu.L("array_append(images, ?)", imageURL)}).
		Where(goqu.Ex{"id": postID}).
		ToSQL()

	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "UpdateImageURL"),
			zap.String("error", err.Error()),
		)
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "UpdateImageURL"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PostRepository) SetSummary(ctx context.Context, postID int64, summary string) error {
	query, _, err := goqu.
		Update("posts").
		Set(goqu.Record{"summary": summary}).
		Where(goqu.Ex{"id": postID}).
		ToSQL()

	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "SetSummary"),
			zap.String("error", err.Error()),
		)
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "SetSummary"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (r *PostRepository) SetQuote(ctx context.Context, postID int64, quote string) error {
	query, _, err := goqu.
		Update("posts").
		Set(goqu.Record{"quote": quote}).
		Where(goqu.Ex{"id": postID}).
		ToSQL()

	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "SetQuote"),
			zap.String("error", err.Error()),
		)
		return err
	}
	_, err = r.db.Exec(query)
	if err != nil {
		r.logger.Error(
			"Post Repository Error",
			zap.String("method", "SetQuote"),
			zap.String("error", err.Error()),
		)
		return err
	}
	return nil
}
