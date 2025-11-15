-- +goose Up
-- +goose StatementBegin
-- Тип для OAuth-провайдеров
CREATE TYPE public.oauthproviderenum AS ENUM ('GOOGLE', 'YANDEX');

-- Таблица пользователей
CREATE TABLE public.users (
                              id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                              name            VARCHAR NOT NULL,
                              is_creator      BOOLEAN,
                              socket_id       VARCHAR,
                              created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                              updated_at      TIMESTAMP WITH TIME ZONE,
                              is_watcher      BOOLEAN DEFAULT FALSE,
                              on_session      BOOLEAN DEFAULT TRUE,
                              email           VARCHAR,
                              hashed_password VARCHAR,
                              is_active       BOOLEAN DEFAULT TRUE,
                              is_verified     BOOLEAN DEFAULT FALSE,
                              oauth_provider  oauthproviderenum,
                              oauth_id        VARCHAR,
                              avatar_url      VARCHAR,
                              is_guest        BOOLEAN
);
ALTER TABLE public.users OWNER TO agile_poker_user;

-- Таблица сессий
CREATE TABLE public.sessions (
                                 id             UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                 name           VARCHAR NOT NULL,
                                 deck_type      VARCHAR NOT NULL,
                                 cards_revealed BOOLEAN,
                                 creator_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 creator_name   VARCHAR NOT NULL,
                                 created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                 updated_at     TIMESTAMP WITH TIME ZONE,
                                 allow_emoji    BOOLEAN DEFAULT TRUE,
                                 auto_reveal    BOOLEAN DEFAULT TRUE,
                                 created_via    VARCHAR NOT NULL
);
ALTER TABLE public.sessions OWNER TO agile_poker_user;

-- Уникальный индекс по email
CREATE UNIQUE INDEX ix_users_email ON public.users (email);

-- Таблица голосов
CREATE TABLE public.votes (
                              id         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                              session_id UUID NOT NULL REFERENCES public.sessions ON DELETE CASCADE,
                              user_id    UUID NOT NULL REFERENCES public.users ON DELETE CASCADE,
                              value      VARCHAR NOT NULL,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                              updated_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE public.votes OWNER TO agile_poker_user;

-- Таблица реакций
CREATE TABLE public.reactions (
                                  id           UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                                  session_id   UUID NOT NULL REFERENCES public.sessions ON DELETE CASCADE,
                                  from_user_id UUID NOT NULL REFERENCES public.users ON DELETE CASCADE,
                                  to_user_id   UUID NOT NULL REFERENCES public.users ON DELETE CASCADE,
                                  emoji        VARCHAR NOT NULL,
                                  created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
ALTER TABLE public.reactions OWNER TO agile_poker_user;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.reactions;
DROP TABLE IF EXISTS public.votes;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.sessions;
DROP TYPE IF EXISTS public.oauthproviderenum;
-- +goose StatementEnd
