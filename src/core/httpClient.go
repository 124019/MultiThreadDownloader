// core/HttpClient.go
package core

import (
    "net"
    "net/http"
    "time"
)

var sharedTransport = &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   15 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          200,         // 全局最大空闲连接
    MaxIdleConnsPerHost:   100,         // ⭐ 关键：默认是 2，必须改大
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1  * time.Second,
    ResponseHeaderTimeout: 20 * time.Second, // 单独限制"等响应头"的时间
}

var SharedHTTPClient = &http.Client{
    Transport: sharedTransport,
    // ⚠️ 注意：这里不要设 Timeout。用 Transport 的各个细分超时来管，
    // 否则又会变成"整个请求 5 秒"
}