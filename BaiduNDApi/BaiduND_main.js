import getFileList from './getFileLists.js';
import getDirectLinks from './BaiduDirectLink.js';
import refreshAccessToken from './getAccessToken.js';
import fs from 'fs';


var BD_Headers = fs.readFileSync(new URL('./BD_Headers.json', import.meta.url), 'utf-8');
var BD_Headers = JSON.parse(BD_Headers);
const BD_Cookies = fs.readFileSync(new URL('./cookies.txt', import.meta.url), 'utf-8');
BD_Headers['cookie'] = BD_Cookies;
// console.log('BD_Headers:', BD_Headers);

const FileList = await getFileList('/');
if (!FileList) {
    console.log('Baidu Server wasn\'t return anything.');
}

let folders = 0;
let total_files_size = 0;
let count = 0;
for (const file of FileList) {
    count += 1;
    folders += file.folder;
    total_files_size += file.size;
}
var total_files = count - folders;

console.log('Found %d folder(s) in %d item(s) on your Account.', folders, count);
console.log('Found %d file(s) from your index page, totaling %d bytes.', total_files, total_files_size);

//Download all files in the index page START

let fidList = []
for (const file of FileList) {
    if (file.folder === 0) {
        fidList.push(file.fid);
    }
}
console.log("File ids", fidList)

// Download all files in the index page END

let AccessToken = fs.readFileSync(new URL('./AccessToken.txt', import.meta.url), 'utf-8');

let retry_count = 0;
try {
  let files = await getDirectLinks(AccessToken, BD_Headers, fidList)
  files.forEach(f => {
    console.log(`${f.filename} (${f.size} bytes)`);
    console.log(JSON.stringify(f, null, 2))
  });
} catch (error) {
  if (error.message === 'InvalidAccessTokenError') {
    retry_count += 1;
    AccessToken = await refreshAccessToken();
    console.log('Access Token is refreshed.');
    let files = await getDirectLinks(AccessToken, BD_Headers, fidList)
    files.forEach(f => {
      console.log(`${f.filename} (${f.size} bytes)`);
      console.log(JSON.stringify(f, null, 2))
    });
  } else {
    console.log('Error:', error.message);
  }
}