// baidu_direct_link.js
const axios = require('axios');

/**
 * 百度网盘直链获取模块
 * 依赖：axios (npm install axios)
 */
class BaiduPan {
  /**
   * @param {string} accessToken - 百度网盘 access_token，需提前获取
   * @param {string} [userAgent] - 自定义User-Agent，可选
   */
  constructor(accessToken, userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36') {
    this.accessToken = accessToken;
    this.userAgent = userAgent;
    this.baseUrl = 'https://pan.baidu.com/rest/2.0/xpan/multimedia';
  }

  /**
   * 获取文件元数据（含直链）
   * @param {number[]} fidList - 文件fs_id数组
   * @param {number} [maxRetry=1] - token失效后重试次数
   * @returns {Promise<Array>} 文件信息列表，每个元素包含:
   *   { server_filename, size, dlink, fs_id, isdir, ... }
   * @throws {Error} 当请求失败或token无效时抛出
   */
  async getFileMetas(fidList, maxRetry = 1) {
    if (!this.accessToken) throw new Error('access_token 不能为空');
    if (!fidList || !fidList.length) throw new Error('fidList 不能为空');

    const fsids = encodeURIComponent(JSON.stringify(fidList));
    const url = `${this.baseUrl}?method=filemetas&app_id=250528&fsids=${fsids}&access_token=${this.accessToken}`;

    try {
      const response = await axios.get(url, {
        headers: { 'User-Agent': this.userAgent },
        timeout: 10000,
      });
      const data = response.data;

      // 处理百度返回码
      if (data.errno === 0) {
        // 成功，返回列表
        return data.list || [];
      } else if (data.errno === 112) {
        throw new Error('页面已过期，请刷新或重新获取 access_token');
      } else if (data.errno === 9019) {
        // token无效，尝试刷新（原脚本通过再次获取token，但这里我们直接抛出，由上层重新设置token）
        if (maxRetry > 0) {
          // 若调用者传入新的token，可再次调用，但这里我们只抛出错误提示
          throw new Error('access_token 无效或已过期，请重新获取');
        } else {
          throw new Error(`获取失败，errno=${data.errno}`);
        }
      } else {
        throw new Error(`百度接口返回错误，errno=${data.errno}, msg=${data.msg || ''}`);
      }
    } catch (error) {
      if (error.response) {
        // 服务器返回了错误状态码
        throw new Error(`HTTP请求失败: ${error.response.status} - ${error.response.statusText}`);
      }
      throw error;
    }
  }

  /**
   * 获取直链列表（便于直接使用）
   * @param {number[]} fidList - 文件fs_id数组
   * @returns {Promise<Array>} 每个元素包含 { filename, size, dlink, fs_id }
   */
  async getDirectLinks(fidList) {
    const list = await this.getFileMetas(fidList);
    // 过滤掉文件夹（isdir=1），只返回文件
    const files = list.filter(item => item.isdir !== 1);
    return files.map(item => ({
      fs_id: item.fs_id,
      filename: item.server_filename || item.filename,
      size: item.size,
      // 原dlink需拼接access_token才能直接下载
      dlink: item.dlink ? `${item.dlink}&access_token=${this.accessToken}` : null,
    }));
  }
}

module.exports = BaiduPan;