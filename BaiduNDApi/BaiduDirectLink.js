import fs from 'fs';

export default function BaiduPan(accessToken) {

  if (!accessToken) throw new Error('no access_token!!');
  var BD_Headers = fs.readFileSync(new URL('./BD_Headers.json', import.meta.url), 'utf-8');
  var BD_Headers = JSON.parse(BD_Headers);
  const BD_Cookies = fs.readFileSync(new URL('./cookies.txt', import.meta.url), 'utf-8');
  BD_Headers['cookie'] = BD_Cookies;
  // console.log('BD_Headers:', BD_Headers);

  async function getFileMetas(fidList, maxRetry = 1) {
    if (!fidList || !fidList.length) throw new Error('not valid fidList');

    const fsids = encodeURIComponent(JSON.stringify(fidList));
    const url = `https://pan.baidu.com/rest/2.0/xpan/multimedia?method=filemetas&dlink=1&fsids=${fsids}&access_token=${accessToken}`;

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 20000);

      const response = await fetch(url, {
        headers: BD_Headers,
        signal: controller.signal,
      });
      clearTimeout(timeoutId);

      if (!response.ok) {
        throw new Error(`HTTP request error: ${response.status} - ${response.statusText}`);
      }

      const data = await response.json();

      if (data.errno === 0) {
        return data.list || [];
      } else if (data.errno === 112) {
        throw new Error('the page is expired,please refresh');
      } else if (data.errno === 9019) {
        if (maxRetry > 0) {
          throw new Error('please refresh the access_token because it is invalid or expired');
        } else {
          throw new Error(`Failed to retrieve, errno=${data.errno}`);
        }
      } else {
        throw new Error(`Baidu API returned error, errno=${data.errno}, msg=${data.msg || ''}`);
      }
    } catch (error) {
      if (error.name === 'AbortError') {
        throw new Error('Request timeout');
      }
      throw error;
    }
  }

  async function getDirectLinks(fidList) {
    const list = await getFileMetas(fidList);
    // const files = list.filter(item => item.isdir !== 1);
    return list.map(item => ({
      fs_id: item.fs_id,
      filename: item.server_filename || item.filename,
      size: item.size,
      document: item.isdir,
      dlink: item.dlink ? `${item.dlink}&access_token=${accessToken}` : null,
    }));
  }

  return {
    getFileMetas,
    getDirectLinks,
  };
}