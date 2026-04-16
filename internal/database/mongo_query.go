package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoFindFn = func(collection *mongo.Collection, ctx context.Context, filter interface{}, findOptions *options.FindOptions) (*mongo.Cursor, error) {
	return collection.Find(ctx, filter, findOptions)
}

var mongoCursorAllFn = func(cursor *mongo.Cursor, ctx context.Context, result interface{}) error {
	return cursor.All(ctx, result)
}

var mongoCursorCloseFn = func(cursor *mongo.Cursor, ctx context.Context) error {
	return cursor.Close(ctx)
}

// FindAll executes a Mongo find query and decodes all documents to []T.
// It keeps repository behavior consistent by returning an empty slice on
// query/decoding failures and invalidating the connection on connectivity errors.
func FindAll[T any](
	ctx context.Context,
	collection *mongo.Collection,
	filter interface{},
	findOptions *options.FindOptions,
	invalidateConnection func(),
) ([]T, error) {
	if collection == nil {
		return []T{}, nil
	}

	cursor, err := mongoFindFn(collection, ctx, filter, findOptions)
	if err != nil {
		if IsMongoConnectionError(err) && invalidateConnection != nil {
			invalidateConnection()
		}
		return []T{}, nil
	}
	defer func() {
		_ = mongoCursorCloseFn(cursor, ctx)
	}()

	var result []T
	if err := mongoCursorAllFn(cursor, ctx, &result); err != nil {
		if IsMongoConnectionError(err) && invalidateConnection != nil {
			invalidateConnection()
		}
		return []T{}, nil
	}

	if result == nil {
		return []T{}, nil
	}

	return result, nil
}
