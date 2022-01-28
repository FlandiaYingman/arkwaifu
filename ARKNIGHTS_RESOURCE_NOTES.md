# Arknights Asset Bundle Notes

## MITM to get assets URL of Arknights

## Format of `hot_update_list.json`

Get this file from `https://ak.hycdn.cn/assetbundle/official/Android/assets/{resVersion}/hot_update_list.json`.

### `fullpack`, `versionId`

Seems useless. Since all fields are empty.

### `countOfTypedRes`

Not quite sure yet.

### `abInfos`

The asset bundle infos. It includes all assets from Arknights.

Retrieve them by `https://ak.hycdn.cn/assetbundle/official/Android/assets/{resVersion}/{name_alt}`, which name_alt is
the `name` field replacing all `/` to `_`, replacing all `#` to `__` and replacing the extension to `dat`.

The retrieved file is a PK ZIP file. Extract it, and you'll get the asset bundle files located in the path `name`.

Let's say there's an asset bundle with the name `arts/building/buffs/buff_image_config.ab`. We transform it
to `arts_building_buffs_buff_image_config.dat` and retrieve it
by `https://ak.hycdn.cn/assetbundle/official/Android/assets/{resVersion}/arts_building_buffs_buff_image_config.dat`. Now
extract `arts_building_buffs_buff_image_config.dat`, there will be an asset bundle file `buff_image_config.ab` in the
directory `arts/building/buffs`.

### `packInfos`

The "large" pack infos (I guess alphabet "l" stands for "large"). They include all asset bundles which have the
corresponded `pid` field. For example, `torappu.ab` has no `pid` field, so it's not included in any "large"
packs; `arts/building/buffs/buff_image_config.ab`'s `pid` field is `lpack_init`, so it's included in the "large"
pack `lpack_init`.

The retrieved file is a PK ZIP file just like `abInfos`. The only one difference is that multiple files shall be
extracted instead of one.

## 剧情资源（已过期）

剧情文本资源集中在`gamedata\story`中；活动剧情文本在`gamedata/story/activities`下（证据：纯文本文件，可直接查找）。

然而，活动的命名似乎没有规律。例如`act15side`（将进酒）与`act15d0`（监狱风云）、`act15d5`（此地之外）虽然都以`act15`开头，但它们似乎没有共同点。

### 图片资源

从剧情的文本资源中，形如`[Background(image="{name}" ...)]`的行是背景“命令”，可在`assets\torappu\dynamicassets\avg\backgrounds`下找到同名文件。

与背景“命令”类似，形如`[Image(image="{name}" ...)]`的行是图片“命令”，可在`assets\torappu\dynamicassets\avg\image`下找到同名文件。

### 剧情名

活动名称在`gamedata\excel\activity_table.json`中。解析该文件即可。



## 剧情资源

切入点：`excel\story_review_table.json`。该文件包含至今（2022年1月23日，将进酒）为止所有的剧情信息（包括将进酒本身）。

读取其每个字段（如`1stact`、`act3d0`、`act5d0`等），可发现它们分别代表着一个章节、活动或干员密录。

它们的`name`字段为它们的名字（主线章节名、活动名等）；`infoUnlockDatas`字段下则是具体的关卡信息。

关卡信息的`storyName`字段是关卡名，`storyCode`字段为关卡代号（如`1-7`），最关键的`storyTxt`字段则为剧情文本资源。（参考“剧情资源（已过期）”一节）

其余的信息可忽略。