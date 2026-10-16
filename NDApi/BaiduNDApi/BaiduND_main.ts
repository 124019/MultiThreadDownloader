//to do list: 
// Set error type

import getFileList from './getFileLists.ts';
import getFileInfo from './getFileInfo.ts';
import refreshAccessToken from './getAccessToken.ts';
import fs from 'fs';
import type { FileListItem } from './getFileLists.ts'

function readLocalHeadersAndCookies() {
  const BD_Headers: Record<string, string> = JSON.parse(
    fs.readFileSync(new URL('./BD_Headers.json', import.meta.url), 'utf-8')
  ) as Record<string, string>;
  const BD_Cookies = fs.readFileSync(new URL('./cookies.txt', import.meta.url), 'utf-8');
  BD_Headers['cookie'] = BD_Cookies;
  // Read the Cookies and Headers from the local files
  return BD_Headers;
}

export async function fetchFileList(ND_Path: string) {
  if (!ND_Path) {
    console.log('No valid path provided. Use the default path: /');
    ND_Path = '/';
  }

  const FileList = await getFileList(ND_Path);
  if (!FileList) {
    console.log('Baidu Server wasn\'t return anything.');
  }

  let fileIdList: FileListItem[] = [];
  for (const file of (FileList as FileListItem[])) {
    if (file.folder === 0) {
      fileIdList.push(file);
    }
  }
  console.log("File List: ", fileIdList);

  let FolderList: FileListItem[] = [];
  for (const file of (FileList as FileListItem[])) {
    if (file.folder === 1) {
      FolderList.push(file);
    }
  }
  console.log("Folder list: ", FolderList);

  // Get info for all files in the index page END

  return FileList as FileListItem[], FolderList as FileListItem[];
}

export async function getDirectLink(FileList: FileListItem[]) {
  let fileIdList: number[] = [];
  for (const file of FileList) {
      if (file.folder === 0) {
          fileIdList.push(file.fid);
      }
  };
  console.log("File ids", fileIdList);

  let BD_Headers: Record<string, string> = readLocalHeadersAndCookies();

  let AccessToken = fs.readFileSync(new URL('./AccessToken.txt', import.meta.url), 'utf-8');

  let retry_count = 0;
  try {
    let files = await getFileInfo(AccessToken, BD_Headers, fileIdList);
    return files;
  } catch (error) {
    if ((error as Error).message === 'InvalidAccessTokenError') {
      retry_count += 1;
      AccessToken = await refreshAccessToken();
      console.log('Access Token is refreshed.');
      let files = await getFileInfo(AccessToken, BD_Headers, fileIdList);
      return files;
    } else {
      console.log('Error:', (error as Error).message);
    }
  }
}