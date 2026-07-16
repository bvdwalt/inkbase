CREATE TABLE pages (
    id           TEXT PRIMARY KEY,
    parent_id    TEXT REFERENCES pages(id) ON DELETE SET NULL,
    slug         TEXT NOT NULL UNIQUE,
    title        TEXT NOT NULL,
    content_json TEXT NOT NULL DEFAULT '',
    content_text TEXT NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pages_parent_id ON pages(parent_id);

CREATE TABLE revisions (
    id           TEXT PRIMARY KEY,
    page_id      TEXT NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    content_json TEXT NOT NULL,
    content_text TEXT NOT NULL,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_revisions_page_id ON revisions(page_id);

-- A standalone (non external-content) FTS5 table: page ids are TEXT (UUIDs),
-- but FTS5's content_rowid must be an integer, so instead of syncing via
-- external-content mode we store page_id as a plain unindexed column and
-- sync it with ordinary INSERT/UPDATE/DELETE, which FTS5 supports directly.
CREATE VIRTUAL TABLE pages_fts USING fts5(
    page_id UNINDEXED,
    title,
    content_text
);

CREATE TRIGGER pages_ai AFTER INSERT ON pages BEGIN
    INSERT INTO pages_fts(page_id, title, content_text) VALUES (new.id, new.title, new.content_text);
END;

CREATE TRIGGER pages_ad AFTER DELETE ON pages BEGIN
    DELETE FROM pages_fts WHERE page_id = old.id;
END;

CREATE TRIGGER pages_au AFTER UPDATE ON pages BEGIN
    DELETE FROM pages_fts WHERE page_id = old.id;
    INSERT INTO pages_fts(page_id, title, content_text) VALUES (new.id, new.title, new.content_text);
END;
