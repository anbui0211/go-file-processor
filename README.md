# Go File Processor

## Giới thiệu

Go File Processor là hệ thống backend hỗ trợ xử lý dữ liệu kế toán, tập trung vào **journal voucher** với các tính năng chính:

- **Import dữ liệu** từ file CSV/Excel
- **Export dữ liệu** ra file CSV/Excel
- **Xử lý batch job** để đảm bảo hiệu suất cao

Hệ thống sử dụng **AWS SQS** để quản lý queue và **S3** để lưu trữ file export.

---

## Công nghệ sử dụng

- **Golang** (Gin Framework)
- **GORM** (ORM cho MySQL)
- **MySQL** (Database)
- **AWS SQS** (Message Queue)\*\*\*\*
- **AWS S3** (Lưu trữ file)
- **Redis** (Cache trạng thái job)

---

## Cách hoạt động của luồng export

1. **Client gửi request export dữ liệu**
2. **Server tạo `jobId`**, gửi vào **AWS SQS Queue**, sau đó trả về cho client `jobId` và `status`
3. **Batch job worker** sẽ lấy dữ liệu từ queue, xử lý export, lưu file vào **AWS S3** và update trạng thái job trong **Redis**
4. **Client gửi request liên tục** để kiểm tra trạng thái job (server sẽ check trạng thái job trong Redis)
5. Khi **job hoàn thành**, client gửi request tải file từ **AWS S3**

---

## Cài đặt & Chạy dự án

### 1️⃣ Cài đặt môi trường

Yêu cầu:

- **Golang** >= 1.19
- **Docker** & **Docker Compose**

### 2️⃣ Clone repository

```sh
git clone https://github.com/your-username/go-file-processor.git
cd go-file-processor
```

### 3️⃣ Cấu hình biến môi trường

### 4️⃣ Khởi chạy hệ thống

Chạy Docker để khởi tạo MySQL, Redis:

```sh
docker-compose up -d
```

Chạy server Golang:

```sh
go run main.go
```

---

## API Documentation

### 1️⃣ **Tạo export job**

**Endpoint:**

```http
POST /create/:type
```

**Response:**

```json
{
  "data": {
    "id": "80f7574b-301a-467e-9e3d-cfdf9e355d8c",
    "status": "pending",
    "file_url": "",
    "export_type": "account"
  }
}
```

### 2️⃣ **Check trạng thái job**

**Endpoint:**

```http
GET /get-status/:id
```

**Response (Khi đang xử lý):**

```json
{
  "status_job": "completed" 
}
```
- "completed" => hoàn tất
- "pending" => chưa xử lý
- "processing" => đang xử lỹ
### 3️⃣ **Download file**

**Endpoint:**

```http
GET /download/:id
```

- Khi test trên postman, nhấn "save response to file" để tải file
- Trên browser, file sẽ tự tải về

---
