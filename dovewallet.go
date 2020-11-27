package dovewallet

import (
    "crypto/hmac"
    "crypto/sha512"
    "encoding/hex"
	"fmt"
	"io/ioutil"
	"time"
	"net/http"
	"sort"
)

// Dovewallet struct
type Dovewallet struct {
	PublicKey, SecretKey string
}
// Request is request dovewallet api
func (dw *Dovewallet) Request(method string, option map[string]string, nonce bool ) ([]byte, error){
	const baseURL = "https://api.dovewallet.com/v1.1"
	var url = baseURL + method
	
	option["apikey"] = dw.PublicKey
	if nonce{
		option["nonce"] = fmt.Sprint(time.Now().UnixNano() / 1000000)
	}
	query := dw.makeQuery(option)
	apisign := dw.getApisign(url+"?"+query)
	body, err := dw.httpRequest(url+"?"+query, apisign)
	if err != nil{
		return nil, err
	}
	return body, nil
}
// makeQuery is convert map to query
func (dw *Dovewallet) makeQuery(option map[string]string)string{
	keys := make([]string,0)
	for k := range option{
		keys = append(keys, k)
	}
	sort.Strings(keys)
	query := ""
	first := true
	for key := range keys{
		if first{
			query +=  keys[key] + "=" + option[keys[key]] 
			first = false
			continue
		}
		query += "&" +  keys[key] + "=" + option[keys[key]]
	}
	return query
}

// httpRequest with url and apisign
func (dw *Dovewallet) httpRequest(url, apisign string) ([]byte, error){
	client := &http.Client{}
	req,err := http.NewRequest("GET", url, nil)
	if err != nil{
		return nil, err
	}

	req.Header.Set("apisign", apisign)
	resp,err := client.Do(req)
	if err != nil{
		return nil, err
	}

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}

// getApisign by url
func (dw *Dovewallet) getApisign(url string) string{
	h := hmac.New(sha512.New, []byte(dw.SecretKey))
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}