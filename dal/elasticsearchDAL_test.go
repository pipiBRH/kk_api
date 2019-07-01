package dal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olivere/elastic"
	database "github.com/pipiBRH/kk_database"
	"github.com/stretchr/testify/assert"
)

func MockEsDAL(url string) (*ElasticsearchDAL, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	esConn := &database.ElasticsearchConnection{
		Ctx:    context.Background(),
		Client: client,
	}

	return &ElasticsearchDAL{Es: esConn}, nil
}

func TestSearchYouBikeInfoByText(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"took": 287,
			"timed_out": false,
			"_shards": {
			  "total": 3,
			  "successful": 3,
			  "skipped": 0,
			  "failed": 0
			},
			"hits": {
			  "total": 8,
			  "max_score": 5.889098,
			  "hits": [
				{
				  "_index": "youbike",
				  "_type": "_doc",
				  "_id": "1003",
				  "_score": 5.889098,
				  "_source": {
					"Sno": 1003,
					"Sna": "汐止區公所",
					"Tot": 46,
					"Sbi": 14,
					"Sarea": "汐止區",
					"Mday": "2019-07-01 10:03:15",
					"Ar": "新台五路一段/仁愛路口(新台五路側汐止地政事務所前機車停車場)",
					"Sareaen": "Xizhi Dist.",
					"Snaen": "Xizhi Dist. Office",
					"Aren": "Sec. 1, Xintai 5th Rd./Ren’ai Rd.",
					"Bemp": 32,
					"Act": 1,
					"Location": [
					  121.6583,
					  25.064161
					]
				  }
				}
			  ]
			}
		  }`

		w.Write([]byte(resp))
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	s, err := MockEsDAL(ts.URL)
	assert.NoError(t, err)

	expectPagination := &Pagination{
		CurrentPage:  1,
		PageSize:     1,
		TotalPage:    8,
		TotalRecords: 8,
	}

	expextYoubikeInfo := []YouBikeInfo{
		YouBikeInfo{
			Sno:     1003,
			Sna:     "汐止區公所",
			Tot:     46,
			Sbi:     14,
			Sarea:   "汐止區",
			Mday:    "2019-07-01 10:03:15",
			Ar:      "新台五路一段/仁愛路口(新台五路側汐止地政事務所前機車停車場)",
			Sareaen: "Xizhi Dist.",
			Snaen:   "Xizhi Dist. Office",
			Aren:    "Sec. 1, Xintai 5th Rd./Ren’ai Rd.",
			Bemp:    32,
			Act:     1,
			Location: []float32{
				121.6583,
				25.064161,
			},
		},
	}

	youbikeInfo, pagination, err := s.SearchYouBikeInfoByText("youbike", "汐止", 1, 1)
	assert.NoError(t, err)

	assert.Equal(t, expextYoubikeInfo, youbikeInfo)
	assert.Equal(t, expectPagination, pagination)
}

func TestSearchTopThreeFreeSpaceYouBikeInfo(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"took": 142,
			"timed_out": false,
			"_shards": {
			  "total": 3,
			  "successful": 3,
			  "skipped": 0,
			  "failed": 0
			},
			"hits": {
			  "total": 556,
			  "max_score": 0,
			  "hits": []
			},
			"aggregations": {
			  "sareas": {
				"doc_count_error_upper_bound": 14,
				"sum_other_doc_count": 356,
				"buckets": [
				  {
					"key": "板橋區",
					"doc_count": 86
				  },
				  {
					"key": "三重區",
					"doc_count": 59
				  },
				  {
					"key": "新莊區",
					"doc_count": 56
				  }
				]
			  }
			}
		  }`

		w.Write([]byte(resp))
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	s, err := MockEsDAL(ts.URL)
	assert.NoError(t, err)

	expectTopThreeFreeSpaceBucket := []TopThreeFreeSpaceBucket{
		TopThreeFreeSpaceBucket{
			Key:      "板橋區",
			DocCount: 86,
		},
		TopThreeFreeSpaceBucket{
			Key:      "三重區",
			DocCount: 59,
		},
		TopThreeFreeSpaceBucket{
			Key:      "新莊區",
			DocCount: 56,
		},
	}

	topThreeFreeSpaceBucket, err := s.SearchTopThreeFreeSpaceYouBikeInfo("youbike")
	assert.NoError(t, err)
	assert.Equal(t, expectTopThreeFreeSpaceBucket, topThreeFreeSpaceBucket)
}

