## executor.go
`executor.go`文件中定义了脚本执行器相关的结构体和方法。

`MakeExecutor`方法会根据文件后缀创建对应的脚本执行器，目前支持的脚本类型有：`python`、`sql`、`bash`、`二进制可执行文件`。

你可以实现自己的`MakeExecutor`和`Executor`来支持更多脚本类型

## finder.go
`finder.go`文件中定义了脚本查找器和管理器相关的结构体和方法。

### Scan
`Scan`方法会扫描指定目录下的所有脚本文件，并按照顺序返回脚本列表。

脚本的命名规则：
- 以数字开头，数字越小，越早执行
- "_"用于分割， 例如：001_init.sh -> 002_init.sh -> 003_init.sh
- 以"rollback"结尾的为回滚脚本，用于正常脚本执行失败时的回滚操作，例如：001_init_rollback.sh

### ScriptManager
`ScriptManager`提供脚本管理的相关方法，用于执行脚本。

`ScriptManager` 将脚本定义为三个阶段， `pre`、`run`和`post`

### Prepare
`ScriptManager.checkers`用于检查是否满足脚本执行条件以及初始化一些数据到`stat`中

`ScriptManager.getters`用于从`stat`中获取指定数据

