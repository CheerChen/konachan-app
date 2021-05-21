# Konakore

[![Go Report Card](https://goreportcard.com/badge/github.com/CheerChen/konakore)](https://goreportcard.com/report/github.com/CheerChen/konakore)
[![Downloads](https://img.shields.io/github/downloads/CheerChen/konakore/total.svg)](https://github.com/CheerChen/konakore/releases)
[![Release](https://img.shields.io/github/release/CheerChen/konakore.svg?label=Release)](https://github.com/CheerChen/konakore/releases)

Manage your Anime wallpaper collection with Kona-kore.

## Features

- Scan local wallpaper files into album
- New post will be sorted in order of relevance to tags in album
- Easy download new wallpaper and manage them

## Quick Start

We recommend deploying Kona-kore with Docker Compose.

See `Dockerfile` to find more details.

wallpaper files in path should contain valid *id* in filename, otherwise it would not be scaned.

```shell
# example
Konachan.com - 286575 banishment original scenic seifuku signed twintails.jpg
or
286575.jpg
```

## License

Kona-kore is released under the MIT license.

## Contributing

- Ranking based on *tf-idf* algorithm
- Using [Danbooru API (version 1.13.0)](https://konachan.com/help/api)
- Please feel free to contact me if you have any ideas
