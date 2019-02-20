DROP TABLE IF EXISTS "emails";
CREATE TABLE emails (
		ID INTEGER NOT NULL PRIMARY KEY,
		msgID text,
		thdID text,
		timeRecieved text,
		proccessed BOOLEAN DEFAULT 0,
		error BOOLEAN DEFAULT 0
	);