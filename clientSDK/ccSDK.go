package clientSDK

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"time"
)

const (
	configFile  = "../clientSDK/sdkConfig.yaml"
	channelID   = "tlchannel"
	chainCodeID = "mycc"
	orgName     = "Org1"
	orgUser     = "User1"

	keyBit = 4096
)

type Copyright struct {
	Name      string    `json:"name"`
	Author    string    `json:"authorName"`
	Press     string    `json:"press"`
	Hash      string    `json:"fileHash"`
	PublicKey PublicKey `json:"public_key"`
}

type CopyRightData struct {
	*Copyright
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

func (d *CopyRightData) CalSig(pri *rsa.PrivateKey) error {
	buf := bytes.NewBuffer(make([]byte, 4096))
	buf.WriteString(d.Name)
	buf.WriteString(d.Author)
	buf.WriteString(d.Press)
	buf.WriteString(d.Hash)
	err := binary.Write(buf, binary.BigEndian, d.Timestamp)
	if err != nil {
		return err
	}

	sha1 := crypto.SHA1.New()
	sha1.Write(buf.Bytes())
	sig, err := rsa.SignPKCS1v15(buf, pri, crypto.SHA1, sha1.Sum(nil))
	if err != nil {
		return err
	}
	d.Signature = Base64Encoding.EncodeToString(sig)
	return nil
}

func Register(copyright *Copyright, pri PrivateKey) (*CopyRightData, error) {
	data := CopyRightData{
		Copyright: copyright,
		Signature: "",
		Timestamp: time.Now().Unix(),
	}
	priKey, err := DecodePrivateKey(pri)
	if err != nil {
		return nil, err
	}
	err = data.CalSig(priKey)
	if err != nil {
		return nil, err
	}
	jsonD, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	//Send to ChainCode
	req, err := client.Execute(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "register",
		Args:        [][]byte{[]byte(data.Hash), []byte(data.Name), []byte(data.Author), []byte(data.Press), jsonD},
	})
	println(req.Payload)
	return &data, err
}

func GetRightByHash(hash string) (*CopyRightData, error) {
	req, err := client.Query(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "query",
		Args:        [][]byte{[]byte(hash)},
	})
	if err != nil {
		return nil, err
	}
	data := new(CopyRightData)
	err = json.Unmarshal(req.Payload, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetRightByInfo(name string, author string, press string) (*CopyRightData, error) {
	req, err := client.Query(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "queryHash",
		Args:        [][]byte{[]byte(name), []byte(author), []byte(press)},
	})
	if err != nil {
		return nil, err
	}
	return GetRightByHash(string(req.Payload))
}

var (
	client *channel.Client
)

func init() {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		log.Fatal(err)
	}
	context := sdk.ChannelContext(channelID, fabsdk.WithOrg(orgName), fabsdk.WithUser(orgUser))
	client, err = channel.New(context)
	if err != nil {
		log.Fatal(err)
	}
	println("init ccSDK OK!")
}
