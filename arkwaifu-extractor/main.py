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
                print(f"{path}", flush=True)
    else:
        for it in src.glob('**/*'):
            if it.is_file():
                list_assets(it, filters)


def unpack(src: Path, dst: Path, filters: List[str]):
    PIL.Image.preinit()
    PIL.Image.init()
    if src.is_file():
        env = UnityPy.load(str(src))
        for path, obj in env.container.items():
            if any(path.startswith(f) for f in filters):
                dest = dst.joinpath(*path.split('/'))
                dest.parent.mkdir(parents=True, exist_ok=True)
                data = obj.read()
                if obj.type.name in ["Texture2D", "Sprite"]:
                    if dest.suffix in PIL.Image.EXTENSION and PIL.Image.EXTENSION[dest.suffix] in PIL.Image.SAVE:
                        data.image.save(dest)
                        print(f"{path}=>{dest}", flush=True)
                    else:
                        print(f"{path} type not supported", file=sys.stderr, flush=True)
                if obj.type.name in ["TextAsset"]:
                    with open(dest, "wb") as file:
                        file.write(bytes(data.script))
                    print(f"{path}=>{dest}", flush=True)
    else:
        print("Searching files...", flush=True)
        with ProcessPoolExecutor() as executor:
            for it in src.glob('**/*'):
                if it.is_file():
                    executor.submit(unpack, it, dst, filters)


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
                unpack(Path(src), Path(args.dst), filters=args.filter)


if __name__ == '__main__':
    main()
