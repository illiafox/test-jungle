package storages

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"jungle-test/app/internal/domain/entity"
	"jungle-test/app/internal/domain/services"
	"jungle-test/app/pkg/apperrors"
)

import sq "github.com/Masterminds/squirrel"

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const usersTableName = "users"

var _ = services.UserStorage(UserStorage{})

type UserStorage struct {
	client *pgxpool.Pool
}

func NewUserStorage(client *pgxpool.Pool) *UserStorage {
	return &UserStorage{client: client}
}

func (s UserStorage) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	sql, args, err := psql.
		Select("user_id", "password_hash", "created_at").
		From(usersTableName).Where(sq.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, apperrors.NewInternal("create squirrel query", err)
	}

	user := &entity.User{
		Username: username,
	}

	err = s.client.QueryRow(ctx, sql, args).Scan(&user.ID, &user.PasswordHash, &user.Created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.NewInternal("client.QueryRow", err)
	}

	return user, nil
}

const imageListTableName = "images"

var _ = services.ImageListStorage(ImageListStorage{})

type ImageListStorage struct {
	client *pgxpool.Pool
}

func NewImageListStorage(client *pgxpool.Pool) *ImageListStorage {
	return &ImageListStorage{client: client}
}

func (s ImageListStorage) AddImage(ctx context.Context, userID uuid.UUID, image entity.Image) error {
	sql, args, err := psql.
		Insert(imageListTableName).
		Columns("user_id", "name", "content_type", "size", "url", "created_at").
		Values(userID, image.Name, image.ContentType, image.Size, image.URL, image.Created).
		ToSql()
	if err != nil {
		return apperrors.NewInternal("create squirrel query", err)
	}

	_, err = s.client.Exec(ctx, sql, args...)
	if err != nil {
		return apperrors.NewInternal("client.Exec", err)
	}

	return nil
}

func (s ImageListStorage) GetImages(ctx context.Context, userID uuid.UUID) (images []entity.Image, err error) {
	sql, args, err := psql.
		Select("name", "content_type", "size", "url", "created_at").
		From(imageListTableName).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, apperrors.NewInternal("create squirrel query", err)
	}

	rows, err := s.client.Query(ctx, sql, args)
	defer rows.Close()

	var image entity.Image
	for rows.Next() {
		err = rows.Scan(&image.Name, &image.ContentType, &image.Size, &image.URL, &image.Created)
		if err != nil {
			return nil, apperrors.NewInternal("scan image", err)
		}
		images = append(images, image)
	}

	return images, err
}
