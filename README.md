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
    r18: true # 易封号
    blacklist: # 黑名单
      - 1781924496
```
