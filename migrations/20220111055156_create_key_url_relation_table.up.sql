CREATE TABLE IF NOT EXISTS IndexRelation(
    keyId  bigint UNSIGNED NOT NULL,
    urlId  bigint UNSIGNED NOT NULL,
    FOREIGN KEY(keyId) REFERENCES `Key`(id),
    FOREIGN KEY(urlId) REFERENCES `Url`(id)
);
