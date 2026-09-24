# Android Telegram Bot

Bot Telegram untuk kontrol `/system/bin/sbfr`, dashboard YACD, serta manajemen core (Clash, Xray, Sing-box, V2Fly, Hysteria) di Android.

---

## 📦 Fitur

- `/help`       - Menampilkan bantuan
- `/menu`       - Menampilkan menu utama
- `/status`     - Ringkasan status sistem Android (battery, uptime, memori)
- `/health`     - Cek status proses core/Mihomo dan koneksi API
- `/log`        - Menampilkan log terakhir dari runs.log
- `/import`     - <path> (default: /data/adb/box/)
- `/export`     - <path/file> (default: /data/adb/box/)
- `/sbfr`       - Menu untuk `/system/bin/sbfr`
- `/yacd`       - Menu dashboard YACD
- `/core`       - Pilih core untuk settings.ini
- `/speedtest`  - Jalankan speed test

---

## ⚙️ Instalasi

1. Clone repository:

```bash
git clone <REPO_URL>
cd <REPO_FOLDER>
```

2. Install dependensi Go:

```bash
go mod tidy
```

3. Buat konfigurasi bot di `bot.ini`:

```ini
[bot]
token = <YOUR_BOT_TOKEN>
owner = <YOUR_TELEGRAM_ID>

[mihomo]
api = http://127.0.0.1:9090
secret = 123456
```

Catatan: untuk kompatibilitas, format lama berikut juga masih diterima:

```ini
[bot]
token = <YOUR_BOT_TOKEN>
owner = <YOUR_TELEGRAM_ID>
mihomo_api = http://127.0.0.1:9090
api_secret = 123456
```

4. Build bot:

```bash
go build -o bot .
```

5. Jalankan bot:

```bash
su -c 'sh bot -c /path/to/bot.ini'
```

Atau bila file bot.ini berada di folder binary yang sama:

```bash
su -c sh bot -c bot.ini
```

> `-c` bersifat opsional. Kalau tidak diberikan, bot akan mencari `bot.ini` di folder executable.

---

## 🔧 Fitur Tambahan

- `health watcher` otomatis: bot akan mengirim notifikasi ke owner jika proses core atau API tidak sehat.
- `status` dashboard: menampilkan ringkasan sistem Android untuk debugging cepat.
- `log` tailing: menampilkan bagian terakhir log untuk diagnosis singkat.

---

## Struktur Folder

```bash
bot/
├── README.md
├── main.go
├── go.mod
├── go.sum
├── LICENSE
├── docs/
│   ├── bot.ini
│   └── bot.sh
├── module/
│   ├── commands.go
│   ├── status.go
│   ├── health.go
│   ├── yacd.go
│   ├── service.go
│   ├── speedtest.go
│   ├── myip.go
│   ├── ipinfo.go
│   ├── hostip.go
│   └── info.go
└── .github/
```
