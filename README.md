add 

make migrate-create NAME=add_sessions


 atlas migrate hash --dir "file://migrations"


 make migrate-up



 # find 

  find . -type f -name "*.go" | grep -v vendor | sort



  step 

   make migrate-create NAME=add_blog_excerpt
Creating migration: add_blog_excerpt
atlas migrate new add_blog_excerpt \
                --dir "file://migrations"
cure@Macbook server %  atlas migrate hash --dir "file://migrations"
cure@Macbook server %  make migrate-up
Applying migrations...
atlas migrate apply \
                --env local \
                --config file://atlas.hcl \
                --allow-dirty
Migrating to version 20260314081025 from 20260309150718 (1 migrations in total):

  -- migrating version 20260314081025
    -> ALTER TABLE blogs
         ADD COLUMN IF NOT EXISTS excerpt TEXT;
    -> CREATE INDEX IF NOT EXISTS idx_blogs_excerpt ON blogs USING gin(to_tsvector('english', coalesce(excerpt, '')));
  -- ok (472.120375ms)

  -------------------------
  -- 1.14376325s
  -- 1 migration
  -- 2 sql statements

  

  