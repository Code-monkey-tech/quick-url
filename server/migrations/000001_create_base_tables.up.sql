create table if not exists public.url
(
    id          serial primary key,
    short_url   text unique not null,
    long_url    text        not null,
    insert_date timestamp default now(),
    expire_date timestamp   null
);

create index url_short_url_idx on public.url (short_url);