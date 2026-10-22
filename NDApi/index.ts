import { Elysia, t } from 'elysia'
import { getDirectLink, fetchFileList } from './BaiduNDApi/BaiduND_main.ts'
import type { FileListItem } from './BaiduNDApi/type.ts';

const BaiduND = new Elysia({'prefix': 'BaiduND'})
  .post('/filelist', async ({ body, status }) => {
    const data = await fetchFileList(body.ND_Path);
    return status(200, data);
  },
    {
      body: t.Object({
        ND_Path: t.String()
      })
    }
  )
  .post('/directlink', async ({ body, status }) => {
      const data = await getDirectLink(body.FileList);
      return status(200, data);
    },
    {
      body: t.Object({
        FileList: t.Array(t.Object({
          name: t.String(),
          path: t.String(),
          isdir: t.Number(),
          fid: t.Number(),
          size: t.Number()
        }))
      })
    }
  )

const QuarkND = new Elysia({'prefix': 'QuarkND'})

const app = new Elysia({'prefix': 'mtd'})
  .use(BaiduND)
  .use(QuarkND)
  .listen(1242)