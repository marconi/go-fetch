package fetcher

import (
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
)

type Downloader struct {
    url      string
    headers  map[string]string
    workers  int
    filename string
}

func (d *Downloader) GetHeaders() (map[string]string, error) {
    resp, err := http.Head(d.url)
    if err != nil {
        return d.headers, err
    }

    if resp.StatusCode != 200 {
        return d.headers, errors.New(resp.Status)
    }

    for key, val := range resp.Header {
        d.headers[key] = val[0]
    }

    return d.headers, err
}

func (d *Downloader) DownloadChunk(url string, out string, start int, stop int, c chan<- string) {
    client := new(http.Client)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, stop))
    resp, _ := client.Do(req)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
        return
    }

    file, err := os.OpenFile(out, os.O_WRONLY, 0600)
    if err != nil {
        if file, err = os.Create(out); err != nil {
            log.Fatalln(err)
            return
        }
    }
    defer file.Close()

    if _, err := file.WriteAt(body, int64(start)); err != nil {
        log.Fatalln(err)
        return
    }

    c <- fmt.Sprintf("Range %d-%d: %d", start, stop, resp.ContentLength)
}

func (d *Downloader) Download() {
    length, _ := strconv.Atoi(d.headers["Content-Length"])
    bytes_chunk := int(math.Ceil(float64(length) / float64(d.workers)))

    fmt.Println("bytes chunk: ", bytes_chunk)
    fmt.Println("file length: ", length)

    if length == 0 {
        return
    }

    c := make(chan string)

    for i := 0; i < d.workers; i++ {
        start := i * bytes_chunk
        stop := start + (bytes_chunk - 1)
        go d.DownloadChunk(d.url, d.filename, start, stop, c)
    }

    for i := 0; i < d.workers; i++ {
        fmt.Println(<-c)
    }

    fmt.Println("\nDownload complete! press <enter> to quit.")
}

func NewDownloader(url string, workers int, filename string) *Downloader {
    d := Downloader{
        url:      url,
        workers:  workers,
        filename: filename,
        headers:  make(map[string]string),
    }
    return &d
}
