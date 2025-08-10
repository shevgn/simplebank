CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    username text NOT NULL,
    refresh_token text NOT NULL,
    user_agent text NOT NULL,
    client_ip text NOT NULL,
    is_blocked boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    expires_at timestamp with time zone NOT NULL
);

ALTER TABLE sessions ADD FOREIGN KEY (username) REFERENCES users (username);
