CREATE TABLE public.prices (
     id             integer     NOT NULL,
     name           text        NULL,
     category       text        NULL,
     price          numeric     NULL,
     create_date    date        NULL,
     CONSTRAINT item_pk PRIMARY KEY (id)
);
