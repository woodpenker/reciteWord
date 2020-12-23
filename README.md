命令行背单词
---

该小程序会随机抽取数据集中的单词, 让你选择正确含义, 直到所有单词都选择对.

最终完成背单词后可查看记录的效果:
![](screenshot/1.png?raw=true)
背单词过程中效果:
![](screenshot/2.png?raw=true)

## 使用
编译:
`go build word.go`

使用:
```shell
$ word --help
-f string
    the data file path (default "我的生词本.json")
-l int
    the max byte to show one line (default 180)
-n int 
    the number of words to recite (default 20)
-t string
    the data file type (default "json")
```

* 使用`-f` 指定自己的数据文件 默认:我的生词本.json
* 使用`-l` 设置一行最大显示字符数量(用于暂时缓解换行导致的显示问题) 默认180
* 使用`-n` 设置需要背诵的单词数量, 默认20
* 使用`-t` 设置自己的数据文件格式, 默认json
