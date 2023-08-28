-- to quickly generate xids use : https://play.golang.org/p/xSMJJ9G9Lt_-

-- Default base partition
INSERT INTO tenants (id, tenant_id, partition_id, name, description) VALUES('9bsv0s0hijjg02qks6jg', '9bsv0s0hijjg02qks6jg', '9bsv0s0hijjg02qks6kg', 'Stawi Development', 'Default base tenant for testing and building stawi');
INSERT INTO partitions (id, tenant_id, partition_id, name, description, properties)
    VALUES('9bsv0s0hijjg02qks6kg', '9bsv0s0hijjg02qks6jg', '9bsv0s0hijjg02qks6kg',
           'Stawi', 'Default stawi partition to serve the masses',
           '{"scope": "openid offline offline_access profile contact", "audience": ["service_chat_engine", "service_profile", "service_stawi_api", "service_files"], "logo_uri": "https://static.stawi.io/logo.png", "redirect_uris": ["https://app.stawi.io"]}');

