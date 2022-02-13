CREATE TABLE tags (  
    id uuid NOT NULL,   
    name varchar(100) NOT NULL,
    channel_id bigint NOT NULL,   
    CONSTRAINT tags_pk PRIMARY KEY (id)
);

CREATE TABLE subscriptions (  
    id bigint NOT NULL,   
    name varchar(150) NOT NULL,
    tag_id uuid REFERENCES tags(id),
    CONSTRAINT subscriptions_pk PRIMARY KEY (id)
);