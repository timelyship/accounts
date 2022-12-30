This is a sample project written in Golang to demonstrate my coding ability in Go! This project represents a microserivce written for user account management features such as LogIn, SignUp etc. It also uses a SQS queue to demonstrate an integration with the email sending service(written in nodejs and SES). I had lot's of fun writing this service and deploying it in a serverless fashion in AWS lambda. Hope you will find some useful patterns here to be used in your application too! Good luck exploring my codebase and please leave a comment if you have any suggession!
# accounts
User management Service

- Config.yaml represents all the settings required to run the project.
- Value present in the Config.yaml file solely purpose of development. If you pass an environment variable it will ignore config.yaml.
- In production, we need to pass all these keys using environment variable.
