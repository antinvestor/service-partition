-- to quickly generate xids use : https://play.golang.org/p/xSMJJ9G9Lt_-

-- Default base partition
INSERT INTO tenants (id, tenant_id, partition_id, name, description) VALUES('9bsv0s0hijjg02qks6dg', '9bsv0s0hijjg02qks6dg', '9bsv0s0hijjg02qks6i0', 'Stawi Development', 'Default base tenant for testing and building stawi');
INSERT INTO partitions (id, tenant_id, partition_id, name, description, properties)
    VALUES('9bsv0s0hijjg02qks6i0', '9bsv0s0hijjg02qks6dg', '9bsv0s0hijjg02qks6i0',
           'Stawi Development', 'Default stawi development partition in stawi development',
           '{"scope": "openid offline_access profile contact", "audience": ["service_chat_engine", "service_profile", "service_stawi_api", "service_files"], "logo_uri": "https://static.stawi.io/logo.png", "redirect_uris": ["https://app-dev.stawi.io", "http://localhost:40455", "http://127.0.0.1:40455"]}');

