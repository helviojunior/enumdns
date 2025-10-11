package writers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	//"path/filepath"
	//"reflect"
	//"io"

	"github.com/bob-reis/enumdns/internal/tools"
	logger "github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/models"
	elk "github.com/elastic/go-elasticsearch/v8"
	esapi "github.com/elastic/go-elasticsearch/v8/esapi"
)

// fields in the main model to ignore
var elkExludedFields = []string{"network"}

const (
	elasticDocumentError   = "cannot create/update document"
	parseResponseBodyError = "failure to parse response body: %s"
)

// JsonWriter is a JSON lines writer
type ElasticWriter struct {
	Client         *elk.Client
	Index          string
	screenshotPath string
}

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

type indexResponse struct {
	ID     string `json:"_id"`
	Index  string `json:"_index"`
	Result string `json:"result"`
	Error  struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
		Cause  struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"caused_by"`
	} `json:"error"`
}

// NewJsonWriter return a new Json lines writer
func NewElasticWriter(uri string) (*ElasticWriter, error) {

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	port := u.Port()
	if port == "" {
		port = "9200"
	}
	index_name := u.EscapedPath()
	index_name = strings.Trim(index_name, "/ ")
	index_name = strings.SplitN(index_name, "/", 2)[0]
	if index_name == "" {
		index_name = "enumdns"
	}

	wr := &ElasticWriter{
		Index:          index_name,
		screenshotPath: "",
	}

	conf := elk.Config{
		Addresses: []string{
			fmt.Sprintf("%s://%s:%s/", u.Scheme, u.Hostname(), port),
		},
		//Username: username,
		//Password: password,
		//CACert:   cert,
		RetryOnStatus: []int{429, 502, 503, 504},
		RetryBackoff: func(i int) time.Duration {
			// A simple exponential delay
			d := time.Duration(math.Exp2(float64(i))) * time.Second
			logger.Debugf("Elastic retry, attempt: %d | Sleeping for %s...\n", i, d)
			return d
		},
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
		},
	}

	if username != "" && password != "" {
		conf.Username = username
		conf.Password = password
	}

	wr.Client, err = elk.NewClient(conf)
	if err != nil {
		return nil, err
	}

	//File Index
	err = wr.CreateIndex(wr.Index, `{
            "settings": {
                    "number_of_replicas": 1,
                    "index": {"highlight.max_analyzed_offset": 10000000}
                },

            "mappings": {
                "properties": {
                    "probed_at": {"type": "date"},
                    "test_id": {"type": "keyword"},
                    "fqdn": {"type": "keyword"},
                    "result_type": {"type": "keyword"},
                    "ipv4": {"type": "keyword"},
                    "ipv6": {"type": "keyword"},
                    "target": {"type": "keyword"},
                    "ptr": {"type": "keyword"},
                    "failed": {"type": "keyword"},
                    "failed_reason": {"type": "text"}
                }
            }
        }`)
	if err != nil {
		return nil, err
	}

	return wr, nil
}

// Write JSON lines to a file
func (ew *ElasticWriter) Write(result *models.Result) error {

	if !result.Exists {
		return nil
	}

	//File
	b_data, err := json.Marshal(*result) //ew.Marshal(*result)
	if err != nil {
		return err
	}

	res, err := ew.Client.Index(ew.Index, bytes.NewReader(b_data), ew.Client.Index.WithDocumentID(result.GetHash()))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		fmt.Printf("Err: %s", res)
		return errors.New(elasticDocumentError)
	}

	return nil
}

func (ew *ElasticWriter) CreateIndex(index string, mapping string) error {

	var raw map[string]interface{}

	response, err := ew.Client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {

		if response.StatusCode == 404 {
			indexReq := esapi.IndicesCreateRequest{
				Index: index,
				Body:  strings.NewReader(string(mapping)),
			}

			logger.Infof("Creating elastic index %s", index)
			res, err := indexReq.Do(context.Background(), ew.Client)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.IsError() {

				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					return fmt.Errorf(parseResponseBodyError, err)
				} else {
					return fmt.Errorf("cannot create/update elastic index [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}

			}

		} else {

			if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
				return fmt.Errorf(parseResponseBodyError, err)
			} else {
				return fmt.Errorf("cannot get elastic index [%d] %s: %s",
					response.StatusCode,
					raw["error"].(map[string]interface{})["type"],
					raw["error"].(map[string]interface{})["reason"],
				)
			}

		}

	}

	return nil

}

func (ew *ElasticWriter) CreateDocBulk(index string, docs map[string][]byte) error {
	var raw map[string]interface{}
	var buf bytes.Buffer
	size := 0
	for id, doc := range docs {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, id, "\n"))
		data := []byte(doc)
		data = append(data, "\n"...)

		size += len(meta) + len(data)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

	}

	logger.Debugf("Elastic bulk %d docs, %d bytes", len(docs), size)

	for i := 0; i < 10; i++ {

		res, err := ew.Client.Bulk(bytes.NewReader(buf.Bytes()), ew.Client.Bulk.WithIndex(index))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {

			if i >= 5 {
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					return fmt.Errorf(parseResponseBodyError, err)
				} else {
					return fmt.Errorf("error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}

			}

			// A successful response might still contain errors for particular documents...
			//
		} else {
			var blk *bulkResponse
			if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
				return fmt.Errorf(parseResponseBodyError, err)
			} else {
				for _, d := range blk.Items {
					// ... so for any HTTP status above 201 ...
					//
					if d.Index.Status > 201 {
						// ... and print the response status and error information ...
						logger.Debugf("  Error: [%d]: %s: %s: %s: %s",
							d.Index.Status,
							d.Index.Error.Type,
							d.Index.Error.Reason,
							d.Index.Error.Cause.Type,
							d.Index.Error.Cause.Reason,
						)
					} else {
						// ... otherwise increase the success counter.
						//

					}
				}
			}
		}

		if res.StatusCode == 200 || res.StatusCode == 201 {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New(elasticDocumentError)
}

func (ew *ElasticWriter) CreateDoc(index string, data []byte, doc_id string) error {
	var raw map[string]interface{}
	for i := 0; i < 10; i++ {
		res, err := ew.Client.Index(index, bytes.NewReader(data), ew.Client.Index.WithDocumentID(doc_id))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {

			if i >= 5 {
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					return fmt.Errorf(parseResponseBodyError, err)
				} else {
					return fmt.Errorf("error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}

			}

			// A successful response might still contain errors for particular documents...
			//
		} else {

			if res.StatusCode == 200 || res.StatusCode == 201 {
				return nil
			}

			//bodyBytes, err := io.ReadAll(res.Body)
			//if err != nil {
			//    return err
			//}
			//bodyString := string(bodyBytes)
			//fmt.Printf("Resp: %s", bodyString)

			var idxRes *indexResponse

			if err := json.NewDecoder(res.Body).Decode(&idxRes); err != nil {
				return fmt.Errorf(parseResponseBodyError, err)
			} else {
				//Debug result
			}
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New(elasticDocumentError)
}

func (ew *ElasticWriter) MarshalAppend(marshalled []byte, new_data map[string]interface{}) ([]byte, error) {
	t_data := make(map[string]interface{})
	if err := json.Unmarshal(marshalled, &t_data); err != nil {
		return marshalled, err
	}

	data := make(map[string]interface{})
	for k, v := range t_data {
		// skip excluded fields
		if tools.SliceHasStr(elkExludedFields, k) {
			continue
		}

		data[k] = v
	}

	for k, v := range new_data {
		data[k] = v
	}

	j_data, err := json.Marshal(data)
	if err != nil {
		return marshalled, err
	}

	return j_data, nil
}

func (ew *ElasticWriter) Marshal(v any) ([]byte, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}

	t_data := make(map[string]interface{})
	if err = json.Unmarshal(j, &t_data); err != nil {
		return []byte{}, err
	}

	data := make(map[string]interface{})
	for k, v := range t_data {
		// skip excluded fields
		if tools.SliceHasStr(elkExludedFields, k) {
			continue
		}

		data[k] = v
	}

	j_data, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	return j_data[:], nil
}

func (ew *ElasticWriter) WriteFqdn(result *models.FQDNData) error {
	return nil
}

func (ew *ElasticWriter) Finish() error {
	return nil
}
