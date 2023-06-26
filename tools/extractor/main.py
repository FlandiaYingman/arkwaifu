import argparse
import concurrent.futures
import functools
import json
import os.path
import sys
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path
from typing import Any
from typing import Dict
from typing import List

import PIL.Image
import UnityPy
from UnityPy.classes import Object
from UnityPy.classes import PPtr
from UnityPy.classes.Object import NodeHelper

# flush every line to prevent blocking outputs
# noinspection PyShadowingBuiltins
print = functools.partial(print, flush=True)

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


def unpack(src: Path, dst: Path, filters: List[str], workers: int | None, tasks_per_child: int | None):
    if src.is_dir():
        print(f"searching files in {src}...")
        with ProcessPoolExecutor(max_workers=workers, max_tasks_per_child=tasks_per_child) as executor:
            futures = []
            for subSrc in src.glob('**/*'):
                if subSrc.is_file():
                    print(f"found {subSrc} in {src}...")
                    futures.append(executor.submit(unpack, subSrc, dst, filters, None, None))
            (done_futures, _) = concurrent.futures.wait(futures, return_when='FIRST_EXCEPTION')
            for future in done_futures:
                future: concurrent.futures.Future
                e = future.exception()
                if e is not None:
                    print(f"exception occurs during concurrently unpacking: {e}", file=sys.stderr)
    elif src.is_file():
        env = UnityPy.load(str(src))
        for container, obj_reader in env.container.items():
            if any(container.startswith(f) for f in filters):
                obj = obj_reader.read()
                container = container.lower()
                obj_name = obj.name.lower()
                container_path = os.path.normpath(os.path.join(container, '..', obj_name))
                path_id_path = dst / os.path.normpath(os.path.join(container, '..', f"{obj_name}.json"))

                path_id_dict = export(obj, dst, container, container_path)
                if len(path_id_dict) > 0:
                    with open(path_id_path, "w", encoding="utf8") as file:
                        json.dump(path_id_dict, file, ensure_ascii=False, indent=4)
    else:
        print(f"WARN: {src} is not dir neither file; skipping")


def export(obj: Object, dst: Path, obj_container: str | None, container_path: str,
           path_id_dict: Dict[int, str] = None) -> Dict[int, str]:
    path_id_dict = {} if path_id_dict is None else path_id_dict

    obj_name = getattr(obj, 'name', '')
    obj_path = obj_container or f"{container_path}/{obj_name}"

    if obj_container is not None:
        obj_container = obj_container.lower()
    container_path = container_path.lower()
    obj_name = obj_name.lower()
    obj_path = obj_path.lower()

    if obj.type.name in ["Texture2D", "Sprite"]:
        dest = dst / obj_container if obj_container else dst / container_path / f"{obj_name}.png"
        dest.parent.mkdir(parents=True, exist_ok=True)
        path_id_dict[obj.path_id] = str(dest.name)
        if dest.suffix in PIL.Image.EXTENSION and PIL.Image.EXTENSION[dest.suffix] in PIL.Image.SAVE:
            obj.image.save(dest)
            print(f"{obj_path}({obj.type.name})=>{dest}")
        else:
            print(f"cannot export {obj_path}({obj.type.name}), format is not supported", file=sys.stderr)

    if obj.type.name in ["TextAsset"]:
        dest = dst / obj_container if obj_container else dst / container_path / f"{obj_name}.txt"
        dest.parent.mkdir(parents=True, exist_ok=True)
        path_id_dict[obj.path_id] = str(dest.name)
        with open(dest, "wb") as file:
            file.write(bytes(obj.script))
        print(f"{obj_path}({obj.type.name})=>{dest}")

    if obj.type.name in ["MonoBehaviour"]:
        script = obj.m_Script.read()
        obj_name = script.name
        dest = dst / obj_container if obj_container else dst / container_path / f"{obj_name}.json"
        dest.parent.mkdir(parents=True, exist_ok=True)
        path_id_dict[obj.path_id] = str(dest.name)
        with open(dest, "w", encoding="utf8") as file:
            json.dump(obj.read_typetree(), file, ensure_ascii=False, indent=4)
        print(f"{obj_path}({obj.type.name})=>{dest}")

    if obj.type.name in ["GameObject"]:
        nodes = traverse(obj)
        container_path = os.path.normpath(os.path.join(container_path, '..', obj_name))
        for node in nodes:
            export(node, dst, None, container_path, path_id_dict)

    return path_id_dict


def traverse(obj: Object) -> List[Object]:
    """
    traverse() traverses through an UnityPy Object, returning all sub UnityPy Objects.
    """

    def traverse_tree(o: Any, r: List[Object]):
        if not o:
            return
        if isinstance(o, PPtr):
            o = o.read()
            r += [o]

        skipping_attr_names = ["m_GameObject", "m_Father", "m_Script"]
        match o:
            case list():
                for i, attr in enumerate(o):
                    traverse_tree(attr, r)
            case dict():
                for name, attr in o.items():
                    traverse_tree(attr, r)
            case NodeHelper():
                for name, attr in o.items():
                    if name in skipping_attr_names:
                        continue
                    traverse_tree(attr, r)
            case _:
                attr_names = dir(o)
                attr_names = list(filter(lambda x: x.startswith('m_') or x == "type_tree", attr_names))
                for attr_name in attr_names:
                    if attr_name in skipping_attr_names:
                        continue
                    sub_obj = getattr(o, attr_name)
                    traverse_tree(sub_obj, r)

    results = []
    traverse_tree(obj, results)
    return results


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
    unpack_parser.add_argument(
        "-t", "--tasks_per_child", nargs="?", default=None,
        help="Specify the tasks per child workers count."
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
                unpack(Path(src), Path(args.dst), filters=args.filter,
                       workers=int(args.workers) if args.workers else None,
                       tasks_per_child=int(args.tasks_per_child) if args.tasks_per_child else None)


if __name__ == '__main__':
    main()
