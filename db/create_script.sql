CREATE TABLE "user"(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    "password" text NOT NULL,
    created_at timestamp(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_user_id PRIMARY KEY(id)
);

CREATE TABLE user_auth(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    expires_at timestamp(0) NOT NULL,
    CONSTRAINT pk_user_auth_id PRIMARY KEY(id),
    CONSTRAINT fk_user_auth_user_id FOREIGN KEY(user_id) REFERENCES "user"(id)
);

CREATE TABLE task(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_name text NOT NULL,
    task_icon text NOT NULL, 
    task_desc text NOT NULL,
    deadline timestamp(0) WITH TIME ZONE,
    starred boolean NOT NULL,
    exec_status text NOT NUll,
    created_at timestamp(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by uuid NOT NULL,
    CONSTRAINT pk_task_id PRIMARY KEY(id),
    CONSTRAINT fk_task_created_by FOREIGN KEY(created_by) REFERENCES "user"(id)
);

CREATE TABLE task_history(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    exec_rating int,
    exec_comment text,
    task_id uuid NOT NULL,
    CONSTRAINT pk_task_history_id PRIMARY KEY(id),
    CONSTRAINT fk_task_history_task_id FOREIGN KEY(task_id) REFERENCES task(id)
);