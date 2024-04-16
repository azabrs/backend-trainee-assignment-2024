package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
    "log"
    "encoding/json"
    "io"
)
func Register()string{
    c := http.Client{}
    q := url.Values{}
    q.Add("login", "apuha")
    q.Add("password", "12345678")
    u := url.URL{
            Scheme: "http",
            Host: "localhost:8080",
            Path: "register",
            RawQuery: q.Encode(),
        }
    fmt.Println(u.String())
    req, err := http.NewRequest(http.MethodPost, u.String(), nil)
    if err != nil { 
        fmt.Println(err) 
        return ""

    } 
    req.Header.Set("is_admin", "true")
    resp, err := c.Do(req)

	if err != nil { 
        fmt.Println(err) 
        return ""

    } 
    defer resp.Body.Close() 
    if resp.StatusCode == http.StatusCreated{
        var j map[string]string
        err = json.NewDecoder(resp.Body).Decode(&j)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s", j)
        return j["Token"]
    } else{
        bytes := make([]byte, 100)
        n, err := resp.Body.Read(bytes)
        bytes = bytes[:n]
        log.Println(string(bytes))
        if err != nil {
            if err == io.EOF {
                return ""
            }
                log.Panic(err)
        }

    }
    
    return ""
}

func TestGetBanner(t *testing.T){
    token := Register()
    fmt.Println(token)
    q := url.Values{}
    c := http.Client{}
    q.Add("tag_id", "2")
    q.Add("feature_id", "1")
    q.Add("use_last_revision", "false")
    u := url.URL{
        Scheme: "http",
        Host: "localhost:8080",
        Path: "user_banner",
        RawQuery: q.Encode(),
    }
    fmt.Println(u.String())
    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    if err != nil { 
        fmt.Println(err) 
        return

    } 
    req.Header.Set("token", token)
    resp, err := c.Do(req)
	if err != nil { 
        fmt.Println(err) 
        return

    } 

    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK{
        t.Fail()
    }
}