-- Built into a temporary .db at test time. Never commit a .db.
CREATE TABLE codes (code TEXT NOT NULL, globdays TEXT);
INSERT INTO codes VALUES ('36415', '000'), ('99213', '000'), ('27447', '090');
