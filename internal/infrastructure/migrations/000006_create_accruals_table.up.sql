CREATE TABLE accruals (
        order_number TEXT PRIMARY KEY,
        amount  float default 0.0,
        state   varchar(100) NOT NULL,
        created  TIMESTAMP WITH TIME ZONE NOT NULL
);