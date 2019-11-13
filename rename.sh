#!/bin/bash

# 通过此脚本直接将 go-simple 重命名为自己的服务工程
repository=$1
oldrepository=${PWD##*/}
if [ "$repository" = "" ] ;then
    echo "没有输入新的項目名"
    exit
fi
echo "項目名重命名 $oldrepository -> $repository"

# 1, 修改go代码中的所有当前mod的引用
function rename_go() {
	folders=`ls`
    for folder in ${folders[@]};do
		if [ -d "$(pwd)/$folder" ] ;then
            cd ./$folder
            rename_go
            cd ../
        else
            if [[ $folder == *.go ]] ;then
                if [[ `cat $(pwd)/$folder | grep "\"$oldrepository"` != "" ]]; then
                    echo "sed 's/$oldrepository/$repository/g' $(pwd)/$folder"
                    sed 's/$oldrepository/$repository/g' $(pwd)/$folder >> $(pwd)/$folder.tmp
                fi
            fi
        fi
	done
}

cd ./src
    rename_go
cd ../
