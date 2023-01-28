module "robo_trade1_lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "4.8.0"
 
  function_name = "robo_trade1"
  description   = "robo trade 1"
  handler       = "robo_trade1"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 10
  architectures = ["x86_64"]

  source_path = "../robo_trade1"

  environment_variables = {
    INTERVAL = ""
  }
}