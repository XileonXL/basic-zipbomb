# basic-zipbomb

A zipbomb is just an amount of zip files automtically generated and compressed together starting from a txt file filled of zeros.

Using this python basic script you can generate a zipbomb with the following command:

```sh
python3 main.py --size 1
```

The valid arguments are:

- `size` (required) --> Size of the initial file in MB. The total decompressed size of the zipbomb will be `copies ** levels * size`.
- `levels` --> Number of nested ZIP levels to create (default: 3).
- `copies` --> Number of copies to include in each level (default: 3).
