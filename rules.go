package rules

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DynamicRules struct {
	StaticParam      string
	ChecksumConstant int
	ChecksumIndexes  []int
	AppToken         string
	RemoveHeaders    []string
	Revision         string
	IsCurrent        bool
	Format           string
}

var dynamicRules = DynamicRules{
	StaticParam:      "",
	ChecksumConstant: 0,
	ChecksumIndexes:  []int{},
	AppToken:         "33d57ade8c02dbc5a333db99ff9ae26a",
	RemoveHeaders:    []string{"user_id"},
	Revision:         "202406181902-06202f45c3",
	IsCurrent:        true,
	Format:           "",
}

func GetXBC(userAgent string) string {
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	}

	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	rand.Seed(time.Now().UnixNano())

	parts := []string{
		currentTime,
		fmt.Sprintf("%d", int64(1e12*rand.Float64())),
		fmt.Sprintf("%d", int64(1e12*rand.Float64())),
		userAgent,
	}

	var encodedParts []string
	for _, part := range parts {
		encodedParts = append(encodedParts, base64.StdEncoding.EncodeToString([]byte(part)))
	}

	msg := strings.Join(encodedParts, ".")
	token := sha1.New()
	token.Write([]byte(msg))
	return fmt.Sprintf("%x", token.Sum(nil))
}

func GetSignAndTime(link string, authID string) (string, string) {
	finalTime := strconv.FormatInt(time.Now().Unix(), 10)

	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", ""
	}

	path := parsedURL.Path
	query := parsedURL.RawQuery

	if query != "" {
		path = path + "?" + query
	}

	encodedArray := []string{
		dynamicRules.StaticParam,
		finalTime,
		path,
		authID,
	}

	message := strings.Join(encodedArray, "\n")
	sha1Sign := sha1.New()
	sha1Sign.Write([]byte(message))
	sha1SignHex := fmt.Sprintf("%x", sha1Sign.Sum(nil))
	sha1Bytes := []byte(sha1SignHex)

	checksum := dynamicRules.ChecksumConstant
	for _, index := range dynamicRules.ChecksumIndexes {
		checksum += int(sha1Bytes[index])
	}

	sign := fmt.Sprintf(dynamicRules.Format, sha1SignHex, checksum)
	return sign, finalTime
}

func Load() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/patchsets/onlyfans-dynamic-rules/refs/heads/main/rules.json", nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0")

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	res_body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	res_body_str := string(res_body)

	res_body_str = strings.ReplaceAll(res_body_str, "\\\"", "\"")

	res_body_str = strings.TrimPrefix(res_body_str, "\"")
	res_body_str = strings.TrimSuffix(res_body_str, "\"")

	json_data := map[string]interface{}{}

	err = json.Unmarshal([]byte(res_body_str), &json_data)

	if err != nil {
		return err
	}

	dynamicRules.StaticParam = json_data["staticParam"].(string)
	dynamicRules.ChecksumConstant = int(json_data["checksumConstant"].(float64))

	dynamicRules.ChecksumIndexes = []int{}

	checksumIndexes := convertSlice[float64](json_data["checksumIndexes"].([]interface{}))

	for _, item := range checksumIndexes {
		dynamicRules.ChecksumIndexes = append(dynamicRules.ChecksumIndexes, int(item))
	}

	dynamicRules.Format = json_data["start"].(string) + ":%s:%x:" + json_data["end"].(string)
	dynamicRules.Format = strings.ReplaceAll(dynamicRules.Format, "{}", "%s")
	dynamicRules.Format = strings.ReplaceAll(dynamicRules.Format, "{:x}", "%x")
	
	return nil
}

func convertSlice[E any](in []any) (out []E) {
	out = make([]E, 0, len(in))
	for _, v := range in {
		out = append(out, v.(E))
	}
	return
}
