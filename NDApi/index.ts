import { Elysia, t } from 'elysia'
import { join } from 'node:path'
import  process_api_info from './src/core/process_api_info.ts'
import type { ProcessedApiInfo } from './src/core/process_api_info.ts';


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


import { getDirectLink, fetchFileList } from './api/official_api/BaiduNDApi/BaiduND_main.ts'
import type * as BaiduNDType from './api/official_api/BaiduNDApi/type.ts'

const BaiduND = new Elysia({'prefix': 'BaiduND'})
  .post('/filelist', async ({ body, status }) => {
      if (!api_info.support_list.includes('BaiduND')) {
        return status(404, { message: 'BaiduND is not supported' });
      } else if (!(api_info.support_details.BaiduND.features).includes('filelist')) {
        return status(404, { message: 'BaiduND filelist feature is not supported' });
      } else {
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
      if (!api_info.support_list.includes('BaiduND')) {
        return status(404, { message: 'BaiduND is not supported' });
      } else if (!(api_info.support_details.BaiduND.features).includes('directlink')) {
        return status(404, { message: 'BaiduND directlink feature is not supported' });
      } else {
        const data = await getDirectLink(body.FileList as BaiduNDType.FileListItem[]);
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

const QuarkND = new Elysia({'prefix': 'QuarkND'})

const app = new Elysia({'prefix': 'mtd'})
  .use(BaiduND)
  .use(QuarkND)
  .listen(1242)