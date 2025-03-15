
-- INSERT

-- Thêm người dùng vào hệ thống
INSERT INTO users (name, email, role) VALUES
('Admin User', 'admin@example.com', 'admin'),
('Accountant A', 'accountant_a@example.com', 'accountant');

-- Thêm danh mục tài khoản kế toán
INSERT INTO accounts (code, name, type) VALUES
('111', 'Tiền mặt', 'asset'),
('112', 'Tiền gửi ngân hàng', 'asset'),
('131', 'Phải thu khách hàng', 'asset'),
('331', 'Phải trả nhà cung cấp', 'liability'),
('511', 'Doanh thu bán hàng', 'revenue'),
('642', 'Chi phí quản lý doanh nghiệp', 'expense');


-- Thêm Journal Voucher
INSERT INTO journal_vouchers (voucher_no, date, description, status, created_by, total_debit, total_credit) VALUES
('JV-20240313-001', '2025-03-13', 'Thanh toán hóa đơn khách hàng', 'approved', 2, 15000000, 15000000),
('JV-20240313-002', '2025-03-13', 'Mua vật tư văn phòng', 'approved', 2, 5000000, 5000000),
('JV-20240313-003', '2025-03-13', 'Trả lương nhân viên', 'approved', 2, 20000000, 20000000),
('JV-20240313-004', '2025-03-13', 'Khách hàng thanh toán đơn hàng', 'approved', 2, 8000000, 8000000),
('JV-20240313-005', '2025-03-13', 'Chi phí điện nước văn phòng', 'approved', 2, 3000000, 3000000),
('JV-20240313-006', '2025-03-13', 'Nhận tạm ứng từ công ty mẹ', 'approved', 2, 10000000, 10000000),
('JV-20240313-007', '2025-03-13', 'Xuất hóa đơn bán hàng', 'approved', 2, 12000000, 12000000),
('JV-20240313-008', '2025-03-13', 'Trả tiền thuê văn phòng', 'approved', 2, 15000000, 15000000),
('JV-20240313-009', '2025-03-13', 'Nhập hàng tồn kho', 'approved', 2, 25000000, 25000000),
('JV-20240313-010', '2025-03-13', 'Thanh toán cho nhà cung cấp', 'approved', 2, 18000000, 18000000),
('JV-20240313-011', '2025-03-13', 'Doanh thu từ dịch vụ tư vấn', 'approved', 2, 9000000, 9000000);


-- Thêm bút toán vào Journal Voucher
INSERT INTO journal_entries (journal_voucher_id, account_id, debit_amount, credit_amount, description) VALUES
(1, 112, 15000000, 0, 'Khách hàng thanh toán qua ngân hàng'),
(1, 131, 0, 15000000, 'Giảm công nợ khách hàng'),

-- JV-002: Mua vật tư văn phòng
(2, 642, 5000000, 0, 'Chi phí vật tư văn phòng'),
(2, 111, 0, 5000000, 'Giảm tiền mặt'),

-- JV-003: Trả lương nhân viên
(3, 642, 20000000, 0, 'Chi phí lương nhân viên'),
(3, 112, 0, 20000000, 'Giảm tiền gửi ngân hàng'),

-- JV-004: Khách hàng thanh toán đơn hàng
(4, 112, 8000000, 0, 'Khách hàng chuyển khoản thanh toán'),
(4, 131, 0, 8000000, 'Giảm công nợ khách hàng'),

-- JV-005: Chi phí điện nước
(5, 642, 3000000, 0, 'Chi phí điện nước văn phòng'),
(5, 111, 0, 3000000, 'Giảm tiền mặt'),

-- JV-006: Nhận tạm ứng từ công ty mẹ
(6, 111, 10000000, 0, 'Nhận tiền tạm ứng từ công ty mẹ'),
(6, 331, 0, 10000000, 'Ghi nhận nợ phải trả'),

-- JV-007: Xuất hóa đơn bán hàng
(7, 131, 12000000, 0, 'Ghi nhận công nợ khách hàng'),
(7, 511, 0, 12000000, 'Doanh thu bán hàng'),

-- JV-008: Trả tiền thuê văn phòng
(8, 642, 15000000, 0, 'Chi phí thuê văn phòng'),
(8, 112, 0, 15000000, 'Giảm tiền ngân hàng'),

-- JV-009: Nhập hàng tồn kho
(9, 153, 25000000, 0, 'Nhập kho hàng hóa'),
(9, 112, 0, 25000000, 'Thanh toán tiền hàng'),

-- JV-010: Thanh toán cho nhà cung cấp
(10, 331, 18000000, 0, 'Thanh toán công nợ cho nhà cung cấp'),
(10, 112, 0, 18000000, 'Giảm tiền ngân hàng'),

-- JV-011: Doanh thu từ dịch vụ tư vấn
(11, 112, 9000000, 0, 'Khách hàng thanh toán tiền dịch vụ tư vấn'),
(11, 511, 0, 9000000, 'Doanh thu dịch vụ tư vấn');
