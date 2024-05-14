CREATE TABLE IF NOT EXISTS users(
    username TEXT PRIMARY KEY,
    role INTEGER,
    pass_hash BLOB
);
INSERT OR REPLACE INTO users(username, role, pass_hash) VALUES ('user1', 0, X'243261243130247878712f76375459727330464f36512e686c4b4f762e52656b504635305359415963354d6d5743764131504d79475235496b726747');
INSERT OR REPLACE INTO users(username, role, pass_hash) VALUES ('user2', 0, X'24326124313024686c634169694f6e4e647a7946314531325737633765486b6768575a793534646e36434663536e58616e4d53557175347453357947');
INSERT OR REPLACE INTO users(username, role, pass_hash) VALUES ('admin', 1, X'24326124313024766d2f486e5a424d684d514f4763674167776a64474f446a68395764536a395943306f6f307074682f627363675074754c334e504b')