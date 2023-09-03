# mdp-preview &middot; ![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

A Cli tool that allows users to convert markdown files to html and preview them in the default browser.

## Installation & Usage

for now you need to clone the repo

```sh
gh repo clone yacinebenkaidali/fs-walk
```

And then build the project locally

```sh
go build .
```

by the end of the build process you'll have a binary that you can use.

## Features

Here's a list of options that this CLI supports now

- `--root` Root directory to start
- `--list` List files only
- `--ext` File extensions to filter out, comma separated. (exp ".jpeg,.png")
- `--log` files to redirect the logs to.
- `--size` Minimum file size
- `--del` Deleted matched files (use this with care)
- `--archive` Archive directory

## Issues

Please report any issues throgh this [link](https://github.com/yacinebenkaidali/fs-walk/issues). and follow the [bug report](https://github.com/yacinebenkaidali/fs-walk/issues/new?assignees=&labels=type%3AEnhancement&title=) template.

## Licenses

MIT

## Collaborators

- Yacine BENKAIDALI <yacinebenkaidali@gmail.com>
