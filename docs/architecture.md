<p align="center">
  <kbd><a href="#한국어">한국어</a></kbd> | <kbd><a href="#english">English</a></kbd>
</p>

---

# 한국어

## 시스템 개요

ecs9s는 AWS ECS 클러스터와 관련 서비스를 터미널에서 관리하는 Go 기반 TUI 애플리케이션입니다. Bubbletea의 Elm 아키텍처(Model-Update-View)를 사용하며, AWS SDK v2로 비동기 API 호출을 처리합니다.

## 레이어별 컴포넌트

### Presentation Layer
| 컴포넌트 | 파일 | 역할 |
|----------|------|------|
| App Shell | `internal/app/app.go` | 루트 모델, 라우팅, 페이지 관리 |
| Pages (21개) | `internal/ui/pages/` | 리소스별 뷰 모델 |
| Components (8개) | `internal/ui/components/` | 재사용 UI (table, tabs, commandbar 등) |
| Styles | `internal/ui/styles/` | Lipgloss 테마별 스타일 |
| Themes | `internal/theme/` | Dark, Light, Blue 프리셋 |

### Data Layer
| 컴포넌트 | 파일 | 역할 |
|----------|------|------|
| ECS Client | `internal/aws/ecs.go` | 클러스터, 서비스, 태스크, 컨테이너, TaskDef |
| CloudWatch Client | `internal/aws/cloudwatch.go` | 로그, 메트릭, 알람 |
| ECR Client | `internal/aws/ecr.go` | 레포지토리, 이미지 |
| ELB Client | `internal/aws/elb.go` | 로드밸런서, 타겟그룹 |
| EC2 Client | `internal/aws/ec2.go` | VPC, 서브넷, 보안그룹, 인스턴스 |
| IAM Client | `internal/aws/iam.go` | 역할, 정책 |
| AutoScaling Client | `internal/aws/autoscaling.go` | 스케일링 타겟/정책 |
| SSM Client | `internal/aws/ssm.go` | Parameter Store |
| Secrets Client | `internal/aws/secrets.go` | Secrets Manager |
| Session | `internal/aws/session.go` | 프로파일/리전 관리 |

### Action Layer
| 컴포넌트 | 파일 | 역할 |
|----------|------|------|
| ECS Exec | `internal/action/exec.go` | 컨테이너 셸 접속 |
| Port Forward | `internal/action/portforward.go` | SSM 포트포워딩 |
| Scale/Deploy | `internal/action/scale.go`, `deploy.go` | 서비스 스케일/배포 |
| Rollback | `internal/action/rollback.go` | 안전한 롤백 |

## 아키텍처 다이어그램

```
┌──────────────────────────────────────────────────────────────┐
│                        main.go                                │
│                   (CLI flags, config)                         │
└──────────────────────┬───────────────────────────────────────┘
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                    App Shell (app.go)                         │
│  ┌─────────┐ ┌──────────┐ ┌───────────┐ ┌────────┐         │
│  │ TabBar  │ │StatusBar │ │CommandBar │ │  Help  │         │
│  └─────────┘ └──────────┘ └───────────┘ └────────┘         │
│  ┌──────────────────────────────────────────────────┐       │
│  │              Active Page (1 of 21)                │       │
│  │  Cluster│Service│Task│Container│TaskDef│Logs│...  │       │
│  └──────────────────────┬───────────────────────────┘       │
└─────────────────────────┼────────────────────────────────────┘
                          ▼
┌──────────────────────────────────────────────────────────────┐
│                    AWS SDK v2 Clients                         │
│  ECS │ CloudWatch │ ECR │ ELB │ EC2 │ IAM │ ASG │ SSM │ SM │
└──────────────────────────┬───────────────────────────────────┘
                           ▼
                    AWS Cloud APIs
```

## 데이터 플로우

```
User Input ▶ App.Update() ▶ Page.Update() ▶ tea.Cmd(AWS API) ▶ fooLoadedMsg ▶ Page.View()
```

## 주요 설계 결정

