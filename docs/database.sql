CREATE DATABASE accounting_system;
USE accounting_system;

-- Bảng Users (Người dùng trong hệ thống)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role ENUM('admin', 'accountant', 'viewer') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Accounts (Danh mục tài khoản kế toán)
CREATE TABLE accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,  -- Ví dụ: 111, 131, 511
    name VARCHAR(255) NOT NULL,  -- Ví dụ: Tiền mặt, Phải thu khách hàng, Doanh thu
    type ENUM('asset', 'liability', 'equity', 'revenue', 'expense') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Journal Vouchers (Chứng từ nhật ký)
CREATE TABLE journal_vouchers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    voucher_no VARCHAR(50) NOT NULL UNIQUE,
    date DATE NOT NULL,
    description TEXT,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    created_by INT NOT NULL,
    approved_by INT DEFAULT NULL,
    total_debit DECIMAL(18,2) NOT NULL DEFAULT 0,
    total_credit DECIMAL(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng Journal Entries (Bút toán kế toán)
CREATE TABLE journal_entries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    journal_voucher_id INT NOT NULL,
    account_id INT NOT NULL,
    debit_amount DECIMAL(18,2) DEFAULT 0,
    credit_amount DECIMAL(18,2) DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
