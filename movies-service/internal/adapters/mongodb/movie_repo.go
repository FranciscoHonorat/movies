package mongodb

import (
	"context"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoMovieRepository struct {
	collection *mongo.Collection
}

type movieDocument struct {
	ID    int32  `bson:"_id"`
	Title string `bson:"title"`
	Year  string `bson:"year"`
}

func NewMongoRepository(collection *mongo.Collection) *MongoMovieRepository {
	return &MongoMovieRepository{
		collection: collection,
	}
}

func toDocument(movie *domain.Movie) movieDocument {
	return movieDocument{
		ID:    movie.ID,
		Title: movie.Title,
		Year:  movie.Year,
	}
}

func toDomain(doc movieDocument) *domain.Movie {
	return &domain.Movie{
		ID:    doc.ID,
		Title: doc.Title,
		Year:  doc.Year,
	}
}

func (m *MongoMovieRepository) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	filters := bson.M{"_id": id}

	var doc movieDocument
	err := m.collection.FindOne(ctx, filters).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrMovieNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return toDomain(doc), nil
}

func (m *MongoMovieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	doc := toDocument(movie)
	_, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return toDomain(doc), nil
}

func (m *MongoMovieRepository) DeleteMovie(ctx context.Context, id int32) error {
	filter := bson.M{"_id": id}

	result, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return domain.ErrInternalServer
	}
	if result.DeletedCount == 0 {
		return domain.ErrMovieNotFound
	}
	return nil
}

func (m *MongoMovieRepository) Count(ctx context.Context, filters output.ListFilters) (int32, error) {
	filter := bson.M{}

	if filters.Title != "" {
		filter["title"] = filters.Title
	}

	if filters.Year != "" {
		filter["year"] = filters.Year
	}

	count, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, domain.ErrInternalServer
	}
	return int32(count), nil
}

func (m *MongoMovieRepository) ListMovies(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
	filter := bson.M{}
	if filters.Title != "" {
		filter["title"] = filters.Title
	}

	if filters.Year != "" {
		filter["year"] = filters.Year
	}

	opts := options.Find()
	opts.SetLimit(int64(pagination.Limit))
	opts.SetSkip(int64((pagination.Page - 1) * pagination.Limit))

	if sorting.By != "" {
		opts.SetSort(bson.D{bson.E{Key: sorting.By, Value: 1}})
	}

	var docs []movieDocument
	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, domain.ErrInternalServer
	}

	var movies []domain.Movie
	for _, doc := range docs {
		movies = append(movies, *toDomain(doc))
	}
	return movies, nil
}