| 결정 | 이유 |
|------|------|
| Bubbletea (Elm arch) | 복잡한 상태 관리에 유리, lipgloss로 세련된 UI |
| 하이브리드 네비게이션 | Tab + Command mode — 초보자와 파워유저 모두 지원 |
| NavContext drill-down | 계층적 리소스 탐색을 자연스럽게 구현 |
| WithDecryption: false | SecureString 값이 메모리에 저장되지 않도록 보안 강화 |
| DescribeTaskDefinition for cost | 이름 기반 추측 대신 실제 리소스 값으로 정확한 비용 추정 |
| SelectedRow() for selection | Cursor() 인덱스는 필터/정렬 시 불일치 — 행 데이터로 원본 조회 |
| tea.ExecProcess for shell | cmd.Run() 대신 사용하여 TUI 상태를 안전하게 중단/복원 |
| Tab 전환 시 context 유지 | DrillDown context를 보존하여 관련 페이지 간 이동 시 데이터 유지 |

---

# English

## System Overview

ecs9s is a Go-based TUI application for managing AWS ECS clusters and related services from the terminal. It uses Bubbletea's Elm architecture (Model-Update-View) and handles asynchronous API calls via AWS SDK v2.

## Components by Layer

### Presentation Layer
| Component | File | Role |
|-----------|------|------|
| App Shell | `internal/app/app.go` | Root model, routing, page management |
| Pages (21) | `internal/ui/pages/` | Per-resource view models |
| Components (8) | `internal/ui/components/` | Reusable UI (table, tabs, commandbar, etc.) |
| Styles | `internal/ui/styles/` | Lipgloss theme-based styles |
| Themes | `internal/theme/` | Dark, Light, Blue presets |

### Data Layer
| Component | File | Role |
|-----------|------|------|
| ECS Client | `internal/aws/ecs.go` | Clusters, services, tasks, containers, task defs |
| CloudWatch Client | `internal/aws/cloudwatch.go` | Logs, metrics, alarms |
| ECR Client | `internal/aws/ecr.go` | Repositories, images |
| ELB Client | `internal/aws/elb.go` | Load balancers, target groups |
| EC2 Client | `internal/aws/ec2.go` | VPCs, subnets, security groups, instances |
| IAM Client | `internal/aws/iam.go` | Roles, policies |
| AutoScaling Client | `internal/aws/autoscaling.go` | Scaling targets/policies |
| SSM Client | `internal/aws/ssm.go` | Parameter Store |
| Secrets Client | `internal/aws/secrets.go` | Secrets Manager |
| Session | `internal/aws/session.go` | Profile/region management |

### Action Layer
| Component | File | Role |
|-----------|------|------|
| ECS Exec | `internal/action/exec.go` | Container shell access |
| Port Forward | `internal/action/portforward.go` | SSM port forwarding |
| Scale/Deploy | `internal/action/scale.go`, `deploy.go` | Service scaling/deployment |
| Rollback | `internal/action/rollback.go` | Safe rollback with verification |

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                        main.go                                │
│                   (CLI flags, config)                         │
└──────────────────────┬───────────────────────────────────────┘
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                    App Shell (app.go)                         │
│  ┌─────────┐ ┌──────────┐ ┌───────────┐ ┌────────┐         │
│  │ TabBar  │ │StatusBar │ │CommandBar │ │  Help  │         │
│  └─────────┘ └──────────┘ └───────────┘ └────────┘         │
│  ┌──────────────────────────────────────────────────┐       │
│  │              Active Page (1 of 21)                │       │
│  │  Cluster│Service│Task│Container│TaskDef│Logs│...  │       │
│  └──────────────────────┬───────────────────────────┘       │
└─────────────────────────┼────────────────────────────────────┘
                          ▼
┌──────────────────────────────────────────────────────────────┐
│                    AWS SDK v2 Clients                         │
│  ECS │ CloudWatch │ ECR │ ELB │ EC2 │ IAM │ ASG │ SSM │ SM │
└──────────────────────────┬───────────────────────────────────┘
                           ▼
                    AWS Cloud APIs
```

## Data Flow

```
User Input ▶ App.Update() ▶ Page.Update() ▶ tea.Cmd(AWS API) ▶ fooLoadedMsg ▶ Page.View()
```

## Key Design Decisions

| Decision | Why |
|----------|-----|
| Bubbletea (Elm arch) | Excellent for complex state management; lipgloss for polished UI |
| Hybrid navigation | Tab + Command mode — serves both beginners and power users |
| NavContext drill-down | Natural hierarchical resource exploration |
| WithDecryption: false | SecureString values never stored in memory |
| DescribeTaskDefinition for cost | Accurate cost from real resource values, not name guessing |
| SelectedRow() for selection | Cursor() index breaks with filter/sort — look up by row data instead |
| tea.ExecProcess for shell | Safely suspend/restore TUI when running interactive ECS Exec |
| Preserve context on tab switch | DrillDown context retained so related pages share data across tabs |
