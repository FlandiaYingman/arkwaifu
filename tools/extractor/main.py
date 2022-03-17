import argparse
import functools
import sys
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path
from typing import List

import PIL.Image
import UnityPy

# flush every line to prevent blocking outputs
# noinspection PyShadowingBuiltins
print = functools.partial(print)

# initialize PIL to preload supported formats
PIL.Image.preinit()
PIL.Image.init()


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


def unpack(src: Path, dst: Path, filters: List[str], workers=None):
    if src.is_dir():
        print(f"searching files in {src}...")
        with ProcessPoolExecutor(max_workers=workers) as executor:
            for it in src.glob('**/*'):
                if it.is_file():
                    print(f"found {it} in {src}...")
                    executor.submit(unpack, it, dst, filters)
    elif src.is_file():
        env = UnityPy.load(str(src))
        for path, obj in env.container.items():
            if any(path.startswith(f) for f in filters):
                dest = dst.joinpath(*path.split('/'))
                dest.parent.mkdir(parents=True, exist_ok=True)
                data = obj.read()
                if obj.type.name in ["Texture2D", "Sprite"]:
                    if dest.suffix in PIL.Image.EXTENSION and PIL.Image.EXTENSION[dest.suffix] in PIL.Image.SAVE:
                        data.image.save(dest)
                        print(f"{path}=>{dest}")
                    else:
                        print(f"type of {path} is not supported", file=sys.stderr)
                if obj.type.name in ["TextAsset"]:
                    with open(dest, "wb") as file:
                        file.write(bytes(data.script))
                    print(f"{path}=>{dest}")
    else:
        print(f"WARN: {src} is not dir neither file; is skipped")


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
    unpack_parser = subparsers.add_parser(
        "unpack",
        help="unpack image assets from sources"
    )
    unpack_parser.add_argument(
        "src", nargs="+",
        help="Path to source file or directory."
    )
    unpack_parser.add_argument(
        "dst",
        help="Path to destination directory."
    )
    unpack_parser.add_argument(
        "-w", "--workers", nargs="?", default=None,
        help="Specify the concurrency workers count."
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
        case "unpack":
            for src in args.src:
                unpack(Path(src), Path(args.dst), filters=args.filter, workers=int(args.workers) if args.workers else None)


if __name__ == '__main__':
    main()
