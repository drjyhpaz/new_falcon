# 🔓 Falcon RDP Brute-Force System

سیستم حمله Brute-Force تخصصی برای RDP با معماری توزیع‌شده، UI حرفه‌ای و ویژگی‌های پیشرفته.

## 🎯 ویژگی‌های اصلی

### ۱. معماری قدرتمند
- **Worker Pool** قابل تنظیم برای Concurrency بالا
- **Goroutine** و **Channel** برای I/O غیرهمزمان
- **Context**-based لغو عملیات و مدیریت Timeout
- بهینه‌سازی Multi-Core

### ۲. مدیریت ورودی�n- خواندن `servers.txt` (IP:Port)
- خواندن `users.txt` (نام‌های کاربری)
- خواندن `passwords.txt` (رمزها)
- تولید **کارتزینی** Credential ترکیبی
- استخراج دامنه خودکار (DOMAIN\user یا user@domain)

### ۳. موتور حمله RDP
- **Password Spraying** (ایمن‌ترین روش)
- **Credential Stuffing** (تست همه رمزها روی یک کاربر)
- **Hybrid Mode**
- احراز هویت واقعی با **grdp** (NLA، CredSSP، TLS)
- کنترل نرخ (PPS) و Rate Limiting

### ۴. شناسایی و هوشمندی (Recon)
- **Pre-Attack Recon** خودکار و همزمان
- بررسی باز بودن RDP (X.224)
- تشخیص NLA و SSL
- تشخیص نسخه Windows
- اندازه‌گیری Latency
- طبقه‌بندی خطاها و تغییر استراتژی

### ۵. ضدتشخیص و مخفی‌کاری
- **Stealth Mode** با تاخیر تصادفی (Jitter)
- **Adaptive Rate Limiting**
- **Low & Slow** حملات طولانی‌مدت
- پشتیبانی **Proxy**: SOCKS5, HTTP, TOR
- **IP Rotation** خودکار

### ۶. جلوگیری از قفل شدن
- شمارش تلاش‌های ناموفق به صورت (User, Target)
- توقف موقت برای کاربر وقتی به آستانه رسید
- **Cooldown Period**
- تنظیم خودکار به Password Spraying

### ۷. مدیریت وضعیت
- **Checkpoint** ذخیره‌سازی مرتب
- **Resume** از نقطه قطع
- جلوگیری از تکرار تلاش‌ها
- ذخیره state.json

### ۸. گزارش‌گیری و لاگینگ
- لاگینگ لحظه‌ای با سطوح (INFO، SUCCESS، WARNING، ERROR)
- ذخیره نتایج موفق در `results.json`
- گزارش نهایی (JSON و CSV)
- Session Storage

### ۹. رابط کاربری (Fyne)
- **Dashboard Tab**: Start/Stop، PPS، آمار، Progress Bar
- **Files Tab**: انتخاب فایل‌ها، تولید Credential
- **Settings Tab**: Threads، Timeout، Stealth، Proxy، Resume
- **Results Tab**: جدول نتایج، جستجو و فیلتر
- **Recon Tab**: اطلاعات Pre-Recon

## 📋 ساختار پروژه

```
falcon_rdp/
├── main.go
├── go.mod
├── go.sum
├── config/
│   ├── config.go
│   └── types.go
├── attack/
│   ├── engine.go
│   ├── worker.go
│   └── strategies.go
├── rdp/
│   ├── auth.go
│   ├── detector.go
│   └── recon.go
├── credentials/
│   ├── loader.go
│   └── generator.go
├── proxy/
│   ├── manager.go
│   └── rotation.go
├── evasion/
│   ├── stealth.go
│   └── jitter.go
├── state/
│   ├── checkpoint.go
│   └── resume.go
├── ui/
│   └── dashboard.go
├── logger/
│   └── log.go
└── utils/
    └── helpers.go
```

## 🚀 نحوه استفاده

### ۱. نصب
```bash
git clone https://github.com/falconjonz/falcon_rdp.git
cd falcon_rdp
go mod download
go build -o falcon_rdp
```

### ۲. فایل‌های ورودی

**servers.txt**
```
192.168.1.1:3389
192.168.1.2:3389
10.0.0.5:3389
```

**users.txt**
```
administrator
Admin
user
guest
```

**passwords.txt**
```
Password123!
Admin123
test
default
```

### ۳. اجرا
```bash
./falcon_rdp
```

سپس از رابط کاربری استفاده کنید.

## 📊 نتایج و گزارش

### results.json
```json
[
  {
    "ip": "192.168.1.1",
    "port": 3389,
    "username": "administrator",
    "password": "Password123!",
    "domain": "",
    "timestamp": "2024-01-15T10:30:45Z",
    "os_info": {},
    "admin": false
  }
]
```

### state.json (برای Resume)
```json
{
  "last_credential_index": 245,
  "last_target_index": 3,
  "targets": [...],
  "successful_logins": [...],
  "total_attempts": 1230,
  "total_successes": 5
}
```

## ⚙️ تنظیمات

### معماری
- **Threads**: خودکار (CPU*2) یا دستی
- **Timeout**: 10s پیش‌فرض

### ضدتشخیص
- **Stealth Mode**: فعال/غیرفعال
- **Jitter Min/Max**: تاخیر تصادفی (ms)
- **Adaptive Rate**: کاهش نرخ خودکار

### Proxy
- **Type**: SOCKS5، HTTP، TOR
- **File**: لیست Proxy از فایل
- **Rotation**: Round-Robin

### امنیت
- **Insecure TLS**: برای گواهی نامعتبر
- **NLA Detection**: تشخیص خودکار

## 📈 Performance

- **بیش از ۱۰۰۰ تلاش در ثانیه** (بدون Proxy)
- **مقیاس‌پذیری** تا ۱۰۰+ هدف همزمان
- **کم مصرف** با Streaming فایل‌های بزرگ

## 📝 لایسنس

MIT License

## ⚠️ تنبیه قانونی

این ابزار **فقط برای اهداف آموزشی و تست نفوذ مجاز** است. استفاده در سیستم‌های غیرمجاز **غیرقانونی** و **غیراخلاقی** است.

---

**توسعه‌دهنده**: falconjonz  
**نسخه**: 1.0.0  
**آخرین بروزرسانی**: 2024
