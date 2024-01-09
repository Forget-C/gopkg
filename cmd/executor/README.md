## executor用于脚本执行工作。

例如：初始化数据库，初始化配置文件等。

`executor` 将会执行指定目录中的所有脚本。

程序将会在`--base-dir`指定的目录下查找scripts目录，执行scripts目录下的所有脚本。

脚本的执行顺序为：pre_init目录下的脚本 -> init目录下的脚本 -> post-init目录下的脚本。

### 脚本名称
脚本的命名规则：
- 以数字开头，数字越小，越早执行
- "_"用于分割， 例如：001_init.sh -> 002_init.sh -> 003_init.sh
- 以"rollback"结尾的为回滚脚本，用于正常脚本执行失败时的回滚操作，例如：001_init_rollback.sh

### 脚本执行顺序
正常执行的脚本将按照文件名前的数字顺序执行，例如：001-init.sh -> 002-init.sh -> 003-init.sh。

```mermaid
graph LR
    001-init.sh --> 002-init.sh --> 003-init.sh
```

回滚脚本是在正常执行的脚本执行失败时执行的， 将会在失败的序号开始，逆序执行rollback脚本：

如当前执行到 002_init.sh , 此脚本执行失败， 那么将会执行 002_init_rollback.sh -> 001_init_rollback.sh

```mermaid
graph LR
002_init_rollback.sh --> 001_init_rollback.sh
```
### 支持的脚本类型
脚本类型是通过文件后缀判断的， 目前支持 python、sql、bash、二进制可执行文件

### 参数
当有`SQL`类型脚本需要设置环境变量提供连接信息：
- MYSQL_HOST 数据库地址
- MYSQL_PORT 数据库端口
- MYSQL_USER 数据库用户名
- MYSQL_PASSWORD 数据库密码
- MYSQL_DATABASE 数据库名称

`--help`查看帮助
```bash
gopkg|master⚡ ⇒ ./executor --help                 
                                      _
   ___  __  __   ___    ___   _   _  | |_    ___    _ __
  / _ \ \ \/ /  / _ \  / __| | | | | | __|  / _ \  | '__|
 |  __/  >  <  |  __/ | (__  | |_| | | |_  | (_) | | |
  \___| /_/\_\  \___|  \___|  \__,_|  \__|  \___/  |_|
The executor is script run/management tools.
For example: initialize the database, initialize the configuration file, etc.
The program will find the scripts directory in the base-dir directory and execute all scripts in the scripts directory.
The script is executed in the following order: scripts in the pre_init directory-> scripts in the init directory-> scripts in the post-init directory.
The script will be executed according to the number sequence before the file name, for example: 001-init.sh -> 002-init.sh -> 003-init.sh.
If there is no number before the script file name, it will be executed in the dictionary order of the file name, for example: a-init.sh -> b-init.sh -> c-init.sh.
Any script execution fails and the program will exit.

Usage:
  executor [flags]

Flags:
      --base-dir string       The base directory for finding default files. (default "./scripts")
      --enable-post-scripts   Enable post-run scripts. (default true)
      --enable-pre-scripts    Enable pre-run scripts. (default true)
  -h, --help                  help for executor
      --skip-error            Skip errors when exec error.
      --skip-rollback-error   Skip errors when exec rollback error. (default true)
```
