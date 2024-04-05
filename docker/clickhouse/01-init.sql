CREATE TABLE IF NOT EXISTS measurements
(
    temperature Float32,
    humidity    Float32,
    timestamp   DateTime('Europe/Moscow'),
)
ENGINE = MergeTree
ORDER BY timestamp;