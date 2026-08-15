-- Enforce unique role names.
ALTER TABLE role
    ADD CONSTRAINT uq_role_role_name UNIQUE (role_name);

-- Enforce unique email addresses.
ALTER TABLE "user"
    ADD CONSTRAINT uq_user_email UNIQUE (email);

-- Phone is optional, but every non-NULL phone number must be unique.
ALTER TABLE "user"
    ADD CONSTRAINT uq_user_phone UNIQUE (phone);

-- Support OTP lookup by user.
CREATE INDEX idx_otp_user_id
    ON otp(user_id);