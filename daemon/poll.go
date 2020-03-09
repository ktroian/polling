package main

import (
    "log"
    "time"
    "context"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

func parseData(body []byte) map[string]*Company {
    var companies []Company
    m := make(map[string]*Company)
    
    json.Unmarshal(body, &companies)

    for i, c := range companies {
        m[c.Name] = &companies[i]
    }

    return m
}

func getData(url string) []byte {
    req, err := http.NewRequest(http.MethodGet, url, nil)
    
    if err != nil {
        log.Println(err)
    }
    client := http.DefaultClient
    resp, err := client.Do(req)
    
    if err != nil {
        log.Println(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        log.Println(err)
    }

    return body
}

func startPolling(ctx context.Context, url string, interval int) {
    go func() {
        delta := make(map[string]*Company)

        for {
            select {
            case <- ctx.Done():
                return

            default:
                data := getData(url)
                update := parseData(data)

                // sync new items
                for name, pointer := range update {
                    if delta[name] == nil {
                        delta[name] = pointer
                        pointer.save()
                    }
                }

                // sync deleted items
                for name, pointer := range delta {
                    if update[name] == nil {
                        c := pointer
                        pointer = nil
                        c.delete()
                    }
                }

                delta = update
                time.Sleep(time.Duration(interval) * time.Second)
            }
        }
    }()
}