# aws

10 AWS SDK v2 client wrappers. One file per service. All clients take `aws.Config` from session.go.

| File | Service | Key Methods |
|------|---------|-------------|
| ecs.go | ECS | ListClusters, ListServices, ListTasks, DescribeContainers, ListTaskDefinitions, mutations (scale, deploy, stop, rollback), CheckExecEnabled, EnableExecOnService |
| cloudwatch.go | CloudWatch + Logs | GetLogEvents, ListLogGroups, GetECSMetrics, ListAlarms |
| ecr.go | ECR | ListRepositories, ListImages |
| elb.go | ELBv2 | ListLoadBalancers, ListTargetGroups (with health) |
| ec2.go | EC2 | ListVPCs, ListSubnets, ListSecurityGroups, ListInstances |
| iam.go | IAM | ListRoles, GetRolePolicies |
| autoscaling.go | App Auto Scaling | ListScalableTargets, ListScalingPolicies |
| ssm.go | SSM Parameter Store | ListParameters, GetParameter |
| secrets.go | Secrets Manager | ListSecrets |
| session.go | AWS Config | NewSession, SwitchProfile, SwitchRegion |

Pattern: client struct wraps SDK client; methods return domain types (not SDK types). Errors wrapped with `fmt.Errorf("action: %w", err)`.

Security: SSM `WithDecryption` must always be `false`. SecureString values masked with `"****"`.
