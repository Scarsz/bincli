package bin

import (
	b64 "encoding/base64"
	"errors"
	"github.com/Scarsz/bincli/crypto"
	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"time"
)

type Bin struct {
	UUID        uuid.UUID
	Key         string
	Hits        int
	Description string
	Expiration  int
	Timestamp   int64
	Files       []File
}

type Options struct {
	Description string
	Expiration  int
	Files       []File
}

func (bin *Bin) URL() string {
	if bin.Key != "" {
		return "https://bin.scarsz.me/" + bin.UUID.String() + "#" + bin.Key
	} else {
		return "https://bin.scarsz.me/" + bin.UUID.String()
	}
}

func (bin *Bin) SaveToTemp() string {
	dir, err := ioutil.TempDir("", "bin-"+bin.UUID.String()+"-*")
	if err != nil {
		panic(err)
	}
	bin.Save(dir)
	return dir
}

func (bin *Bin) Save(path string) {
	for _, file := range bin.Files {
		if !file.Available() {
			continue
		}

		err := ioutil.WriteFile(path+"/"+file.Name, file.Content, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func Create(options Options) Bin {
	key := uniuri.NewLen(32)
	keyBytes := []byte(key)

	payload := make(map[string]interface{})

	payload["expiration"] = options.Expiration

	if options.Description != "" {
		payload["description"] = crypto.EncryptString(keyBytes, options.Description)
	}

	filesPayload := make([]map[string]interface{}, 0, len(options.Files))
	for _, file := range options.Files {
		filesPayload = append(filesPayload, file.SerializeMap(keyBytes))
	}
	payload["files"] = filesPayload

	//payloadJson, err := json.Marshal(payload)
	//fmt.Println(string(payloadJson))

	resp, err := req.Post("https://bin.scarsz.me/v1/post", req.BodyJSON(payload))
	if err != nil {
		panic(err)
	}
	json := resp.String()

	return Bin{
		UUID:        uuid.MustParse(gjson.Get(json, "bin").String()),
		Key:         key,
		Hits:        0,
		Description: options.Description,
		Expiration:  options.Expiration,
		Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		Files:       options.Files,
	}
}

func Retrieve(id uuid.UUID, key string) (Bin, error) {
	resp, err := req.Get("https://bin.scarsz.me/" + id.String() + ".json")
	if err != nil {
		panic(err)
	}
	json := resp.String()

	if resp.Response().StatusCode == 404 {
		return Bin{}, errors.New("Bin " + id.String() + " doesn't exist (404)")
	}

	var descriptionBytes []byte
	if gjson.Get(json, "description").Exists() {
		descriptionBytes, err = b64.StdEncoding.DecodeString(gjson.Get(json, "description").String())
		if err != nil {
			return Bin{}, err
		}
	}

	b := Bin{
		UUID:       id,
		Key:        key,
		Hits:       int(gjson.Get(json, "hits").Int()),
		Expiration: int(gjson.Get(json, "hits").Int()),
		Timestamp:  gjson.Get(json, "time").Int(),
	}
	if len(key) == 32 {
		if len(descriptionBytes) > 0 {
			b.Description = string(crypto.Decrypt([]byte(key), descriptionBytes))
		}
		for _, fileMap := range gjson.Get(json, "files").Array() {
			b.Files = append(b.Files, FileFromEncryptedMap(fileMap.Map(), key))
		}
	} else {
		for _, fileMap := range gjson.Get(json, "files").Array() {
			b.Files = append(b.Files, File{UUID: uuid.MustParse(fileMap.Map()["id"].String())})
		}
	}

	return b, nil
}
