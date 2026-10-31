# MTD Frontend Development Documentation
> This document will keep updating as the project progresses.

### Development Guidings
> 因为开发便捷所以暂时直接使用了flutter/material，同时在此基础上做一些改动。

##### General
- 一些按钮（或者不动的元素）保持使用`material`组件
##### Animation
- 需要移动的组件改为使用`Widget`**单独绘制**，并使用`material`组件提供的颜色，同时另外计算**加速度与变形**，此情况仅适用于点按，即**可预测落点的情况**。另外，减少或避免制作允许拖拽的组件（不可预测落点）
- 在按钮被点按以后，会出现简单的按压形变与平滑转色（原色->按压后：预设->松手：原色）,可能需要将原本的按钮动画覆盖掉
- 切换页面时，按钮动画为扩散至整个屏幕同时逐渐变为背景色，然后才关闭。同时页面的其他元素在下一页面充分加载以后再关闭（控制栏例外）
##### Color
- 所有的颜色都应该使用`material`提供的颜色，非必要不使用自定义颜色

### UI/X Design Guidelines

- 所有组件都是色块
- 所有点击都应该有反馈
- 能够移动组件或色块就不应该让其消失
- 色块不应该是硬且不可变的，应该能够随着点击和移动（因加速度的变形）而变化

### Structure

- `lib`目录下存放所有代码
- `lib/main.dart`为入口文件

### System Support (Temporary)

**Desktop**
- Windows AMD64 (dev)
- ~~Linux~~[^1]
- ~~MacOS~~

**~~Mobile~~**[^2]
- Android Arm64 *
- iOS/iPadOS Arm64 *

**~~Web~~**[^3]

[^1]:我没有Linux和MacOS设备用于测试，因此不会进行开发
[^2]:在Desktop版本开发完成之前，暂时不会支持移动设备
[^3]:不会考虑另外支持Web