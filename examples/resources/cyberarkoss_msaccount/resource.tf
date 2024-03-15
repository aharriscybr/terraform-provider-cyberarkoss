resource "cyberarkoss_msaccount" "mskey" {
  name             = "user-ms"
  address          = "1.2.3.4"
  username         = "user-ms"
  platform         = "MS_TF"
  safe             = "TF_TEST_SAFE"
  secrettype       = "password"
  secret           = "SincerelySecure2#24!"
  sm_manage        = false
  sm_manage_reason = "No CPM Associated with Safe."
  ms_appid         = "ApplicationID"
  ms_appobjid      = "ApplicationObjectID"
  ms_keyid         = "KeyID"
  ms_adid          = "ADKeyID"
  ms_duration      = "300"
  ms_pop           = "yes"
  ms_keydesc       = "key descriptiong with spaces"
}