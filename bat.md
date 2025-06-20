- 显示信息 -- echo
- 注释语句 -- rem
- 注释符号 -- ::
- 列文件名 -- dir
- 目录切换 -- cd
- 切换到上层目录 -- cd..
- 切换到D盘 -- D:
- 引用当前完整路径 -- %cd%
- 
- 关闭命令回显 -- @+命令
- 打开回显 -- echo on
- 关闭回显 -- echo off
- 输出空行 -- echo.
- 暂停 -- pause
- 回复命令提问 -- echo xx|命令
```
@echo off
echo Y|rd /s c:\abc
pause
```
- 将内容输出到文件 -- echo xx > filename
