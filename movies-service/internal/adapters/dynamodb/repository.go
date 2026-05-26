package dynamodb

import (
	"context"
	"errors"
	"strconv"

	"github.com/FranciscoHonorat/movies/movies-service/internal/core/domain"
	"github.com/FranciscoHonorat/movies/movies-service/internal/core/ports/output"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "movies"

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoRepository(ctx context.Context, endpoint string) (*DynamoDBRepository, error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint}, nil
				},
			),
		),
	)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	var resourceInUse *types.ResourceInUseException
	if err != nil && !errors.As(err, &resourceInUse) {
		return nil, err
	}

	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoDBRepository) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(int(id))},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, domain.ErrMovieNotFound
	}
	return &domain.Movie{
		ID:    id,
		Title: result.Item["title"].(*types.AttributeValueMemberS).Value,
		Year:  result.Item["year"].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (r *DynamoDBRepository) CreateMovie(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]types.AttributeValue{
			"id":    &types.AttributeValueMemberN{Value: strconv.Itoa(int(movie.ID))},
			"title": &types.AttributeValueMemberS{Value: movie.Title},
			"year":  &types.AttributeValueMemberS{Value: movie.Year},
		},
	})
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (r *DynamoDBRepository) ListMovies(ctx context.Context, filters output.ListFilters, pagination output.Pagination, sorting output.Sorting) ([]domain.Movie, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}
	movies := make([]domain.Movie, 0, len(result.Items))
	for _, item := range result.Items {
		idInt, _ := strconv.ParseInt(item["id"].(*types.AttributeValueMemberN).Value, 10, 32)
		movies = append(movies, domain.Movie{
			ID:    int32(idInt),
			Title: item["title"].(*types.AttributeValueMemberS).Value,
			Year:  item["year"].(*types.AttributeValueMemberS).Value,
		})
	}
	return movies, nil
}

func (r *DynamoDBRepository) DeleteMovie(ctx context.Context, id int32) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(int(id))},
		},
	})
	if err != nil {
		return domain.ErrInternalServer
	}
	return nil
}

func (r *DynamoDBRepository) Count(ctx context.Context, filters output.ListFilters) (int32, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return 0, err
	}
	return int32(len(result.Items)), nil
}
