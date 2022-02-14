package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/Cloud-Recording-Golang/schemas"

	"github.com/spf13/viper"
)

var Regions = map[int]string{
	0:  "us-east-1",
	1:  "us-east-2",
	2:  "us-west-1",
	3:  "us-west-2",
	4:  "eu-west-1",
	5:  "eu-west-2",
	6:  "eu-west-3",
	7:  "eu-central-1",
	8:  "ap-southeast-1",
	9:  "ap-southeast-2",
	10: "ap-northeast-1",
	11: "ap-northeast-2",
	12: "sa-east-1",
	13: "ca-central-1",
	14: "ap-south-1",
	15: "cn-north-1",
	16: "cn-northwest-1",
	17: "us-gov-west-1",
}

// Recorder manages cloud recording
type Recorder struct {
	http.Client
	CallInfo schemas.CallInfo
	Token    string
	UID      int
	RID      string
	SID      string
}

type StatusStruct struct {
	Resourceid     string `json:"resourceId"`
	Sid            string `json:"sid"`
	Serverresponse struct {
		Filelistmode string `json:"fileListMode"`
		Filelist     []struct {
			Filename       string `json:"filename"`
			Tracktype      string `json:"trackType"`
			UID            string `json:"uid"`
			Mixedalluser   bool   `json:"mixedAllUser"`
			Isplayable     bool   `json:"isPlayable"`
			Slicestarttime int64  `json:"sliceStartTime"`
		} `json:"fileList"`
		Status         int   `json:"status"`
		Slicestarttime int64 `json:"sliceStartTime"`
	} `json:"serverResponse"`
}

// Acquire runs the acquire endpoint for Cloud Recording
func (rec *Recorder) Acquire() (string, error) {
	creds, err := GenerateUserCredentials(rec.CallInfo.Channel)
	if err != nil {
		return "", err
	}

	rec.UID = creds.UID
	rec.Token = creds.Rtc

	requestBodyMap := map[string]interface{}{
		"cname": rec.CallInfo.Channel,
		"uid":   strconv.Itoa(rec.UID),
		"clientRequest": map[string]interface{}{
			"resourceExpiredHour": 24,
		},
	}

	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return "", err
	}

	// requestBody := fmt.Sprintf(`
	// 	{
	// 		"cname": "%s",
	// 		"uid": "%d",
	// 		"clientRequest": {
	// 			"resourceExpiredHour": 24
	// 		}
	// 	}
	// `, rec.CallInfo.Channel, rec.UID)

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/acquire",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	resp, err := rec.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	rec.RID = result["resourceId"]
	b, _ := json.Marshal(result)

	fmt.Println("-----Acquire-----")
	fmt.Printf("Response Body: %s\n\n", b)

	return string(b), nil
}

