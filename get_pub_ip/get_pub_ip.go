package main

import (
    "container/list"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "sort"
    "strings"
    "sync"
    "time"
)

type ParallelHttpCtx struct {
    ResultCh        chan *WorkItem
    AllResultDoneCh chan struct{}
    Results         *list.List
    Wg              sync.WaitGroup
    WaitCtx         context.Context
}
type IpGetter func(r io.Reader) (string, error)

type WorkItem struct {
    Uri      string
    IpGetter IpGetter
    Result   string
}

// do http.get with context.Context
// context can be cancel()
// not care about err, only push result to chan
func fetchIpRoutine(wk *WorkItem, ctx *ParallelHttpCtx) {
    defer ctx.Wg.Done()
    defer func() {
        // always push result
        ctx.ResultCh <- wk
    }()
    //
    clt := &http.Client{}
    req, err := http.NewRequest(http.MethodGet, wk.Uri, nil)
    if err != nil {
        return
    }
    req = req.WithContext(ctx.WaitCtx)

    resp, err := clt.Do(req)
    if err != nil {
        return
    }
    if resp.StatusCode != 200 {
        return
    }
    v, err := wk.IpGetter(resp.Body)
    if err != nil {
        return
    }
    wk.Result = strings.TrimSpace(string(v))
}

// wait all sub routines result
// when all result done before cancel() then notify a chan
// can be cancel()
func waitAllResultRoutine(ctx *ParallelHttpCtx, cnt int) {
    defer ctx.Wg.Done()
    //
    for i := 0; i < cnt; i += 1 {
        select {
        case <-ctx.WaitCtx.Done():
            return
        case v := <-ctx.ResultCh:
            ctx.Results.PushBack(v)
        }
    }
    close(ctx.AllResultDoneCh)
}

// find the most often result
type Pair struct {
    Key   string
    Value int
}
type PairArray []Pair

func (p PairArray) Len() int               { return len(p) }
func (p PairArray) Less(i int, j int) bool { return p[i].Value < p[j].Value }
func (p PairArray) Swap(i int, j int)      { p[i], p[j] = p[j], p[i] }
func getTop(ctx *ParallelHttpCtx) string {
    const defaultr = "0.0.0.1"
    if ctx.Results.Len() <= 0 {
        return defaultr
    }
    //
    rm := make(map[string]int, ctx.Results.Len())
    for e := ctx.Results.Front(); e != nil; e = e.Next() {
        v := e.Value.(*WorkItem)
        if v.Result != "" {
            rm[v.Result] ++
        }
    }
    if len(rm) <= 0 {
        return defaultr
    }
    ra := make(PairArray, len(rm))
    i := 0
    for k, v := range rm {
        ra[i] = Pair{k, v}
        i++
    }
    sort.Stable(sort.Reverse(ra))
    return ra[0].Key
}

// parse result for https://ip.nf/me.json
func GetIpInJsonIpIp(r io.Reader) (string, error) {
    b, err := ioutil.ReadAll(r)
    if err != nil {
        return "", err
    }
    m := make(map[string]interface{})
    _ = json.Unmarshal(b, &m)
    ipObj, ok := m["ip"].(map[string]interface{})
    if ok {
        ip, ok := ipObj["ip"].(string)
        if ok {
            return ip, nil
        }
    }
    return "", fmt.Errorf("not found ip in %v", string(b))
}

// parse result for http://ip-api.com/json
func GetIpInJsonQuery(r io.Reader) (string, error) {
    b, err := ioutil.ReadAll(r)
    if err != nil {
        return "", err
    }
    m := make(map[string]interface{})
    _ = json.Unmarshal(b, &m)
    ip, ok := m["query"].(string)
    if ok {
        return ip, nil
    }
    return "", fmt.Errorf("not found ip in %v", string(b))
}

// parse result for https://wtfismyip.com/json
func GetIpInJsonYourFuck(r io.Reader) (string, error) {
    b, err := ioutil.ReadAll(r)
    if err != nil {
        return "", err
    }
    m := make(map[string]interface{})
    _ = json.Unmarshal(b, &m)
    ip, ok := m["YourFuckingIPAddress"].(string)
    if ok {
        return ip, nil
    }
    return "", fmt.Errorf("not found ip in %v", string(b))
}
func GetIpInPlainText(r io.Reader) (string, error) {
    b, err := ioutil.ReadAll(r)
    if err != nil {
        return "", err
    }
    v := strings.TrimSpace(string(b))
    return v, nil
}

// 1 并行发起http请求
// 2 给定超时，有几个返回结果用几个返回结果
//   且如果在超时时间内全部请求得到返回，这将会是更好的场面，
//   我们就不必要一直死等超时，直接取用结果
// 3 没有必要设置捕获 signal 信号，CTRL +C 可以在任意时刻退出，go 保证
//   我们也没有要优雅退出的需求
func main() {
    // no need to use https://api.ipify.org/?format=json
    pubSrvs := &[...]WorkItem{
        {Uri: "https://ip.nf/me.json", IpGetter: GetIpInJsonIpIp},
        {Uri: "http://ip-api.com/json", IpGetter: GetIpInJsonQuery},
        {Uri: "https://wtfismyip.com/json", IpGetter: GetIpInJsonYourFuck},
        {Uri: "https://api.ipify.org", IpGetter: GetIpInPlainText},
        {Uri: "https://ip.seeip.org", IpGetter: GetIpInPlainText},
        {Uri: "https://ifconfig.me/ip", IpGetter: GetIpInPlainText},
        {Uri: "https://ifconfig.co/ip", IpGetter: GetIpInPlainText},
    }
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Printf("pid= %v", os.Getpid())
    ctx := new(ParallelHttpCtx)
    ctx.Results = list.New()
    ctx.ResultCh = make(chan *WorkItem, len(pubSrvs))
    ctx.AllResultDoneCh = make(chan struct{})
    var cancel context.CancelFunc
    ctx.WaitCtx, cancel = context.WithCancel(context.Background())

    for i := 0; i < len(pubSrvs); i += 1 {
        ctx.Wg.Add(1)
        go fetchIpRoutine(&pubSrvs[i], ctx)
    }

    ctx.Wg.Add(1)
    go waitAllResultRoutine(ctx, len(pubSrvs))
    // wait timeout or all done
    select {
    case <-time.After(time.Second * 3):
        log.Printf("main timeup, cancel it beforehand")
        cancel()
    case <-ctx.AllResultDoneCh:
    }
    //log.Printf("main wait sub routine")
    ctx.Wg.Wait()
    //
    log.Printf("fetch result cnt= %v from %v", ctx.Results.Len(), len(pubSrvs))
    //fmt.Printf("The pub ip= %v\n", getTop(ctx))
    for i := 0; i < len(pubSrvs); i += 1 {
        log.Printf("%v -> %v", pubSrvs[i].Uri, pubSrvs[i].Result)
    }
    fmt.Printf("%v", getTop(ctx))
    log.Printf("main exit")
}
