variable "region" {
  default = "us-west-2"
}

variable "app_name" {
  default = "equilibria"
}

variable "mysql_version" {
  default = "8.0"
}

# Environment variables
variable "mysql_root_password" {}

variable "mysql_root_user" {}

variable "database_user" {}

variable "database_password" {}

variable "database_name" {}

variable "log_level" {}

variable "open_ai_api_key" {}

variable "twilio_sid" {}

variable "twilio_auth_token" {}

variable "twilio_phone_number" {}

variable "twilio_status_callback_url" {}

variable "twilio_verify_service_sid" {}
