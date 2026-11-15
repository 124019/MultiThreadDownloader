import { Elysia, t } from 'elysia'
import { join } from 'node:path'
import  process_api_info from './src/lib/process_api_info.ts'
import type { ProcessedApiInfo } from './src/lib/process_api_info.ts';


async function GetAllApiInfos(): Promise<Array<ProcessedApiInfo>> {
  const apiDirPath = join(import.meta.dir, 'api');
  const glob = new Bun.Glob("*/api_info.json");
  // const result: Record<string, ApiInfo> = Object.create(null);
  const return_result: Array<ProcessedApiInfo> = [];

  for await (const relPath of glob.scan({ cwd: apiDirPath, absolute: false })) { // Get all paths of api_info.json in api directory
    const absPath = join(apiDirPath, relPath); // Get absolute path of all api_info.json in the first-level subdirectories of the 'api' folder
    try {
      const data: Record<string, unknown> = (await Bun.file(absPath).json());
      // console.log(data);
      const result: ProcessedApiInfo = await process_api_info(data);
      return_result.push(result); // Push the processed api_info.json into the return_result array
    } catch (err) {
      console.error(`Skip: Unable to parse: ${absPath}:`, err);
    };
  };
  return return_result;
}

const api_infos: Array<ProcessedApiInfo> = await GetAllApiInfos();
console.log(JSON.stringify(api_infos, null, 2));

// GET ALL API SUPPORT PLATFORMS LIST
// const support_list = function (): Array<string> {
//   let support_list: Array<string> = [];
//   for (const api_info of api_infos) {
//     const each_support_list = api_info.support_list;
//     if (!isStringArray(each_support_list)) {
//       throw new Error(`each_support_list is not a string array: ${each_support_list}`);
//     };
//     for (const each_support_platform of each_support_list) {
//       support_list.push(each_support_platform);
//     };
//   };
//   return support_list;
// }()
// console.log(support_list)

const api_info = api_infos[0];// Temporary code, just get the first api. It will be More complex in the future

import type * as BaiduNDType from './api/official_api/BaiduNDApi/BaiduND_main.ts'
import type * as BaiduFileType from './api/official_api/BaiduNDApi/type.ts'
// import { fetchFileList, getDirectLink } from './api/official_api/BaiduNDApi/BaiduND_main.ts'

type BaiduNDShape = {
  getDirectLink: typeof BaiduNDType.getDirectLink | null
  fetchFileList: typeof BaiduNDType.fetchFileList | null
}

const { getDirectLink, fetchFileList } = await (async (): Promise<BaiduNDShape> => {
  const EMPTY: BaiduNDShape = { getDirectLink: null, fetchFileList: null }
  if (!api_info.support_list.includes('BaiduNetdisk')) {
    return EMPTY;
  }
  try {
    if (api_info.support_list.includes('BaiduNetdisk')) {
      console.log('BaiduNetdisk is supported');
      const module = await import('./api/official_api/BaiduNDApi/BaiduND_main.ts')
      return { getDirectLink: module.getDirectLink, fetchFileList: module.fetchFileList }
    } else {
      console.log('BaiduNetdisk is not supported');
    }
  } catch (err) {
    console.warn("Api Core is corrupted: The api_info file is incorrect or Api Core is incomplete");
    console.log("Api Core is corrupted: The api_info file is incorrect or Api Core is incomplete");
    return EMPTY;
  }
  console.log('BaiduNetdisk is not supported: EMPTY');
  return EMPTY;
})();

const BaiduNetdisk = new Elysia({'prefix': 'BaiduNetdisk'})
  .post('/filelist', async ({ body, status }) => {
      if (!api_info.support_list.includes('BaiduNetdisk')) {
        return status(404, { message: `BaiduNetdisk is not supported, support list: ${api_info.support_list}` });
      } else if (!(api_info.support_details.BaiduNetdisk.features).includes('filelist')) {
        return status(404, { message: `BaiduNetdisk filelist feature is not supported, supported features: ${api_info.support_details.BaiduNetdisk.features}` });
      } else {
        if (fetchFileList === null) {
          return status(500, { message: `BaiduNetdisk filelist feature is not implemented: The api_info file is incorrect or Api Core is incomplete: ${api_info.support_details.BaiduNetdisk.features}` });
        };
        const data = await fetchFileList(body.ND_Path as string);
        return status(200, data);
      }
    },
    {
      body: t.Object({
        ND_Path: t.String()
      }, { additionalProperties: false })
    }
  )
  .post('/directlink', async ({ body, status }) => {
      if (!api_info.support_list.includes('BaiduNetdisk')) {
        return status(404, { message: 'BaiduNetdisk is not supported' });
      } else if (!(api_info.support_details.BaiduNetdisk.features).includes('directlink')) {
        return status(404, { message: 'BaiduNetdisk directlink feature is not supported' });
      } else {
        if (getDirectLink === null) {
          return status(500, { message: 'BaiduNetdisk directlink feature is not implemented: The api_info file is incorrect or Api Core is incomplete' });
        };
        const data = await getDirectLink(body.FileList as BaiduFileType.FileListItem[]);
        return status(200, data);
      }
    },
    {
      body: t.Object({
        FileList: t.Array(t.Object({
          name: t.String(),
          path: t.String(),
          isdir: t.Number(),
          fid: t.Number(),
          size: t.Number()
        }, { additionalProperties: false })
      )
    }, { additionalProperties: false })
  })

const QuarkNetdisk = new Elysia({'prefix': 'QuarkNetdisk'})

const app = new Elysia({ prefix: '/mtd' })
  .onError(({ code, error, path, set }) => {
    console.error('[Elysia Error]', code, path, error);
    set.headers['content-type'] = 'application/json';
    return { code, message: String(error) };
  })
  .use(BaiduNetdisk)
  .use(QuarkNetdisk)
  .listen(1242)
  .use(BaiduNetdisk)
  .use(QuarkNetdisk)
  .listen(1242)