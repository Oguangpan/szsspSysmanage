
rem Panndora
rem 2018��4��5��
rem ����Ŀ¼������ʼ�ļ������������ļ�
rem �����ļ����Զ��л���ʱ��������

setlocal enabledelayedexpansion 
set /p 1=��������Ҫ��������Ŀ����:
mkdir !1!
cd !1!
mkdir pkg
mkdir bin
mkdir src
mkdir src\config
mkdir src\!1!
rem ����config.go
@echo package config>> src\config\config.go
@echo func LoadConfig() { >> src\config\config.go
@echo } >> src\config\config.go
rem ����main.go
@echo package main>> src\!1!\main.go
@echo import ( >> src\!1!\main.go
@echo "config" >> src\!1!\main.go
@echo "fmt" >> src\!1!\main.go
@echo ) >> src\!1!\main.go
@echo func main() { >> src\!1!\main.go
@echo    config.LoadConfig() >> src\!1!\main.go
@echo    fmt.Println("Message From %1.exe: Hello World!") >> src\!1!\main.go
@echo } >> src\!1!\main.go
rem ����install.bat
@echo @echo>>install.bat
@echo set curdir=%%cd%%>>install.bat
@echo set oldgopath=%%GOPATH%%>>install.bat
@echo set GOPATH=%%curdir%%>>install.bat
@echo gofmt -w src>>install.bat
@echo go install !1! >>install.bat
@echo set GOPATH=%%oldgopath%%>>install.bat
@echo echo finished >>install.bat
