import BaiduPan from './BaiduDirectLink.js';
import fs from 'fs';


const accessToken = fs.readFileSync(new URL('./AccessToken.txt', import.meta.url), 'utf-8');
const pan = new BaiduPan(accessToken);


const fidList = [753104397715006];

pan.getDirectLinks(fidList)
  .then(files => {
    console.log('url:');
    files.forEach(f => {
      console.log(`${f.filename} (${f.size} bytes): ${f.dlink}`);
      console.log(JSON.stringify(f, null, 2))
    });
  })
  .catch(err => {
    console.error('Failed to get direct links:', err.message);
  });