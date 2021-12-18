package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
// const hostUrl string = "http://127.0.0.1:9200"

var client *elastic.Client
var index string
var dataPool []string
var DataChannel chan string

func init() {
	dataPool = []string{}
	DataChannel = make(chan string, 1000)
}

func BulkAdd(data []string) error {
	log.Debugf("开始批量导入 %d", len(data))
	bulkRequest := client.Bulk()
	for n, x := range data {
		log.Debugf("%s %d 装载批量弹药 %d Pool %d", fmt.Sprintf("ABC%d%d", time.Now().UnixNano(), n), n, len(DataChannel), len(dataPool))
		req := elastic.NewBulkIndexRequest().Index(index).Type("doc").Id(fmt.Sprintf("ABC%d%d", time.Now().UnixNano(), n)).Doc(x)
		bulkRequest = bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}
	if bulkResponse != nil {
		indexed := bulkResponse.Indexed()
		if len(indexed) != 1 {
			log.Debug("indexed is not length 1")
		}
		if indexed[0].Status != 201 {
			log.Error("Status ", indexed[0].Status)
		}
	}
	return nil
}

func InitEs(hostUrl, indexName string) {
	go func() {
		log.Info("开启发送通道")
		var mutex sync.Mutex
		for {
			mutex.Lock()
			select {
			case v := <-DataChannel:
				log.Debugf("新增数据Channel %d 数据池 %d", len(DataChannel), len(dataPool))
				dataPool = append(dataPool, v)
				mutex.Unlock()
			case <-time.After(time.Nanosecond * 100):
				if len(dataPool) != 0 {
					if len(dataPool) > 100 {
						log.Debugf("任务数 %d 执行", len(dataPool))
						err := BulkAdd(dataPool)
						if err != nil {
							log.Error(err.Error())
						} else {
							dataPool = []string{}
							log.Debug("====> 完成批量上传 清空数据池")
						}
					}
				}
				mutex.Unlock()
			}
		}
	}()
	// errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	// client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(hostUrl))
	client, err = elastic.NewClient(elastic.SetURL(hostUrl))
	if err != nil {
		log.Fatal(err.Error())
	}
	info, code, err := client.Ping(hostUrl).Do(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Infof("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(hostUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Infof("Elasticsearch version %s\n", esversion)

	index = fmt.Sprintf("%s-%s", indexName, time.Now().Format("2006-01-02"))
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		// Handle error
		log.Fatal(err.Error())
		// return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			log.Fatal(err.Error())
			// return err
		}
		log.Infof("Add new Index %s\n", index)
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	} else {
		log.Infof("Index %s 已存在", index)
	}

	// 每一小时检查索引名是否更新达到按天更新的效果
	go func(name string) {
		d := time.Duration(time.Hour)
		t := time.NewTicker(d)
		defer t.Stop()

		for {
			<-t.C
			if index != fmt.Sprintf("%s-%s", name, time.Now().Format("2006-01-02")) {
				index = fmt.Sprintf("%s-%s", name, time.Now().Format("2006-01-02"))
				exists, err := client.IndexExists(index).Do(context.Background())
				if err != nil {
					// Handle error
					log.Fatal(err.Error())
					// return err
				}
				if !exists {
					// Create a new index.
					createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(context.Background())
					if err != nil {
						// Handle error
						log.Fatal(err.Error())
						// return err
					}
					log.Infof("Add new Index %s\n", index)
					if !createIndex.Acknowledged {
						// Not acknowledged
					}
				} else {
					log.Infof("Index %s 已存在", index)
				}
			}

		}
	}(indexName)
}

func CreateEs(data, typed, id string) error {
	put2, err := client.Index().Index(index).Type(typed).Id(id).BodyString(data).Do(context.Background())
	if err != nil {
		return err
	}
	log.Debugf("Indexed tweet %s to index %s, type %s data %s\n", put2.Id, put2.Index, put2.Type, data)
	return nil
}

func Search(index, typs string) {
	var res *elastic.SearchResult
	var err error

	res, err = client.Search(index).Type(typs).Do(context.Background())
	printResult(res, err)
}

func printResult(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}

	log.Debugln("length:", res.Hits.TotalHits)
}
