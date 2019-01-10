DROP TABLE IF EXISTS "shifts";
CREATE TABLE shifts (
		ID INTEGER NOT NULL PRIMARY KEY,
		Summery text,
		description text,
		TimeZone text,
		EventDateStart text,
		EventDateEnd text,
		Processed BOOLEAN DEFAULT FALSE,
		proccessTime text,
		eventID text,
		msgID int NOT NULL,
		FOREIGN KEY (msgID) REFERENCES emails ("msgID")
	);

INSERT INTO "shifts" ("ID", "Summery", "description", "TimeZone", "EventDateStart", "EventDateEnd", "Processed", "proccessTime", "eventID", "msgID") VALUES ('1', 'a', 's', 't', 'a', 'a', '1', 'd', 'd', '***REMOVED***');
