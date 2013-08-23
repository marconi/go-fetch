package main

import (
    "flag"
    "fmt"
    "github.com/marconi/gofetch/fetcher"
)

var f_url string
var f_workers int
var f_name string

func init() {
    flag.StringVar(&f_url, "url", "", "URL of the file to download")
    flag.StringVar(&f_name, "filename", "", "Name of downloaded file")
    flag.IntVar(&f_workers, "workers", 2, "Number of download workers")
}

func main() {
    flag.Parse()
    d := fetcher.NewDownloader(f_url, f_workers, f_name)
    _, err := d.GetHeaders()
    if err != nil {
        fmt.Println(err)
    } else {
        d.Download()
        var input string
        fmt.Scanln(&input)
    }
}
