#!/usr/bin/env bash

export SYBASE=/opt/sybase
source /opt/sybase/SYBASE.sh
# 清除文件末尾空行
sed -i ":a; /^\n*$/{\$d; N; ba}" $SYBASE/$SYBASE_ASE/install/RUN_MYSYBASE

# https://github.com/DataGrip/docker-env/issues/13#issuecomment-1741782611
echo "-T11889\\" >> $SYBASE/$SYBASE_ASE/install/RUN_MYSYBASE

# 设置SSL
if [ "${ENABLE_SSL}" == true ]; then
    cat /MYSYBASE.crt /MYSYBASE.key > $SYBASE/$SYBASE_ASE/certificates/MYSYBASE.crt
    cat /MYSYBASE.crt > $SYBASE/$SYBASE_ASE/certificates/MYSYBASE.txt
    cat /MYSYBASE.crt >> $SYBASE/config/trusted.txt
    sed -i 's/enable ssl = DEFAULT/enable ssl = 1/' $SYBASE/$SYBASE_ASE/MYSYBASE.cfg
    sed -i 's/5000/5000 ssl/g' $SYBASE/interfaces
fi

sh /opt/sybase/SYBASE.sh && sh /opt/sybase/ASE-16_0/install/RUN_MYSYBASE > /dev/null &

#waiting for sybase to start
export STATUS=0
i=0
echo "STARTING... (about 30 sec)"
while [[ $STATUS -eq 0 ]] || [[ $i -lt 30 ]]; do
        sleep 1
        i=$((i+1))
        STATUS=$(grep "server name is 'MYSYBASE'" /opt/sybase/ASE-16_0/install/MYSYBASE.log | wc -c)
done

SYBASE_LOGIN_ROOT=root
SYBASE_LOGIN_TEST=test
# 密码开头结尾不能有空白符
SYBASE_PASSWORD="Pass}word"
cat <<-EOSQL > init1.sql
use master
go
disk resize name='master', size='60m'
go
create login $SYBASE_LOGIN_ROOT with password "$SYBASE_PASSWORD"
go
create login $SYBASE_LOGIN_TEST with password "$SYBASE_PASSWORD"
go
grant role sso_role to $SYBASE_LOGIN_ROOT
go
sp_adduser $SYBASE_LOGIN_TEST,$SYBASE_LOGIN_TEST
go
EOSQL

/opt/sybase/OCS-16_0/bin/isql -Usa -PmyPassword -SMYSYBASE -i"./init1.sql"

echo "Sybase initialized successfully"

#trap
while [ "$END" == '' ]; do
  sleep 1
  trap "/etc/init.d/sybase stop && END=1" INT TERM
done
