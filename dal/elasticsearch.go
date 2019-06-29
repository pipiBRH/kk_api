package dal

import (
	"encoding/json"
	"math"

	"github.com/olivere/elastic"
	database "github.com/pipiBRH/kk_database"
)

type ElasticsearchDAL struct {
	Es *database.ElasticsearchConnection
}

func NewElasticsearchDAL(client *database.ElasticsearchConnection) *ElasticsearchDAL {
	return &ElasticsearchDAL{
		Es: client,
	}
}

type Pagination struct {
	CurrentPage  int
	PageSize     int
	TotalPage    int
	TotalRecords int
}

type YouBikeInfo struct {
	Sno      int
	Sna      string
	Tot      int
	Sbi      int
	Sarea    string
	Mday     string
	Ar       string
	Sareaen  string
	Snaen    string
	Aren     string
	Bemp     int
	Act      int
	Location []float32
}

func (es *ElasticsearchDAL) SearchYouBikeInfoByText(index, text string, page, pageSize int) ([]YouBikeInfo, *Pagination, error) {
	offset := (page - 1) * pageSize
	query := elastic.NewMultiMatchQuery(text, "Ar", "Aren", "Sna", "Snaen")
	searchResult, err := es.Es.Client.Search().
		Index(index).
		Type("_doc").
		Query(query).
		From(offset).
		Size(pageSize).
		Do(es.Es.Ctx)
	if err != nil {
		return nil, nil, err
	}

	results := make([]YouBikeInfo, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		var temp YouBikeInfo
		err := json.Unmarshal(*hit.Source, &temp)
		if err != nil {
			return nil, nil, err
		}
		results[i] = temp
	}

	pagination := &Pagination{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPage:    int(math.Ceil(float64(searchResult.Hits.TotalHits) / float64(pageSize))),
		TotalRecords: int(searchResult.Hits.TotalHits),
	}

	return results, pagination, nil
}

type (
	TopThreeFreeSpaceBuckets struct {
		Buckets []TopThreeFreeSpaceBucket `json:"buckets"`
	}

	TopThreeFreeSpaceBucket struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	}
)

func (es *ElasticsearchDAL) SearchTopThreeFreeSpaceYouBikeInfo(index string) ([]TopThreeFreeSpaceBucket, error) {
	aggs := elastic.NewTermsAggregation().Field("Sarea").OrderByCountDesc().Size(3)
	searchResult, err := es.Es.Client.Search().
		Index(index).
		Type("_doc").
		Aggregation("sareas", aggs).
		Size(0).
		Do(es.Es.Ctx)
	if err != nil {
		return nil, err
	}

	var results TopThreeFreeSpaceBuckets
	err = json.Unmarshal(*searchResult.Aggregations["sareas"], &results)
	if err != nil {
		return nil, err
	}

	return results.Buckets, nil
}

type (
	MaxAndMinBuckets struct {
		Buckets []MaxAndMinBucket `json:"buckets"`
	}

	MaxAndMinBucket struct {
		Date   string `json:"key_as_string"`
		Sareas Sareas `json:"sareas"`
	}

	Sareas struct {
		Buckets []Sarea `json:"buckets"`
	}

	Sarea struct {
		Sarea string `json:"key"`
		Count int    `json:"doc_count"`
		Stats Stats  `json:"bemp_stats"`
	}

	Stats struct {
		Count int     `json:"count"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Avf   float64 `json:"avg"`
		Sum   float64 `json:"sum"`
	}
)

func (es *ElasticsearchDAL) SearchMaxAndMinFreeSpaceYouBikeInfoByTimeRange(index, from, to string) ([]MaxAndMinBucket, error) {
	query := elastic.NewRangeQuery("Mday").Gte(from).Lte(to)

	statsAggs := elastic.NewStatsAggregation().Field("Bemp")

	TermsAggs := elastic.NewTermsAggregation().
		Field("Sarea").
		OrderByCountDesc().
		SubAggregation("bemp_stats", statsAggs)

	dateHistogramAggs := elastic.NewDateHistogramAggregation().
		Field("Mday").
		Interval("10m").
		OrderByKeyAsc().
		SubAggregation("sareas", TermsAggs)

	searchResult, err := es.Es.Client.Search().
		Index(index).
		Type("_doc").
		Query(query).
		Aggregation("range", dateHistogramAggs).
		Size(0).
		Do(es.Es.Ctx)
	if err != nil {
		return nil, err
	}

	var results MaxAndMinBuckets
	err = json.Unmarshal(*searchResult.Aggregations["range"], &results)
	if err != nil {
		return nil, err
	}

	return results.Buckets, nil
}
