CREATE TYPE user_role AS ENUM ('ADMIN','PSYCHOLOGIST','STUDENT','PATIENT');
CREATE TYPE request_status AS ENUM ('OPEN','IN_TRIAGE','ASSIGNED','CLOSED');

CREATE TABLE users (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	role user_role NOT NULL,
	crp TEXT,
	institution TEXT,
	is_approved BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE service_requests (
	id UUID PRIMARY KEY,
	patient_id UUID REFERENCES users(id),
	is_anonymous BOOLEAN DEFAULT false,
	description TEXT NOT NULL,
	area TEXT,
	urgency INT DEFAULT 1, -- 1 a 5
	status request_status DEFAULT 'OPEN',
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE assignments (
	id UUID PRIMARY KEY,
	request_id UUID REFERENCES service_requests(id) ON DELETE CASCADE,
	assignee_id UUID REFERENCES users(id),
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE sessions (
	id UUID PRIMARY KEY,
	request_id UUID REFERENCES service_requests(id) ON DELETE SET NULL,
	professional_id UUID REFERENCES users(id),
	student_id UUID REFERENCES users(id),
	date TIMESTAMPTZ NOT NULL,
	duration_minutes INT NOT NULL,
	type TEXT NOT NULL, -- online, chat, etc
	notes TEXT, -- anotações sigilosas
	supervisor_crp TEXT,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE feedbacks (
	id UUID PRIMARY KEY,
	request_id UUID REFERENCES service_requests(id) ON DELETE CASCADE,
	rating_acolhimento INT CHECK (rating_acolhimento BETWEEN 1 AND 5),
	rating_tecnica INT CHECK (rating_tecnica BETWEEN 1 AND 5),
	rating_resultado INT CHECK (rating_resultado BETWEEN 1 AND 5),
	comment TEXT,
	flag_issue BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE chat_rooms (
	id UUID PRIMARY KEY,
	request_id UUID REFERENCES service_requests(id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE chat_messages (
	id UUID PRIMARY KEY,
	room_id UUID REFERENCES chat_rooms(id) ON DELETE CASCADE,
	sender TEXT NOT NULL, -- "patient" ou "professional"
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now()
);
