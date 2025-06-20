# 项目名称

辅助记忆问答工具

## 项目背景

该项目的研发初衷主要是用于个人学习记忆使用，根据基于约定模板记录的笔记生成问答题，应用程序是命令行终端的形式。

## 项目结构

```
.
├── dissove/            # 解析笔记文档为符合结构的文本块树
├── quiz/               # 核心业务逻辑
│   ├── engine.go       # 问答引擎逻辑
│   └── parse.go        # 生成问答题逻辑
├── /util               # 工具包
├── main.go             # 主程序入口
├── meta.go             # 元数据结构定义
├── quiz.go             # 问答提示及插件功能注入逻辑
├── scan.go             # 根据正则表达式匹配抽取文本块的逻辑
├── store.go            # 将可构成树的文本块结构json串保存到数据库中
├── show_type.go        # 按行展示的对照打字功能，类似于金山打字，不校检用户打的字。
├── README.md           # 项目说明文档
├── bat.md              # 用于生成问题内容的笔记示例
├── english.md          # 用于对照打字的英文笔记示例
└── go.mod              # Go 依赖管理文件
```

## 使用方法

### 笔试文件说明

笔记类型为markdown类型，请看项目中的<kbd>bat.md</kbd>示例
其中需要注意的是在笔记中需要记忆的内容分为提示和答案两部分，这两部分用" -- "连接
其它的按markdown文件要求即可，


### 关于可执行文件
可执行文件是个命令行终端工具，如编译完可执行文件名称是<kbd>feelorder.exe</kbd>,使用示例如下
```
feelorder.exe -scan  [filename]  //尝试扫描抽取笔记内容，看看是否出问题
feelorder.exe -test  [filename]  //扫描抽取笔记内容并进入记忆测试
feelorder.exe -store  [filename] //扫描抽取笔记内容并保存到数据库中
feelorder.exe -load  [filename]  //从数据库中加载处理好的指定笔记内容并进入测试
feelorder.exe -list  [filename]  //从数据库中展示保存的笔记列表，用户可从中选择并进入测试
feelorder.exe -show  [dirname]   //从指定文件夹中展示需要进行打字记忆的文件并进入训练

```

### 关于保存笔记内容的数据库表
```
CREATE TABLE `notebook` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `prev_id` int(11) DEFAULT NULL,
  `note_name` varchar(20) DEFAULT NULL COMMENT '笔记名称',
  `content` mediumtext COMMENT '笔记内容，json格式',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19139 DEFAULT CHARSET=utf8 COMMENT='笔记'

```


### 使用示例

![使用示例](https://raw.githubusercontent.com/BryceSun/images/refs/heads/main/bryceuseway.png)


### 关于打字记忆训练功能

没什么特别的功能，主要是按行展示文本中的内容，用户进行跟打，不作对错校验

![文件选择](https://raw.githubusercontent.com/BryceSun/images/refs/heads/main/type_list.png)


![练习示例](https://raw.githubusercontent.com/BryceSun/images/refs/heads/main/type_way.png)