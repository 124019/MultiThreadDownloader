import getFileList from './getFileLists.js';

const FileList = await getFileList('/我的资源');
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

console.log('Found %d folder(s) in %d item(s)', folders, count);
console.log('Found %d Bytes in %d file(s)', total_files_size, total_files);