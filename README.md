# MiraiGo-module-setu

ID: `com.aimerneige.setu`

Module for [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 鸣谢

本项目调用了 [Lolicon API](https://api.lolicon.app/#/setu) 获取色图。

## 功能

- 在群内收到“来点色图”时在群内发送“不可以色色！”并通过私聊发送色图给发言者。
- 在私聊收到“来点色图”时发送色图。
- 在群内收到“来点r18色图”时在群内发送“太色了！不可以！”并通过私聊发送 r18 色图给发言者。（r18 模式关闭后不可用）
- 在私聊收到“来点r18色图”时发送 r18 色图。（r18 模式关闭后不可用）
- 在群内收到“来点[关键词]色图”时在群内发送“不可以色色！”并通过私聊发送指定 tag 的色图给发言者。
- 在私聊收到“来点[关键词]色图”时发送指定 tag 的色图。

## 使用方法

在适当位置引用本包

```go
package example

imports (
    // ...

    _ "github.com/yukichan-bot-module/MiraiGo-module-setu"

    // ...
)

// ...
```

在全局配置文件中写入配置

```yaml
aimerneige:
  setu:
    private: true # 是否允许私聊
    r18: true # 易封号
    blacklist: # 黑名单
      - 1781924496
    allowed: # 开启功能的群
      - 857066811
      - 328521977
      - 306979312
```

## LICENSE

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
<img src="https://www.gnu.org/graphics/agplv3-155x51.png">
</a>

本项目使用 `AGPLv3` 协议开源，您可以在 [GitHub](https://github.com/yukichan-bot-module/MiraiGo-module-setu) 获取本项目源代码。为了整个社区的良性发展，我们强烈建议您做到以下几点：

- **间接接触（包括但不限于使用 `Http API` 或 跨进程技术）到本项目的软件使用 `AGPLv3` 开源**
- **不鼓励，不支持一切商业使用**
