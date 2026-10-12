# Android Telegram Bot

Bot Telegram untuk mengelola `sbfr`, dashboard API Mihomo/YACD, dan core proxy di perangkat Android. Bot ini ditujukan untuk perangkat Android yang sudah di-root dan hanya menerima perintah dari Telegram ID pemilik yang dikonfigurasi.

## Fitur dan perintah

| Perintah | Fungsi |
| --- | --- |
| `/help` | Menampilkan bantuan |
| `/menu` | Membuka menu `sbfr` |
| `/sbfr` | Membuka menu pengelolaan `sbfr` |
| `/yacd` | Mengelola proxy melalui API Mihomo |
| `/core` | Memilih core: Clash, Sing-box, Xray, V2Fly, atau Hysteria |
| `/speedtest` | Menjalankan SpeedTest |
| `/myip` | Menampilkan informasi IP publik |
| `/info` | Menampilkan informasi baterai dan uptime Android |
| `/ipinfo <domain>` | Melihat informasi IP/domain; bisa juga membalas pesan berisi domain |
| `/hostip <domain>` | Mencari alamat IP domain; bisa juga membalas pesan berisi domain |
| `/import [direktori]` | Membalas dokumen dengan perintah ini untuk menyimpannya; direktori default `/data/adb/box` |
| `/export <nama-file>` | Mengirim file dari `/data/adb/box` |
| `/log` | Membaca `/data/adb/box/run/runs.log` |
| `/sh <perintah>` | Menjalankan perintah shell |

**Perhatian:** `/sh` menjalankan perintah shell yang diberikan. Gunakan hanya di chat tepercaya dan berhati-hatilah dengan perintah yang dijalankan.

## Persyaratan

- Perangkat Android yang sudah di-root dengan `sbfr` dan berkas terkait tersedia di lokasi yang diharapkan bot.
- Go 1.25 atau lebih baru untuk membangun aplikasi.
- Token bot Telegram dan ID numerik Telegram pemilik.
- API Mihomo yang dapat dijangkau dari perangkat (untuk fitur YACD).

## Instalasi dan konfigurasi

1. Clone repository dan masuk ke direktorinya:

   ```sh
   git clone <REPO_URL>
   cd <REPO_FOLDER>
   ```

2. Salin konfigurasi contoh, lalu isi token dan ID Telegram Anda:

   ```sh
   cp docs/bot.ini bot.ini
   ```

   ```ini
   [bot]
   token = <YOUR_BOT_TOKEN>
   owner = <YOUR_TELEGRAM_ID>

   [mihomo]
   api = http://127.0.0.1:9090
   secret = <MIHOMO_API_SECRET>
   ```

   `owner` harus berupa ID numerik Telegram. `secret` boleh dikosongkan jika API Mihomo tidak memerlukan autentikasi.

3. Bangun aplikasi dengan Go untuk arsitektur dan sistem operasi perangkat tujuan:

   ```sh
   go build -o bot .
   ```

4. Jalankan bot dengan hak akses root dan berikan lokasi berkas konfigurasi:

   ```sh
   su -c './bot -c ./bot.ini'
   ```

   Jika opsi `-c` tidak diberikan, bot mencari `bot.ini` di direktori tempat binary berada.

## Struktur repository

```text
.
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── main.go
├── module/
│   ├── commands.go
│   ├── hostip.go
│   ├── info.go
│   ├── ipinfo.go
│   ├── myip.go
│   ├── service.go
│   ├── speedtest.go
│   └── yacd.go
└── docs/
    ├── bot.ini
    └── bot.sh
```
