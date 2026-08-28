import fs from 'fs';

async function getFileMetas(AccessToken, BD_Headers, fidList, maxRetry = 1) {
  if (!fidList || !fidList.length) throw new Error('not valid fidList');

  const fsids = encodeURIComponent(JSON.stringify(fidList));

  try {
    const url = `https://pan.baidu.com/rest/2.0/xpan/multimedia?method=filemetas&dlink=1&fsids=${fsids}&access_token=${AccessToken}`;
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

    if (response.status === 400) { // Invalid access_token or expired access_token
      throw new Error('InvalidAccessTokenError');
    }

    const data = await response.json();

    if (data.errno === 0) {
      return data.list || [];
    } else if (data.errno === 112) {
      throw new Error('the page is expired,please refresh');
    } else if (data.errno === 9019) {
      if (maxRetry > 0) {
        throw new Error('InvalidAccessTokenError');
      } else {
        throw new Error(`Failed to retrieve, errno=${data.errno}`);
      }
    } else {
      throw new Error(`Baidu API returned error, errno=${data.errno}, msg=${data.msg || ''}`);
    }
  } catch (error) {
    if (error.name === 'AbortError') {
      throw new Error('Request timeout');
    } else {
    throw error;
    }
  }
}


export default async function getDirectLinks(AccessToken, BD_Headers, fidList) {
  // if (retry_count > 2) throw new Error('OverRetryError'); // Too many retries, giving up
  if (!AccessToken) throw new Error('no access_token!!');

  const list = await getFileMetas(AccessToken, BD_Headers, fidList);
  // const files = list.filter(item => item.isdir !== 1);
  return list.map(item => ({
    fs_id: item.fs_id,
    filename: item.server_filename || item.filename,
    size: item.size,
    document: item.isdir,
    dlink: item.dlink ? `${item.dlink}&access_token=${AccessToken}` : null,
  }));
}