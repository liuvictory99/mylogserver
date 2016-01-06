#!/bin/sh
clean_days=7
dir=/home/ec2-user/logserver/logs/

scan(){

  if [ ! -d $1 ];then
     return
  fi

  current_day=`date +"%Y%m%d"`
  clean_cmd=1
  cur_dir=$1
  cur_dir=${cur_dir##$dir}

  for((i=$((clean_day+0));$i>-1;i--))
  do
    drop_days=`date -d "$i days ago" +"%Y%m%d"`
    if [ $drop_days = $cur_dir ];then
       clean_cmd=0
       break
    fi
  done

  if [ $clean_cmd -eq 1 ];then
     rm -rf $1
  fi
}

for file in $dir*
do
  scan $file
done
