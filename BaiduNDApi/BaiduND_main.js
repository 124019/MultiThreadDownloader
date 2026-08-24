import getFileList from './getFileLists.js';
import BaiduPan from './BaiduDirectLink.js'
import fs from 'fs';

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

const accessToken = fs.readFileSync(new URL('./AccessToken.txt', import.meta.url), 'utf-8');
const pan = new BaiduPan(accessToken);

pan.getDirectLinks(fidList)
  .then(files => {
    files.forEach(f => {
      console.log(`${f.filename} (${f.size} bytes)`);
      console.log(JSON.stringify(f, null, 2))
    });
  })
  .catch(err => {
    console.error('Failed to get direct links:', err.message);
  });