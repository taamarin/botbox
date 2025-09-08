#!/system/bin/sh

scripts=$(realpath "$0")
scripts_dir=$(dirname "${scripts}")

if ! command -v busybox &> /dev/null; then
  export PATH="/data/adb/magisk:/data/adb/ksu/bin:/data/adb/ap/bin:$PATH:/system/bin"
fi

BOT_EXEC="$scripts_dir/bot -c $scripts_dir/bot.ini"       # Path ke executable bot
PID_FILE="$scripts_dir/bot.pid"

# Fungsi: cek apakah bot sedang berjalan
check_bot() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null 2>&1; then
            return 0
        else
            rm -f "$PID_FILE"
            return 1
        fi
    else
        return 1
    fi
}

# Fungsi: start bot
start_bot() {
    chmod +x $scripts_dir/bot
    if check_bot; then
        echo "Bot sudah berjalan dengan PID $(cat $PID_FILE)"
    else
        nohup busybox setuidgid 0:0 $BOT_EXEC >/dev/null 2>&1 &
        PID=$!
        echo -n $PID > $PID_FILE
        echo "Bot berhasil dijalankan dengan PID $(cat $PID_FILE)"
    fi
}

# Fungsi: stop bot
stop_bot() {
    if check_bot; then
        PID=$(cat $PID_FILE)
        su -c "kill $PID"
        rm -f "$PID_FILE"
        echo "Bot dengan PID $PID berhasil dihentikan"
    else
        echo "Bot tidak berjalan"
    fi
}

# Fungsi: status bot
status_bot() {
    if check_bot; then
        PID=$(cat $PID_FILE)
        echo "Bot sedang berjalan dengan PID $PID"
        su -c "ps -p $PID -o cmd,pid,pcpu,%mem,rss"
    else
        echo "Bot tidak berjalan"
    fi
}

# Fungsi: restart bot
restart_bot() {
    stop_bot
    sleep 1
    start_bot
}

# Menu interaktif
case $1 in
    start|r) start_bot ;;
    stop|k) stop_bot ;;
    restart|l) restart_bot ;;
    status|s) status_bot ;;
    help)
        echo "Available commands:"
        echo "  start   - Menjalankan bot"
        echo "  stop    - Menghentikan bot"
        echo "  restart - Restart bot"
        echo "  status  - Cek status bot"
        ;;
    *) echo "Opsi tidak valid! Gunakan '$0 help'" ;;
esac