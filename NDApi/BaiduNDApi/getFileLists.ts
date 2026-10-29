import axios from 'axios';
import * as fs from 'node:fs';
import type { FileListItem } from './type.ts';
import { readFileSync } from 'fs';

export default async function getFileList(dir_path = '/') { // For Root Menu,use '/{folder_name}' to join folders
  
  const AccessToken = fs.readFileSync(new URL('./AccessToken.txt', import.meta.url), 'utf-8');

  const url = `https://pan.baidu.com/rest/2.0/xpan/file?method=list&access_token=${AccessToken}&dir=${encodeURIComponent(dir_path)}`;
  console.log('request URL:', url);
  
  const cookies = readFileSync(new URL('./cookies.txt', import.meta.url), 'utf-8');
  console.log('Cookies:', cookies);
  try {
    const response = await axios.get(url, {
      headers: {
        "Cookie": cookies,
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36 Edg/156.0.0.0",
      },
      timeout: 10000,
    });

    const data = response.data;
    if (data.errno !== 0) {
      console.error('baidu return error:', data);
      return;
    }
    const DataList = data.list || [];

    const fileList = DataList.map((item: any) => ({
        name: item.server_filename,
        path: item.path,
        isdir: item.isdir,
        fid: item.fs_id,
        size: item.size
    }));
    console.log('Successfully get fileList , get %d files&folders', fileList.length);
    // console.log(JSON.stringify(fileList, null, 2));
    return fileList as FileListItem[];
  } catch (error: any) {
    console.error('request failed:', (error as Error).message);
  }
}


const result = await getFileList('/');
console.log(result);