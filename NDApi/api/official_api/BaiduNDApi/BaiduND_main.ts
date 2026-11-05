import getFileList from './getFileLists.ts';
import getFileInfo from './getFileInfo.ts';
import refreshAccessToken from './getAccessToken.ts';
import type { FileListItem } from './type.ts';

async function readLocalHeadersAndCookies() {
  const BD_Headers: Record<string, string> = JSON.parse(
    await Bun.file(new URL('./BD_Headers.json', import.meta.url)).json()
  ) as Record<string, string>;
  const BD_Cookies = await Bun.file(new URL('./cookies.txt', import.meta.url)).text();
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

  let fileList: FileListItem[] = [];
  for (const file of (FileList as FileListItem[])) {
    if (file.isdir === 0) {
      fileList.push(file);
    }
  }
  console.log("File List: ", fileList);

  let FolderList: FileListItem[] = [];
  for (const file of (FileList as FileListItem[])) {
    if (file.isdir === 1) {
      FolderList.push(file);
    }
  }
  console.log("Folder list: ", FolderList);

  // Get info for all files in the index page END

  let jsonList = {
    FileList: fileList,
    FolderList: FolderList
  };
  return jsonList;
}

export async function getDirectLink(FileList: FileListItem[]) {
  const fileIdList = FileList.filter(f => f.isdir === 0).map(f => f.fid);
  console.log('File ids', fileIdList);
  if (fileIdList.length === 0) return [];

  const BD_Headers = await readLocalHeadersAndCookies();
  let AccessToken = await Bun.file(new URL('./AccessToken.txt', import.meta.url)).text();

  try {
    return await getFileInfo(AccessToken, BD_Headers, fileIdList);
  } catch (error) {
    if ((error as Error).message !== 'InvalidAccessTokenError') throw error; // 别吞
    AccessToken = await refreshAccessToken();
    console.log('Access Token is refreshed.');
    return await getFileInfo(AccessToken, BD_Headers, fileIdList);          // 失败会正常冒泡
  }
}