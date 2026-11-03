还在开发中，这是一个用于开发的仓库

这将会是一个本体以`GNU GPLv3`协议开源的多线程的网盘资源下载器，它支持以插件的形式插入适配的网盘Api插件以防止我内置的插件失效，同时也会内置一个由我开发的官方插件，它将会以`GNU AGPLv3`协议开源，Star我的仓库以在第一时间获取最新信息。

我正在学习前端开发，可能会需要很长时间用于设计并实现ui，在发布正式版之前只支持自行配置环境进行源码调试。

本软件可能使用了某些其他Api服务器的一些服务，可能会在一些地方有违规之处，如果违反了你的服务器使用协定或者你作为服务器拥有者不希望我的项目使用你的服务，请在第一时间联系我，我会在正式版Release发布之前，删除使用你的服务的部分或者对你的协议进行合规性兼容。同时引用的项目列表与项目作者列表会在发布Release之前加入。

This project will no longer use web development for the frontend and will switch to Flutter instead.

The backend will be written in Golang and TypeScript. Go for downloader and TS for allow users add or exchange other Api plugins.

### Development Progress:

##### Downloader
- downloader `Go` `dev`
- ~~Main Backend Interface Entry Points~~ `planned`

##### Netdisk Apis 
- Netdisk Api Interface Entry Points `ElysiaJS / Bun TS` `beta`
- Baidu Netdisk Api Plugin Core `Bun TS` `alpha`
- ~~Quark Netdisk Api Plugin Core~~ `planned`

##### Frontend
- ~~Windows AMD64 Desktop Frontend~~ `Flutter / Dart` `dev`

### Development Document List
- /frontend/README.md -> For develop the frontend.
- /doc
