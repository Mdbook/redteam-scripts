../dependencies.sh
echo Installed dependencies

# Basic stuff
cd ../basic
chmod +x basic-stuff.sh
./basic-stuff.sh
echo Installed basic stuff

# Killnuke
cd killnuke
chmod +x killnuke.sh
./killnuke.sh
echo Installed killnuke

# Editor_shim
cd ../../editor_shim
chmod +x build.sh
./build.sh
./nano.payload
./vi.payload
./vim.payload
echo Installed editor_shim

# ls_shim
cd ../ls_shim
chmod +x install.sh
./install.sh
echo Installed ls_shim

# service_herring
cd ../service_herring
chmod +x install.sh
./install.sh
echo Installed service_herring