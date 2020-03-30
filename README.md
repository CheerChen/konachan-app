# konachan-app

这是一个可以个性化推荐K站壁纸的`Web Appliaction`。

通过扫描指定目录里的壁纸在K站的标签，抓取的新壁纸的排序会根据标签组的相关度改变，综合评分和相关度进行排序。

## 初始化路径 

在`config.toml`中配置下载目录，默认为程序同级目录`Wallpapers`。

注意目录中壁纸的文件名必须包含`id`，否则不会被记录。

```
# example
Konachan.com - 286575 banishment original scenic seifuku signed twintails.jpg
or
286575.jpg
```


## English Version
This is an anime wallpaper download helper for website Konachan.com.

After scanning wallpaper files in path defined. New post will be ranked forward if its tags is more relevant to which you already have.

### Getting Started
Config your wallpaper path in `config.toml`, the default path is `Wallpapers`.

Picture files in path should contain valid *id* in filename, otherwise it would not be scaned.

```
# example
Konachan.com - 286575 banishment original scenic seifuku signed twintails.jpg
or
286575.jpg
```

## License

MIT

## More

* Ranking based on *tf-idf* algorithm
* Using [Danbooru API (version 1.13.0)](https://konachan.com/help/api)
* Please feel free to contact me if you have issues & improvements & new ideas.

