# konachan-app

这是一个帮助你在 Konachan.com 获得自己感兴趣的~~色图~~壁纸的工具。

This is an anime wallpaper download helper for website Konachan.com.

通过导入已经下载的 Konachan 壁纸，或者使用 konachan-app 下载一段时间，konachan-app 的壁纸的排序会根据你的壁纸库的相关度改变，越相关越靠前。

After importing wallpaper files, or just start using konachan-app to download a few times, the order of posts in konachan-app will be changed. One post will be ranked forward if its tags is more relevant to which you have downloaded.

## Getting Started

### 初始化相册路径 

启动 konachan-app 后，在命令行输入你的壁纸目录，可以是空目录。

Input your wallpaper path after running konachan-app，can be an empty directory.

```
Please input wallpaper path >E:\Wallpaper
```

注意目录中壁纸的文件名必须包含*id*，否则不会被记录。

The picture files in path should contain valid *id* in filename, otherwise it would not be recorded.

```
Konachan.com - 286575 banishment original scenic seifuku signed twintails.jpg
or
286575.jpg
```

### 创建快捷方式

不想每次都输入地址，可以通过快捷方式启动。

You can create a shortcut to avoid inputting path each time.

右键 - 新建快捷方式 - 输入

Right click - new shortcut - input

```
konachan-app.exe -p "E:\Wallpaper"
```

### 

## License

MIT

## More

* Ranking based on *tf-idf* algorithm
* Using [Danbooru API (version 1.13.0)](https://konachan.com/help/api)
* Please feel free to contact me if you have issues & improvements & new ideas.