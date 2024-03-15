resource "cyberarkoss_dbaccount" "pgdb" {
  name             = "user-db"
  address          = "1.2.3.4"
  username         = "user-db"
  platform         = "PROD_PostgreSQL"
  safe             = "TF_TEST_SAFE"
  secrettype       = "password"
  secret           = "SincerelySecure2#24!"
  sm_manage        = false
  sm_manage_reason = "No CPM Associated with Safe."
  db_port          = "8432"
  db_dsn           = "dsn"
  dbname           = "dbo.services"
}