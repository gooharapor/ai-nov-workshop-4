## เพิ่ม Feature ให้สามารถ Transfer points ระหว่าง users ได้ (ใช้ swagger)
- ใช้ API specs swagger ไฟล์ `workshop-4/specs/transfer.yml` เป็นตัวกำกับและบอกเพิ่ม feature โดยเพิ่มอีก 2 table โดยใช้ query ด้านล่างนี้อ้างอิง

1. transfers — เก็บคำสั่งโอน + idemKey (ใช้ค้นหา)
```sql
CREATE TABLE IF NOT EXISTS transfers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  from_user_id INTEGER NOT NULL,
  to_user_id INTEGER NOT NULL,
  amount INTEGER NOT NULL CHECK (amount > 0),
  status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
  note TEXT,
  idempotency_key TEXT NOT NULL UNIQUE,     -- ใช้เป็น id ใน GET /transfers/{id}
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  completed_at TEXT,
  fail_reason TEXT,
  FOREIGN KEY (from_user_id) REFERENCES users(id),
  FOREIGN KEY (to_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_transfers_from     ON transfers(from_user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_to       ON transfers(to_user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_created  ON transfers(created_at);
```

2. point_ledger — สมุดบัญชีแต้ม (append-only)
```sql
CREATE TABLE IF NOT EXISTS point_ledger (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  change INTEGER NOT NULL,                  -- +รับโอน / -โอนออก
  balance_after INTEGER NOT NULL,
  event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
  transfer_id INTEGER,                      -- อ้างถึง transfers.id (ไอดีภายใน)
  reference TEXT,
  metadata TEXT,                            -- JSON text
  created_at TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (transfer_id) REFERENCES transfers(id)
);

CREATE INDEX IF NOT EXISTS idx_ledger_user     ON point_ledger(user_id);
CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id);
CREATE INDEX IF NOT EXISTS idx_ledger_created  ON point_ledger(created_at);
```

* ขั้นตอนนี้ เราจะตรวจจาก swagger ของ project ให้ push swagger ของ project ตัวเองขึ้นมา (จะใช้ library generate หรือ AI generate ก็ได้) = ผลลัพธ์เป็น `result.yml` เป็น swagger format
** (optional) สามารถ challenge เพิ่มได้โดยการเปลี่ยน user เป็น Authentication