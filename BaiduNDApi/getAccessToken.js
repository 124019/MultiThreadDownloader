import axios from 'axios';
import { readFileSync } from 'fs';
import { writeFileSync } from 'fs';

const cookies = readFileSync(new URL('./cookies.txt', import.meta.url), 'utf-8');
console.log('Cookies:', cookies);

const baidu_verify_url = readFileSync(new URL('./baidu_verify_url.txt', import.meta.url), 'utf-8');
console.log('Baidu Verify URL:', baidu_verify_url);

const config = {
  method: 'get',
  url: baidu_verify_url,
  headers: {
    'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
    'accept-encoding': 'gzip, deflate, br, zstd',
    'accept-language': 'zh-CN,zh-TW;q=0.9,zh;q=0.8',
    'cookie': cookies,
    'priority': 'u=0, i',
    'sec-ch-ua': '"Chromium";v="148", "Google Chrome";v="148", "Not/A)Brand";v="99"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'document',
    'sec-fetch-mode': 'navigate',
    'sec-fetch-site': 'none',
    'sec-fetch-user': '?1',
    'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36'
  },
  maxRedirects: 5,       // 自动跟随重定向（最多5次）
  validateStatus: (status) => status >= 200 && status < 400
};

async function getTokenFromRedirect() {
  try {
    const response = await axios(config);
    const finalUrl = response.request.res.responseUrl || response.request.path;
    if (!finalUrl) {
      console.log('finalUrl:', response.config.url);
    }
    const match = finalUrl.match(/access_token=([^&]+)/);
    if (match) {
      return match[1];
    } else {
      console.warn('Cannot find access_token in the final URL:', finalUrl);
    }
  } catch (error) {
    console.error('Request failed:', error.message);
  }
}

var AccessToken = await getTokenFromRedirect();
console.log('Access Token:', AccessToken);

writeFileSync(new URL('./AccessToken.txt', import.meta.url), AccessToken, 'utf-8');
console.log('Access Token saved to AccessToken.txt');
