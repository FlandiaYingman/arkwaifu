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