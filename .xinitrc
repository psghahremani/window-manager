# xclock &
# xeyes &
xterm &
# xfce4-terminal &
# xterm &
# xfce4-terminal &

# sleep 7 && xterm &

sleep 2
cd ../wmwm && echo -e "Building Go..." && go build -o wm ./*.go && echo -e "Running Go..." && ./wm
