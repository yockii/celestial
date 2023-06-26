package search

import (
	"encoding/json"
	"errors"
	"github.com/meilisearch/meilisearch-go"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
)

var client *meilisearch.Client
var index = "celestial"

func InitMeiliSearch(host, apiKey string, idx ...string) error {
	c := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: apiKey,
	})
	stats, err := c.GetStats()
	if err != nil {
		logger.Errorln(err)
		return err
	}
	client = c
	logger.Info("search database: ", stats.DatabaseSize)
	if len(idx) > 0 {
		index = idx[0]
	}
	return nil
}

func DeleteDocuments(idList ...uint64) (success bool, err error) {
	if client == nil {
		err = errors.New("No Search client initialized. ")
		return
	}

	for _, id := range idList {
		var taskInfo *meilisearch.TaskInfo
		taskInfo, err = client.Index(index).DeleteDocument(strconv.FormatUint(id, 10))
		if err != nil {
			logger.Errorln(err)
			return
		}
		if taskInfo != nil {
			success = success && taskInfo.Status != meilisearch.TaskStatusFailed
		}
	}
	return
}

func AddDocument(documents []*Document) (success bool, err error) {
	if client == nil {
		err = errors.New("No Search client initialized. ")
		return
	}
	var task *meilisearch.TaskInfo
	task, err = client.Index(index).AddDocuments(documents, "id")
	if err != nil {
		logger.Errorln(err)
		return
	}
	if task != nil {
		success = task.Status != meilisearch.TaskStatusFailed
	}
	return
}

func UpdateDocument(documents []*Document) (success bool, err error) {
	if client == nil {
		err = errors.New("No Search client initialized. ")
		return
	}
	var task *meilisearch.TaskInfo
	task, err = client.Index(index).UpdateDocuments(documents, "id")
	if err != nil {
		logger.Errorln(err)
		return
	}
	if task != nil {
		success = task.Status != meilisearch.TaskStatusFailed
	}
	return
}

func Search(keyword string, limit, offset int64) (total int64, result []*Document, err error) {
	if client == nil {
		err = errors.New("No Search client initialized. ")
		return
	}
	var searchResp *meilisearch.SearchResponse
	searchResp, err = client.Index(index).Search(keyword, &meilisearch.SearchRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	total = searchResp.EstimatedTotalHits
	for _, hit := range searchResp.Hits {
		var bs []byte
		bs, err = json.Marshal(hit)
		j := gjson.ParseBytes(bs)
		result = append(result, &Document{
			ID:         j.Get("id").Uint(),
			Title:      j.Get("title").String(),
			Content:    j.Get("content").String(),
			Route:      j.Get("route").String(),
			CreateTime: j.Get("createTime").Int(),
			UpdateTime: j.Get("updateTime").Int(),
		})
	}
	return
}
