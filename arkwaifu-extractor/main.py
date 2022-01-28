import argparse
import sys
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor
from pathlib import Path
from typing import List

import PIL.Image
import UnityPy


def list_assets(src: Path, filters: List[str]):
    if src.is_file():
        env = UnityPy.load(str(src))
        for path, obj in env.container.items():
            if any(path.startswith(f) for f in filters):
                print(f"{path}")
    else:
        for it in src.glob('**/*'):
            if it.is_file():
                list_assets(it, filters)


def unpack_image_assets(src: Path, dst: Path, filters: List[str]):
    PIL.Image.preinit()
    PIL.Image.init()
    if src.is_file():
        env = UnityPy.load(str(src))
        for path, obj in env.container.items():
            if obj.type.name not in ["Texture2D", "Sprite"]:
                continue
            if any(path.startswith(f) for f in filters):
                dest = dst.joinpath(*path.split('/'))
                dest.parent.mkdir(parents=True, exist_ok=True)
                data = obj.read()

                if dest.suffix in PIL.Image.EXTENSION and PIL.Image.EXTENSION[dest.suffix] in PIL.Image.SAVE:
                    data.image.save(dest)
                    print(f"{path}=>{dest}")
                else:
                    print(f"{path} type not supported", file=sys.stderr)
    else:
        with ProcessPoolExecutor() as executor:
            for it in src.glob('**/*'):
                if it.is_file():
                    executor.submit(unpack_image_assets, it, dst, filters)


def unpack_text_assets(src: Path, dst: Path, filters: List[str]):
    if src.is_file():
        env = UnityPy.load(str(src))
        for path, obj in env.container.items():
            if obj.type.name not in ["TextAsset"]:
                continue
            if any(path.startswith(f) for f in filters):
                dest = dst.joinpath(*path.split('/'))
                dest.parent.mkdir(parents=True, exist_ok=True)
                data = obj.read()

                with open(dest, "wb") as file:
                    file.write(bytes(data.script))
                print(f"{path}=>{dest}")
    else:
        with ProcessPoolExecutor() as executor:
            for it in src.glob('**/*'):
                if it.is_file():
                    executor.submit(unpack_text_assets, it, dst, filters)


def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(
        dest="command",
        help="commands",
        required=True,
    )
    list_parser = subparsers.add_parser(
        "list",
        help="list assets from sources"
    )
    list_parser.add_argument(
        "src", nargs="+",
        help="Path to source file or directory."
    )
    unpack_image = subparsers.add_parser(
        "unpack-image",
        help="unpack image assets from sources"
    )
    unpack_image.add_argument(
        "src", nargs="+",
        help="Path to source file or directory."
    )
    unpack_image.add_argument(
        "dst",
        help="Path to destination directory."
    )
    unpack_text = subparsers.add_parser(
        "unpack-text",
        help="unpack text assets from sources"
    )
    unpack_text.add_argument(
        "src", nargs="+",
        help="Path to source file or directory."
    )
    unpack_text.add_argument(
        "dst",
        help="Path to destination directory."
    )
    parser.add_argument(
        "-f", "--filter", nargs="+", default=[""],
        help="Specify a path prefix. Only process the assets which match the prefix."
    )
    args = parser.parse_args()
    match args.command:
        case "list":
            for src in args.src:
                list_assets(Path(src), filters=args.filter)
        case "unpack-image":
            for src in args.src:
                unpack_image_assets(Path(src), Path(args.dst), filters=args.filter)
        case "unpack-text":
            for src in args.src:
                unpack_text_assets(Path(src), Path(args.dst), filters=args.filter)


if __name__ == '__main__':
    main()
