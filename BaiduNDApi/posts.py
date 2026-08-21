import requests
import json
import base64

response = requests.post('https://api.youxiaohou.com/config/v2?ver=6.2.7&a=YouXiaoHou')
tx = base64.b64decode(response.text).decode('utf-8')
BD_verify = json.loads(tx).get('pcs')
Baidu_verify_url = BD_verify.get('3')
print(Baidu_verify_url)

with open('baidu_verify_url.txt', 'w') as f:
    f.write(Baidu_verify_url)