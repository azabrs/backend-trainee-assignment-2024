package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)


func TestGetBanner(t *testing.T){
    q := url.Values{}
    c := http.Client{}
    q.Add("tag_id", "2")
    q.Add("feature_id", "1")
    q.Add("use_last_revision", "true")
    u := url.URL{
        Scheme: "http",
        Host: "localhost:8080",
        Path: "user_banner",
        RawQuery: q.Encode(),
    }
    fmt.Println(u.String())
    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    if err != nil { 
        t.Fatalf("Request Error: %v", err) 
        return

    } 
    req.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMyMDE4ODksImlzQWRtaW4iOnRydWV9.f7moW378Gp5MUBE0LGbDCEdSJ_eYbdboiDJZjYS6rj8")
    resp, err := c.Do(req)
	if err != nil { 
        t.Fatalf("Response Error: %v", err) 
        return

    }

    defer resp.Body.Close()
    for{
             
        bs := make([]byte, 1014)
        n, err := resp.Body.Read(bs)
        fmt.Println(string(bs[:n]))
         
        if n == 0 || err != nil{
            break
        }
    }
}