// Start starts the recording
func (rec *Recorder) Start() (string, error) {
	now := time.Now()

	var uidList []string
	for _, streamer := range rec.CallInfo.Streamers {
		uidList = append(uidList, streamer.Uid)
	}

	requestBodyMap := map[string]interface{}{
		"cname": rec.CallInfo.Channel,
		"uid":   strconv.Itoa(rec.UID),
		"clientRequest": map[string]interface{}{
			"token": rec.Token,
			"recordingConfig": map[string]interface{}{
				"channelType":  0,
				"audioProfile": 2,
				"maxIdleTime":  30,
				"streamTypes":  2,
				"transcodingConfig": map[string]interface{}{
					"width":            500,
					"height":           1080,
					"fps":              30,
					"bitrate":          1710,
					"mixedVideoLayout": 0,
					"backgroundColor":  "#000000",
					"backgroundConfig": rec.CallInfo.Streamers,
				},
				"subscribeAudioUids": uidList,
				"subscribeVideoUids": uidList,
			},
			"recordingFileConfig": map[string]interface{}{
				"avFileType": []string{"hls", "mp4"},
			},
			"storageConfig": map[string]interface{}{
				"vendor":    viper.GetInt("RECORDING_VENDOR"),
				"region":    viper.GetInt("RECORDING_REGION"),
				"bucket":    viper.GetString("BUCKET_NAME"),
				"accessKey": viper.GetString("BUCKET_ACCESS_KEY"),
				"secretKey": viper.GetString("BUCKET_ACCESS_SECRET"),
				"fileNamePrefix": []string{
					fmt.Sprintf("%d", now.Year()),
					fmt.Sprintf("%02d", now.Month()),
					fmt.Sprintf("%02d", now.Day()),
					rec.CallInfo.Channel,
				},
			},
		},
	}
	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return "", err
	}

	// requestBody := fmt.Sprintf(`
	// 	{
	// 		"cname": "%s",
	// 		"uid": "%d",
	// 		"clientRequest": {
	// 			"token": "%s",
	// 			"recordingConfig": {
	// 				"channelType": 0,
	// 				"audioProfile": 2,
	// 				"maxIdleTime": 30,
	// 				"streamTypes": 2,
	// 				"transcodingConfig": {
	// 					"width": 500,
	// 					"height": 1080,
	// 					"fps": 30,
	// 					"bitrate": 1710,
	// 					"mixedVideoLayout": 0,
	// 					"backgroundColor": "#000000",
	// 					"backgroundConfig": "%s",
	// 				},
	// 				"subscribeAudioUids": "%s",
	// 				"subscribeVideoUids": "%s",
	// 			},
	// 			"recordingFileConfig": {
	// 				"avFileType": ["hls", "mp4"]
	// 			},
	// 			"storageConfig": {
	// 				"vendor": %d,
	// 				"region": %d,
	// 				"bucket": "%s",
	// 				"accessKey": "%s",
	// 				"secretKey": "%s",
	// 				"fileNamePrefix": ["%d", "%d", "%d", "%s"]
	// 			}
	// 		}
	// 	}
	// `, rec.CallInfo.Channel, rec.UID, rec.Token, backgroundConfig, uidListJson, uidListJson, viper.GetInt("RECORDING_VENDOR"), viper.GetInt("RECORDING_REGION"), viper.GetString("BUCKET_NAME"),
	// 	viper.GetString("BUCKET_ACCESS_KEY"), viper.GetString("BUCKET_ACCESS_SECRET"),
	// 	now.Year(), now.Month(), now.Day(), rec.CallInfo.Channel)

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rec.RID+"/mode/mix/start",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	resp, err := rec.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	rec.SID = result["sid"]
	b, _ := json.Marshal(result)

	fmt.Println("-----Start-----")
	fmt.Printf("Response Body: %s\n\n", b)

	return string(b), nil
}

// Stop stops the cloud recording
func Stop(channel string, uid int, rid string, sid string) (string, error) {
	requestBodyMap := map[string]interface{}{
		"cname": channel,
		"uid":   strconv.Itoa(uid),
		"clientRequest": map[string]interface{}{
			"async_stop": true,
		},
	}
	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rid+"/sid/"+sid+"/mode/mix/stop",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	b, _ := json.Marshal(result)

	fmt.Println("-----Stop-----")
	fmt.Printf("Response Body: %s\n\n", b)

	return string(b), nil
}

func Query(rid string, sid string) (StatusStruct, error) {
	url := "https://api.agora.io/v1/apps/" + viper.GetString("APP_ID") + "/cloud_recording/resourceid/" + rid + "/sid/" + sid + "/mode/mix/query"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return StatusStruct{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return StatusStruct{}, err
	}

	defer resp.Body.Close()

	var result StatusStruct
	json.NewDecoder(resp.Body).Decode(&result)
	// // b, _ := json.Marshal(result)
	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//     return "", err
	// }
	// result := string(bodyBytes)
	return result, nil
}

// Update updates the cloud recording
func Update(channel string, uid int, rid string, sid string, streamers []schemas.Streamer) (string, error) {
	var uidList []string
	for _, streamer := range streamers {
		uidList = append(uidList, streamer.Uid)
	}

	requestBodyMap := map[string]interface{}{
		"cname": channel,
		"uid":   strconv.Itoa(uid),
		"clientRequest": map[string]interface{}{
			"streamSubscribe": map[string]interface{}{
				"audioUidList": map[string]interface{}{
					"subscribeAudioUids": uidList,
				},
				"videoUidList": map[string]interface{}{
					"subscribeVideoUids": uidList,
				},
			},
		},
	}
	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rid+"/sid/"+sid+"/mode/mix/update",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	b, _ := json.Marshal(result)
	return string(b), nil
}

// UpdateLayout updates the layout of cloud recording
func UpdateLayout(channel string, uid int, rid string, sid string, streamers []schemas.Streamer) (string, error) {
	requestBodyMap := map[string]interface{}{
		"cname": channel,
		"uid":   uid,
		"clientRequest": map[string]interface{}{
			"mixedVideoLayout": 0,
			"backgroundColor":  "#000000",
			"backgroundConfig": streamers,
		},
	}
	requestBody, err := json.Marshal(requestBodyMap)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rid+"/sid/"+sid+"/mode/mix/updateLayout",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	b, _ := json.Marshal(result)
	return string(b), nil
}
