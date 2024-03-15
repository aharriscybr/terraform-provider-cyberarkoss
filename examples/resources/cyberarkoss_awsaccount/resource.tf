resource "cyberarkoss_awsaccount" "awskey" {
  name              = "user-aws"
  username          = "user-aws"
  platform          = "AWS_TF"
  safe              = "TF_TEST_SAFE"
  secrettype        = "key"
  secret            = "secret_key"
  sm_manage         = false
  sm_manage_reason  = "No CPM Associated with Safe."
  aws_kid           = "9876543210"
  aws_accountid     = "0123456789"
  aws_alias         = "aws_alias"
  aws_accountregion = "us-east-2"
}