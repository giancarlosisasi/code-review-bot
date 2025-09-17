# Code Review Bot 🤖 (WIP)

An intelligent Slack bot that automatically assigns code reviewers to GitLab merge requests using smart team balancing and workload management algorithms.

> **📚 Hey there!**  
> This is my **Golang learning project** - I'm teaching myself Go by building something useful! I wrote all the code myself (no AI code generation), just used Claude to help me understand Go concepts and find the right documentation when I got stuck. It's been a great way to learn the language hands-on! 🚀

## 🚀 Features

- **Smart Reviewer Assignment**: Automatically selects reviewers based on seniority levels and current workload
- **Workload Balancing**: Ensures fair distribution of review tasks (max 3 MRs per person)
- **Team Guild Support**: Organizes team members by guilds (frontend, backend, devops, qa)
- **Slack Integration**: Simple `/code-review <mr-url>` slash command
- **GitLab Integration**: Seamlessly works with GitLab merge requests
- **Intelligent Escalation**: Automatically escalates when team capacity is reached

## 🏗️ **Architecture**

The bot consists of two main services running concurrently:

1. **Web Server**: Handles GitLab webhooks and provides REST API endpoints
2. **Slack Bot**: Processes slash commands and manages team communication

### Core Components

- **Assignment Algorithm**: Prioritizes senior developers while balancing workload
- **Team Management**: JSON-based team configuration with seniority weights
- **In-Memory Storage**: Fast, temporary storage for MR assignments and workload tracking
- **GitLab API Client**: Fetches MR details and assigns reviewers

## 📋 Prerequisites

- Go 1.24.3 or later
- GitLab API access token
- Slack app with appropriate permissions

## 🛠️ Setup

### 1. Clone the Repository

```bash
git clone https://github.com/giancarlosisasi/code-review-bot.git
cd code-review-bot
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080
APP_ENV=development
DB_URL=postgres://postgres:postgres@localhost:5433/thedbname

# GitLab Configuration
GITLAB_API_TOKEN=your_gitlab_token_here
GITLAB_ORG_SLUG=your_organization_slug

# Slack Configuration
SLACK_SOCKET_MODE_TOKEN=xapp-your_socket_mode_token
SLACK_BOT_OAUTH_TOKEN=xoxb-your_bot_oauth_token
SLACK_SIGNING_SECRET=your_signing_secret
```

> **Note**: The `DB_URL` is currently not used as the application runs with in-memory storage. PostgreSQL integration is planned for future releases.

### 4. Team Configuration

Create a `users.json` file with your team members:

```json
[
  {
    "id": "user@example.com",
    "name": "John Doe",
    "email": "user@example.com",
    "gitlab_member_id": "12345",
    "slack_member_id": "U1234567890",
    "seniority_weight": 55,
    "team_guild": "frontend",
    "role": "admin"
  }
]
```

**Seniority Weights:**
- Senior: 55
- Semi-Senior: 30-50
- Junior: 15

### 5. Run the Application

```bash
go run cmd/server/main.go
```

## 🎯 Usage

### Slack Command

Use the `/code-review` slash command in any Slack channel:

```
/code-review https://gitlab.com/your-org/your-project/-/merge_requests/123
```

The bot will:
1. Parse the GitLab MR URL
2. Fetch MR details from GitLab API
3. Apply the assignment algorithm
4. Assign reviewers to the MR
5. Notify the team in Slack

### Assignment Algorithm

The bot uses a sophisticated algorithm to select reviewers:

1. **Excludes**: MR author and already assigned reviewers
2. **Prioritizes by**: Seniority level (Senior → Semi-Senior → Junior)
3. **Balances workload**: Prefers team members with fewer active MRs
4. **Escalates when needed**: Increases limits if no reviewers are available
5. **Assigns**: Up to 2 reviewers per MR (configurable)

## 🔧 Configuration

### Team Guilds

The bot supports multiple team guilds:
- `frontend`
- `backend` 
- `devops`
- `qa`

### Workload Limits

- **Default**: 3 MRs per person maximum
- **Escalation**: Automatically increases limits when team is at capacity
- **Admin notification**: Alerts admin when manual intervention needed

## 🚧 Roadmap

- [x] GitLab integration
- [x] Smart reviewer assignment algorithm
- [x] Workload balancing
- [x] Slack slash commands
- [ ] Add GitHub support
- [ ] Database persistence (PostgreSQL)
- [ ] Containerize with Docker
- [ ] Webhook automation
- [ ] Analytics and reporting
- [ ] Configuration UI

## 🏃‍♂️ Development

### Project Structure

```
code-review-bot/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── database/            # Database layer (in-memory + PostgreSQL)
├── gitlab_client/       # GitLab API integration
├── models/              # Data models
├── repository/          # Data access layer
├── server/              # HTTP server and handlers
├── slackbot/            # Slack bot implementation
├── utils/               # Utility functions
└── users.json           # Team configuration
```

### Key Dependencies

- **Gin**: HTTP web framework
- **Slack Go SDK**: Slack API integration
- **GitLab Go SDK**: GitLab API client
- **Viper**: Configuration management
- **Zerolog**: Structured logging
- **PostgreSQL**: Database (future)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For questions or issues, please:
1. Check the existing issues
2. Create a new issue with detailed information
3. Contact the development team

---

**Built with ❤️ for better code review processes**