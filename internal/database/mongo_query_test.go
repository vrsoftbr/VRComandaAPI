package database

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type queryTestModel struct {
	Value string
}

func resetMongoQueryFns() {
	mongoFindFn = func(collection *mongo.Collection, ctx context.Context, filter interface{}, findOptions *options.FindOptions) (*mongo.Cursor, error) {
		return collection.Find(ctx, filter, findOptions)
	}

	mongoCursorAllFn = func(cursor *mongo.Cursor, ctx context.Context, result interface{}) error {
		return cursor.All(ctx, result)
	}

	mongoCursorCloseFn = func(cursor *mongo.Cursor, ctx context.Context) error {
		return cursor.Close(ctx)
	}
}

func TestFindAllReturnsEmptyWhenCollectionIsNil(t *testing.T) {
	result, err := FindAll[queryTestModel](context.Background(), nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
}

func TestFindAllReturnsEmptyAndInvalidatesOnFindConnectionError(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return nil, errors.New("server selection timeout")
	}

	invalidated := 0
	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, func() {
		invalidated++
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
	if invalidated != 1 {
		t.Fatalf("expected invalidate callback once, got %d", invalidated)
	}
}

func TestFindAllReturnsEmptyWithoutInvalidatingOnFindNonConnectionError(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return nil, errors.New("validation failed")
	}

	invalidated := 0
	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, func() {
		invalidated++
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
	if invalidated != 0 {
		t.Fatalf("expected invalidate callback zero times, got %d", invalidated)
	}
}

func TestFindAllReturnsEmptyAndInvalidatesOnAllConnectionError(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return &mongo.Cursor{}, nil
	}
	mongoCursorAllFn = func(_ *mongo.Cursor, _ context.Context, _ interface{}) error {
		return errors.New("topology changed")
	}
	mongoCursorCloseFn = func(_ *mongo.Cursor, _ context.Context) error {
		return nil
	}

	invalidated := 0
	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, func() {
		invalidated++
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
	if invalidated != 1 {
		t.Fatalf("expected invalidate callback once, got %d", invalidated)
	}
}

func TestFindAllReturnsEmptyWithoutInvalidatingOnAllNonConnectionError(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return &mongo.Cursor{}, nil
	}
	mongoCursorAllFn = func(_ *mongo.Cursor, _ context.Context, _ interface{}) error {
		return errors.New("decode failed")
	}
	mongoCursorCloseFn = func(_ *mongo.Cursor, _ context.Context) error {
		return nil
	}

	invalidated := 0
	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, func() {
		invalidated++
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
	if invalidated != 0 {
		t.Fatalf("expected invalidate callback zero times, got %d", invalidated)
	}
}

func TestFindAllReturnsEmptyWhenDecodedResultIsNil(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return &mongo.Cursor{}, nil
	}
	mongoCursorAllFn = func(_ *mongo.Cursor, _ context.Context, result interface{}) error {
		return nil
	}
	mongoCursorCloseFn = func(_ *mongo.Cursor, _ context.Context) error {
		return nil
	}

	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
}

func TestFindAllReturnsDecodedItems(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return &mongo.Cursor{}, nil
	}
	mongoCursorAllFn = func(_ *mongo.Cursor, _ context.Context, result interface{}) error {
		typed, ok := result.(*[]queryTestModel)
		if !ok {
			t.Fatalf("expected *[]queryTestModel, got %T", result)
		}
		*typed = []queryTestModel{{Value: "A"}, {Value: "B"}}
		return nil
	}
	mongoCursorCloseFn = func(_ *mongo.Cursor, _ context.Context) error {
		return nil
	}

	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0].Value != "A" || result[1].Value != "B" {
		t.Fatalf("unexpected decoded result: %+v", result)
	}
}

func TestFindAllWithNilInvalidatorOnConnectionError(t *testing.T) {
	defer resetMongoQueryFns()

	mongoFindFn = func(_ *mongo.Collection, _ context.Context, _ interface{}, _ *options.FindOptions) (*mongo.Cursor, error) {
		return nil, errors.New("connection reset by peer")
	}

	result, err := FindAll[queryTestModel](context.Background(), &mongo.Collection{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got len=%d", len(result))
	}
}
