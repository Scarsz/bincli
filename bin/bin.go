package bin

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Scarsz/bincli/crypto"
	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"github.com/imroc/req"
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

	var response map[string]interface{}
	err = json.Unmarshal([]byte(resp.String()), &response)
	if err != nil {
		panic(err)
	}

	return Bin{
		UUID:        uuid.MustParse(response["bin"].(string)),
		Key:         key,
		Hits:        0,
		Description: options.Description,
		Expiration:  options.Expiration,
		Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		Files:       options.Files,
	}
}

func Retrieve(uuid uuid.UUID, key string) (Bin, error) {
	resp, err := req.Get("https://bin.scarsz.me/" + uuid.String() + ".json")
	if err != nil {
		panic(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(resp.String()), &response)
	if err != nil {
		return Bin{}, err
	}

	if resp.Response().StatusCode == 404 {
		return Bin{}, errors.New("Bin " + uuid.String() + " doesn't exist (404)")
	}

	var descriptionBytes []byte
	if response["description"] != nil {
		descriptionBytes, err = b64.StdEncoding.DecodeString(response["description"].(string))
		if err != nil {
			return Bin{}, err
		}
	}

	b := Bin{
		UUID:       uuid,
		Key:        key,
		Hits:       int(response["hits"].(float64)),
		Expiration: int(response["expiration"].(float64)),
		Timestamp:  int64(response["time"].(float64)),
		Files:      nil,
	}
	if len(descriptionBytes) > 0 {
		b.Description = string(crypto.Decrypt([]byte(key), descriptionBytes))
	}

	return b, nil
}
