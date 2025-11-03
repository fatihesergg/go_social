DROP TABLE IF EXISTS posts CASCADE;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'post_visibility') THEN
        DROP TYPE post_visibility;
    END IF;
END
$$;
