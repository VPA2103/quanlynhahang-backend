# quanlynhahang-backend

Backend quản lý nhà hàng (Go, Gin). Module: `github.com/vpa/quanlynhahang-backend`.

## Cấu trúc thư mục

```
quanlynhahang-backend/
├── config/                 # Cấu hình app, env, Cloudinary, thanh toán
│   ├── cloudinary.go
│   ├── config.go
│   ├── env.go
│   └── payment.go
├── controllers/            # HTTP handlers / controller
│   ├── auth_handler.go
│   ├── ban_an.go
│   ├── dat_ban_controller.go
│   ├── goi_mon.go
│   ├── hoa_don.go
│   ├── lien_he_controller.go
│   ├── loai_mon_an.go
│   ├── mon_an.go
│   ├── nhan_vien.go
│   ├── notification_controller.go
│   ├── payment.go
│   └── upload_handler.go
├── middleware/             # Auth, phân quyền
│   ├── auth.go
│   ├── auth_roles.go
│   └── role.go
├── models/                 # Model dữ liệu / entity
│   ├── banan.go
│   ├── chitiet_hoadon.go
│   ├── contact.go
│   ├── datban.go
│   ├── hoadon.go
│   ├── Images.go
│   ├── khachhang.go
│   ├── lien_he.go
│   ├── loaimonan.go
│   ├── monan.go
│   ├── nhanvien.go
│   ├── notification.go
│   └── thanhtoan.go
├── realtime/               # WebSocket / realtime hub
│   ├── client.go
│   ├── handler.go
│   ├── hub.go
│   └── push.go
├── routes/                 # Đăng ký route API
│   ├── admin/
│   │   └── admin_Routes.go
│   ├── realtime/
│   │   └── websocket.go
│   ├── banan_routes.go
│   ├── dat_ban_api.go
│   ├── goi_mon.go
│   ├── hoa_don.go
│   ├── lien_he_route.go
│   ├── loai_mon_an.go
│   ├── mon_an.go
│   ├── nhanvien_Routes.go
│   ├── payment_routes.go
│   ├── router.go
│   └── upload.go
├── send_mail/              # Gửi email
│   └── send_mail.go
├── services/               # Logic nghiệp vụ
│   ├── goi_mon.go
│   └── hoa_don.go
├── utils/                  # JWT, mail, QR, upload, secret
│   ├── jwt.go
│   ├── mail.go
│   ├── qr.go
│   ├── secret.go
│   └── upload_image.go
├── Document/               # Tài liệu (thư mục trong repo)
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go                 # Điểm vào ứng dụng
├── test.html
├── .env.example            # Mẫu biến môi trường
└── .gitignore
```
