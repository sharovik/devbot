create table events
(
	id integer
		constraint events_pk
			primary key autoincrement,
	alias varchar not null,
	installed_version varchar
);

create unique index events_name_uindex
	on events (alias);

create table questions_regex
(
	id integer not null
		constraint question_regex_pk
			primary key autoincrement,
	regex varchar not null,
	regex_group varchar
);

create unique index question_regex_regex_uindex
	on questions_regex (regex);

create table scenarios
(
	id integer
		constraint scenarios_pk
			primary key autoincrement,
	name varchar not null,
	event_id integer not null
		references events
			on delete set null
);

create table questions
(
	id integer
		constraint questions_pk
			primary key autoincrement,
	question varchar not null,
	answer varchar not null,
	scenario_id int not null
		references scenarios
			on delete cascade,
	regex_id integer
		references questions_regex
			on delete set null
);

create index questions_question_index
	on questions (question);

create unique index scenarios_name_uindex
	on scenarios (name);

create table migration
(
	id integer
		constraint migration_pk
			primary key autoincrement,
	version varchar not null
);

create unique index migration_version_uindex
	on migration (version);


insert into main.events (id, alias, installed_version) values (2, 'dictionary', null);
insert into main.events (id, alias, installed_version) values (4, 'hello', null);

insert into main.questions (id, question, answer, scenario_id, regex_id) values (2, 'New answer', 'Well, well, well. Ok, I''ll put this data in my mind.', 2, 2);
insert into main.questions (id, question, answer, scenario_id, regex_id) values (4, 'Hello', 'Hey', 4, 0);
insert into main.questions (id, question, answer, scenario_id, regex_id) values (5, 'Hey', 'Yo', 5, 0);
insert into main.questions (id, question, answer, scenario_id, regex_id) values (6, 'Here?', 'Yes, what do you need?', 6, 4);
insert into main.questions (id, question, answer, scenario_id, regex_id) values (7, 'Say hello to John', 'Hey %s', 7, 5);
insert into main.questions (id, question, answer, scenario_id, regex_id) values (10, 'Wellcome John', 'Wellcome %s', 10, 8);

insert into main.questions_regex (id, regex, regex_group) values (2, '(?i)(New answer)', '');
insert into main.questions_regex (id, regex, regex_group) values (4, '(?i)(here\?)', '');
insert into main.questions_regex (id, regex, regex_group) values (5, '(?i)Say hello to (?P<name>.+)', 'name');
insert into main.questions_regex (id, regex, regex_group) values (8, '(?i)Wellcome (?P<name>.+)', 'name');

insert into main.scenarios (id, name, event_id) values (2, 'Scenario #2', 2);
insert into main.scenarios (id, name, event_id) values (4, 'Scenario #4', 4);
insert into main.scenarios (id, name, event_id) values (5, 'Scenario #5', 4);
insert into main.scenarios (id, name, event_id) values (6, 'Scenario #6', 4);
insert into main.scenarios (id, name, event_id) values (7, 'Scenario #7', 4);
insert into main.scenarios (id, name, event_id) values (10, 'Scenario #8', 4);