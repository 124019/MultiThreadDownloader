import requests
import json

with open("cookies.json", "r") as f:
    cookies = json.load(f)

headers = {
    "Accept": "application/json, text/plain, */*",
    "Accept-Encoding": "gzip, deflate, br, zstd",
    "Accept-Language": "zh-CN,zh-HK;q=0.9,zh;q=0.8,en;q=0.7,en-GB;q=0.6,en-US;q=0.5",
    "DNT": "1",
    "Referer": "https://pan.baidu.com/disk/main",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-origin",
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36 Edg/151.0.0.0",
    "X-Requested-With": "XMLHttpRequest",
    "sec-ch-ua": '"Not=A?Brand";v="99", "Microsoft Edge";v="151", "Chromium";v="151"',
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": '"Windows"',
    "sec-gpc": "1"
}

response = requests.get(url="https://pan.baidu.com/disk/main#/index?category=all&path=%2F%E6%88%91%E7%9A%84%E8%B5%84%E6%BA%90%2F%E7%88%B1%E8%A8%80%E5%8F%B65%E5%8A%A8%E6%80%81%E7%B4%A0%E6%9D%90%E5%8C%85", cookies=cookies, headers=headers)
print(response.text)
