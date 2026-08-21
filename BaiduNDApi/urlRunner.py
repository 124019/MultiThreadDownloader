import requests
import json


with open('./BaiduNDApi/RunUrlHeader.json', 'r') as f:
    headers = json.load(f)

url = 'https://d.pcs.baidu.com/file/d16cfed6cr01ea3aa76c1eae5142890d?fid=1101615540288-250528-352041860641814&rt=pr&sign=FDtAERK-DCb740ccc5511e5e8fedcff06b081203-nelU%2BeXMQ6z%2BCOVFSUmHBAagGxQ%3D&expires=8h&chkbd=0&chkv=0&dp-logid=2573646359366602230&dp-callid=0&dstime=1787333424&r=311453749&vuk=1101615540288&origin_appid=15195230&file_type=0&access_token=123.e2d376f48a69a36231154e5f982d3051.YgbeIGGgmNW7LVl9KBbgP4quf2jXF-xNm5hmDGx.YYkxuA' # example url

response = requests.get(url, headers=headers)
print(response.status_code)
print(response.text)