func TestSearchMaxAndMinFreeSpaceYouBikeInfoByTimeRange(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"took": 20,
			"timed_out": false,
			"_shards": {
			  "total": 3,
			  "successful": 3,
			  "skipped": 0,
			  "failed": 0
			},
			"hits": {
			  "total": 1102,
			  "max_score": 0,
			  "hits": []
			},
			"aggregations": {
			  "range": {
				"buckets": [
				  {
					"key_as_string": "2019-07-01 00:00:00",
					"key": 1561939200000,
					"doc_count": 1102,
					"sareas": {
					  "doc_count_error_upper_bound": 0,
					  "sum_other_doc_count": 236,
					  "buckets": [
						{
						  "key": "板橋區",
						  "doc_count": 170,
						  "bemp_stats": {
							"count": 170,
							"min": 7,
							"max": 99,
							"avg": 26.111764705882354,
							"sum": 4439
						  }
						},
						{
						  "key": "三重區",
						  "doc_count": 118,
						  "bemp_stats": {
							"count": 118,
							"min": 2,
							"max": 76,
							"avg": 22.64406779661017,
							"sum": 2672
						  }
						},
						{
						  "key": "新莊區",
						  "doc_count": 112,
						  "bemp_stats": {
							"count": 112,
							"min": 1,
							"max": 72,
							"avg": 22.098214285714285,
							"sum": 2475
						  }
						},
						{
						  "key": "中和區",
						  "doc_count": 86,
						  "bemp_stats": {
							"count": 86,
							"min": 7,
							"max": 71,
							"avg": 21.837209302325583,
							"sum": 1878
						  }
						},
						{
						  "key": "新店區",
						  "doc_count": 86,
						  "bemp_stats": {
							"count": 86,
							"min": 0,
							"max": 44,
							"avg": 22.325581395348838,
							"sum": 1920
						  }
						},
						{
						  "key": "樹林區",
						  "doc_count": 66,
						  "bemp_stats": {
							"count": 66,
							"min": 6,
							"max": 55,
							"avg": 21.393939393939394,
							"sum": 1412
						  }
						},
						{
						  "key": "汐止區",
						  "doc_count": 66,
						  "bemp_stats": {
							"count": 66,
							"min": 4,
							"max": 56,
							"avg": 23.37878787878788,
							"sum": 1543
						  }
						},
						{
						  "key": "土城區",
						  "doc_count": 62,
						  "bemp_stats": {
							"count": 62,
							"min": 3,
							"max": 74,
							"avg": 22.20967741935484,
							"sum": 1377
						  }
						},
						{
						  "key": "林口區",
						  "doc_count": 50,
						  "bemp_stats": {
							"count": 50,
							"min": 10,
							"max": 49,
							"avg": 20.36,
							"sum": 1018
						  }
						},
						{
						  "key": "永和區",
						  "doc_count": 50,
						  "bemp_stats": {
							"count": 50,
							"min": 3,
							"max": 33,
							"avg": 15.92,
							"sum": 796
						  }
						}
					  ]
					}
				  }
				]
			  }
			}
		  }`

		w.Write([]byte(resp))
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	s, err := MockEsDAL(ts.URL)
	assert.NoError(t, err)

	expectMaxAndMin := []MaxAndMinBucket{
		MaxAndMinBucket{
			Date: "2019-07-01 00:00:00",
			Sareas: Sareas{
				Buckets: []Sarea{
					Sarea{
						Sarea: "板橋區",
						Count: 170,
						Stats: Stats{
							Count: 170,
							Min:   7,
							Max:   99,
							Avf:   26.111764705882354,
							Sum:   4439,
						},
					},
					Sarea{
						Sarea: "三重區",
						Count: 118,
						Stats: Stats{
							Count: 118,
							Min:   2,
							Max:   76,
							Avf:   22.64406779661017,
							Sum:   2672,
						},
					},
					Sarea{
						Sarea: "新莊區",
						Count: 112,
						Stats: Stats{
							Count: 112,
							Min:   1,
							Max:   72,
							Avf:   22.098214285714285,
							Sum:   2475,
						},
					},
					Sarea{
						Sarea: "中和區",
						Count: 86,
						Stats: Stats{
							Count: 86,
							Min:   7,
							Max:   71,
							Avf:   21.837209302325583,
							Sum:   1878,
						},
					},
					Sarea{
						Sarea: "新店區",
						Count: 86,
						Stats: Stats{
							Count: 86,
							Min:   0,
							Max:   44,
							Avf:   22.325581395348838,
							Sum:   1920,
						},
					},
					Sarea{
						Sarea: "樹林區",
						Count: 66,
						Stats: Stats{
							Count: 66,
							Min:   6,
							Max:   55,
							Avf:   21.393939393939394,
							Sum:   1412,
						},
					},
					Sarea{
						Sarea: "汐止區",
						Count: 66,
						Stats: Stats{
							Count: 66,
							Min:   4,
							Max:   56,
							Avf:   23.37878787878788,
							Sum:   1543,
						},
					},
					Sarea{
						Sarea: "土城區",
						Count: 62,
						Stats: Stats{
							Count: 62,
							Min:   3,
							Max:   74,
							Avf:   22.20967741935484,
							Sum:   1377,
						},
					},
					Sarea{
						Sarea: "林口區",
						Count: 50,
						Stats: Stats{
							Count: 50,
							Min:   10,
							Max:   49,
							Avf:   20.36,
							Sum:   1018,
						},
					},
					Sarea{
						Sarea: "永和區",
						Count: 50,
						Stats: Stats{
							Count: 50,
							Min:   3,
							Max:   33,
							Avf:   15.92,
							Sum:   796,
						},
					},
				},
			},
		},
	}

	maxAndMin, err := s.SearchMaxAndMinFreeSpaceYouBikeInfoByTimeRange("youbike_history", "2019-07-01 00:00:00", "2019-07-01 00:10:00")
	assert.NoError(t, err)
	assert.Equal(t, expectMaxAndMin, maxAndMin)
}
