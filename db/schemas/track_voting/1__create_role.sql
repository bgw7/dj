DO 
$$ 
BEGIN IF NOT EXISTS(
    SELECT
    FROM
        pg_roles
    WHERE
        rolname = 'api_user'
) THEN CREATE USER api_user LOGIN ENCRYPTED PASSWORD 'api_password';

END IF;
END;
$$;
