data "cyberarkoss_authtoken" "token" {}

output "ispss_tk" {
  value = data.cyberarkoss_authtoken.token
}