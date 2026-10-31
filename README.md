还在开发中，这是一个用于开发的仓库，代码暂时不支持除了调试以外的任何形式的直接使用

这将会是一个本体以`GNU GPLv3`协议开源的全自动的高速网盘下载器，它将会支持以插件的形式插入适配的网盘Api插件以防止我内置的插件失效，同时也会内置一个由我开发的官方插件，它将会以`GNU AGPLv3`协议开源，Star我的仓库以在第一时间获取最新信息。

我正在学习前端开发，可能会需要很长时间用于设计并实现ui，在发布正式版之前只支持自行配置环境进行源码调试。

本软件可能使用了某些其他Api服务器的一些服务，可能会在一些地方有违规之处，如果违反了你的服务器使用协定或者你作为服务器拥有者不希望我的项目使用你的服务，请在第一时间联系我（包括但不限于提issue），我会在正式版Release发布之前（但愿有那一天）删除使用你的服务的部分或者对你的协议进行合规性兼容。

This project will no longer use web development for the frontend and will switch to Flutter instead.

The backend will be written in Golang and TypeScript. Go for downloader and TS for allow users add or exchange other Api plugins.

### Development Progress:

##### Downloader
- downloader `Go` `alpha`
- ~~Main Backend Interface Entry Points~~ `planned`

##### Netdisk Apis 
- Netdisk Api Interface Entry Points `ElysiaJS / Bun TS` `beta`
- Baidu Netdisk Api Plugin `Bun TS` `beta`
- ~~Quark Netdisk Api Plugin~~ `planned`

##### Frontend
- ~~Windows AMD64 Desktop Frontend~~ `Flutter / Dart` `dev`
