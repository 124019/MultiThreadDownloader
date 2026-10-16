import { Elysia, t } from 'elysia'
import { getDirectLink, fetchFileList } from './BaiduNDApi/BaiduND_main.ts'

const BaiduND = new Elysia({'prefix': 'BaiduND'})
  .post('/filelist', ({ body, status }) => {
      return status(200, fetchFileList(body.ND_Path));
    },
    {
      body: t.Object({
        ND_Path: t.String()
      })
    }
  )
  .post('/directlink', ({ body, status }) => {
      return status(200, getDirectLink(body.FileList));
    },
    {
      body: t.Object({
        FileList: t.Array(t.Object({
          name: t.String(),
          path: t.String(),
          folder: t.Number(),
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