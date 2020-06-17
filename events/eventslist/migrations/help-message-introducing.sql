/* Create in-memory temp table for variables */
BEGIN;

PRAGMA temp_store = 2;
CREATE TEMP TABLE vars
(
	name  TEXT PRIMARY KEY,
	value TEXT
);

/* Declaring a variables */
INSERT INTO vars (name)
VALUES ('eventID');
INSERT INTO vars (name)
VALUES ('scenarioID');
INSERT INTO vars (name)
VALUES ('regexID');

/* Assigning a variable (pick the right storage class) */
UPDATE vars
SET value = (select id from events where alias = 'eventslist')
WHERE Name = 'eventID';

/* Creation of new event question */

insert into main.scenarios (name, event_id)
values ('Scenario help', (SELECT value FROM vars WHERE Name = 'eventID' LIMIT 1));

UPDATE vars
SET value = (select last_insert_rowid())
WHERE Name = 'scenarioID';

insert into main.questions_regex (regex, regex_group)
values ('(?i)^(help)$', '');

UPDATE vars
SET value = (select last_insert_rowid())
WHERE Name = 'regexID';
insert into main.questions (question, answer, scenario_id, regex_id)
values ('Help',
		'If you want to see the list of my functions, please try ask me the following question `events list`. This will printout all possible phrases what currently I can understand.',
		(SELECT value FROM vars WHERE Name = 'scenarioID' LIMIT 1),
		(SELECT value FROM vars WHERE Name = 'regexID' LIMIT 1));

DROP TABLE vars;
END;