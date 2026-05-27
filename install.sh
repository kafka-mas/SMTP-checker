#!/bin/bash

echo "Start installation"

BINARY_PATH=$PWD/$1
ZABBIX_PATH=/etc/zabbix/zabbix_agentd.d/

if [ "$EUID" -ne 0 ]; then
  echo "Error, user must be root"
  exit 1
fi

if [ ! -d "/etc/zabbix/" ]; then
    echo "Error, /etc/zabbix not found"
	exit 1
fi

if [ ! -d $ZABBIX_PATH ]; then
	mkdir -p $ZABBIX_PATH
fi

touch $ZABBIX_PATH/smtp-checker.conf

cat > $ZABBIX_PATH/smtp-checker.conf << EOF
UserParameter=smtpchecker.master[*], $BINARY_PATH \$1
EOF

cp $PWD/templates/.env-temp $PWD/build/.env
$EDITOR $PWD/build/.env

chmod +x $BINARY_PATH
chown -R zabbix:zabbix $BINARY_PATH

systemctl restart zabbix-agent.service

echo "Sucessfully installed"
