package database

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

type VectorDatabase struct {
	client *qdrant.Client
}

type VectorPoint struct {
	ID      uint64
	Vector  []float64 // Изменили на float64
	Payload map[string]interface{}
}

type SearchResult struct {
	ID      uint64
	Score   float32
	Payload map[string]interface{}
}

func NewVectorDatabase(host string, port int, apiKey string) (*VectorDatabase, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   host,
		Port:                   port,
		UseTLS:                 true,
		SkipCompatibilityCheck: true,
		APIKey:                 apiKey,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	return &VectorDatabase{
		client: client,
	}, nil
}

func (v *VectorDatabase) CreateCollection(name string, vectorSize int) error {
	err := v.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return fmt.Errorf("failed to create collection %s: %w", name, err)
	}

	return nil
}

func (v *VectorDatabase) DeleteCollection(name string) error {
	err := v.client.DeleteCollection(context.Background(), name)

	if err != nil {
		return fmt.Errorf("failed to delete collection %s: %w", name, err)
	}

	return nil
}

func (v *VectorDatabase) DeletePointsByFilter(collectionName string, filter map[string]interface{}) error {
	// Создаем фильтр в формате Qdrant с использованием правильной структуры
	qdrantFilter := &qdrant.Filter{
		Must: make([]*qdrant.Condition, 0, len(filter)),
	}

	// Преобразуем каждую пару ключ-значение в условие
	for key, value := range filter {
		var condition *qdrant.Condition

		switch v := value.(type) {
		case string:
			condition = &qdrant.Condition{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: key,
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Keyword{
								Keyword: v,
							},
						},
					},
				},
			}
		case int:
			condition = &qdrant.Condition{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: key,
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Integer{
								Integer: int64(v),
							},
						},
					},
				},
			}
		case int64:
			condition = &qdrant.Condition{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: key,
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Integer{
								Integer: v,
							},
						},
					},
				},
			}
		case float64:
			// Для float64 преобразуем в целое, если возможно
			if v == float64(int64(v)) {
				condition = &qdrant.Condition{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: key,
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Integer{
									Integer: int64(v),
								},
							},
						},
					},
				}
			} else {
				return fmt.Errorf("float64 values that are not integers are not supported in filter for key %s", key)
			}
		case bool:
			condition = &qdrant.Condition{
				ConditionOneOf: &qdrant.Condition_Field{
					Field: &qdrant.FieldCondition{
						Key: key,
						Match: &qdrant.Match{
							MatchValue: &qdrant.Match_Boolean{
								Boolean: v,
							},
						},
					},
				},
			}
		default:
			return fmt.Errorf("unsupported filter type for key %s", key)
		}

		qdrantFilter.Must = append(qdrantFilter.Must, condition)
	}

	// Удаляем точки с использованием созданного фильтра
	_, err := v.client.Delete(context.Background(), &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Filter{
				Filter: qdrantFilter,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete points from collection %s by filter %v: %w", collectionName, filter, err)
	}

	return nil
}

func (v *VectorDatabase) CollectionExists(name string) (bool, error) {
	collections, err := v.client.ListCollections(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to list collections: %w", err)
	}

	for _, collection := range collections {
		if collection == name {
			return true, nil
		}
	}

	return false, nil
}

func (v *VectorDatabase) UpsertPoints(collectionName string, points []VectorPoint) error {
	qdrantPoints := make([]*qdrant.PointStruct, len(points))

	for i, point := range points {
		payload := make(map[string]*qdrant.Value)
		for key, value := range point.Payload {
			switch v := value.(type) {
			case string:
				payload[key] = &qdrant.Value{
					Kind: &qdrant.Value_StringValue{StringValue: v},
				}
			case int:
				payload[key] = &qdrant.Value{
					Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(v)},
				}
			case int64:
				payload[key] = &qdrant.Value{
					Kind: &qdrant.Value_IntegerValue{IntegerValue: v},
				}
			case float64:
				payload[key] = &qdrant.Value{
					Kind: &qdrant.Value_DoubleValue{DoubleValue: v},
				}
			case bool:
				payload[key] = &qdrant.Value{
					Kind: &qdrant.Value_BoolValue{BoolValue: v},
				}
			}
		}

		// Конвертируем float64 в float32 для Qdrant
		vector32 := make([]float32, len(point.Vector))
		for j, val := range point.Vector {
			vector32[j] = float32(val)
		}

		qdrantPoints[i] = &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{
					Num: point.ID,
				},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{Data: vector32},
				},
			},
			Payload: payload,
		}
	}

	_, err := v.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         qdrantPoints,
	})

	if err != nil {
		return fmt.Errorf("failed to upsert points to collection %s: %w", collectionName, err)
	}

	return nil
}

func (v *VectorDatabase) Search(collectionName string, queryVector []float64, limit int, scoreThreshold float32) ([]SearchResult, error) {
	// Конвертируем float64 в float32 для Qdrant
	queryVector32 := make([]float32, len(queryVector))
	for i, val := range queryVector {
		queryVector32[i] = float32(val)
	}

	limitUint64 := uint64(limit)
	searchResult, err := v.client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: collectionName,
		Limit:          &limitUint64,
		ScoreThreshold: &scoreThreshold,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search in collection %s: %w", collectionName, err)
	}

	results := make([]SearchResult, len(searchResult))
	for i, point := range searchResult {
		payload := make(map[string]interface{})
		for key, value := range point.Payload {
			switch v := value.Kind.(type) {
			case *qdrant.Value_StringValue:
				payload[key] = v.StringValue
			case *qdrant.Value_IntegerValue:
				payload[key] = v.IntegerValue
			case *qdrant.Value_DoubleValue:
				payload[key] = v.DoubleValue
			case *qdrant.Value_BoolValue:
				payload[key] = v.BoolValue
			}
		}

		results[i] = SearchResult{
			ID:      point.Id.GetNum(),
			Score:   point.Score,
			Payload: payload,
		}
	}

	return results, nil
}

func (v *VectorDatabase) DeletePoints(collectionName string, pointIDs []uint64) error {
	ids := make([]*qdrant.PointId, len(pointIDs))
	for i, id := range pointIDs {
		ids[i] = &qdrant.PointId{
			PointIdOptions: &qdrant.PointId_Num{Num: id},
		}
	}

	_, err := v.client.Delete(context.Background(), &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{Ids: ids},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete points from collection %s: %w", collectionName, err)
	}

	return nil
}

func (v *VectorDatabase) CountPoints(collectionName string) (uint64, error) {
	info, err := v.client.GetCollectionInfo(context.Background(), collectionName)

	if err != nil {
		return 0, fmt.Errorf("failed to get collection info for %s: %w", collectionName, err)
	}

	return *info.PointsCount, nil
}

func (v *VectorDatabase) GetCollectionInfo(collectionName string) (map[string]interface{}, error) {
	info, err := v.client.GetCollectionInfo(context.Background(), collectionName)

	if err != nil {
		return nil, fmt.Errorf("failed to get collection info for %s: %w", collectionName, err)
	}

	return map[string]interface{}{
		"points_count":          info.PointsCount,
		"segments_count":        info.SegmentsCount,
		"status":                info.Status.String(),
		"indexed_vectors_count": info.IndexedVectorsCount,
	}, nil
}

func (v *VectorDatabase) Close() error {
	// В текущей версии go-client нет явного метода Close
	return nil
}
