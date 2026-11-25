package http_config

import (
    "net"
    "net/http"
    "time"
)

var sharedTransport = &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   5 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          200,
    MaxIdleConnsPerHost:   50,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   5 * time.Second,
    ExpectContinueTimeout: 1  * time.Second,
    ResponseHeaderTimeout: 20 * time.Second,
}

var SharedHTTPClient = &http.Client{
    Transport: sharedTransport,
}