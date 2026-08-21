const BaiduPan = require('./baidu_direct_link');

// 假设已从浏览器或其他方式获取到 access_token
const accessToken = 'YOUR_ACCESS_TOKEN';
const pan = new BaiduPan(accessToken);

// 文件ID列表（从百度网盘页面获取，例如通过浏览器开发者工具）
const fidList = [123456789, 987654321];

pan.getDirectLinks(fidList)
  .then(files => {
    console.log('获取到的直链：');
    files.forEach(f => {
      console.log(`${f.filename} (${f.size} bytes): ${f.dlink}`);
    });
  })
  .catch(err => {
    console.error('获取失败：', err.message);
  });