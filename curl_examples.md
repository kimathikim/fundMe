# DealFlow API Endpoints - curl Examples

## Add a startup to dealflow
```bash
curl -X POST http://localhost:8080/api/v1/dealflow \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "UserID": "60d21b4667d0d8992e610c85",
    "StartupName": "AI Solutions",
    "Industry": "AI/ML",
    "FundingStage": "Series A",
    "Location": "San Francisco",
    "FundRequired": 500000,
    "MatchScore": 85.5,
    "Tags": ["AI", "Machine Learning"]
  }'
```

## Get dealflow by ID
```bash
curl -X GET http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## List all dealflow entries
```bash
curl -X GET http://localhost:8080/api/v1/dealflow \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Update dealflow
```bash
curl -X PUT http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "MatchScore": 92.5,
    "Tags": ["AI", "Machine Learning", "NLP"]
  }'
```

## Delete dealflow
```bash
curl -X DELETE http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Invest in startup
```bash
curl -X POST http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/invest \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "investmentAmount": 250000
  }'
```

## Add meeting
```bash
curl -X POST http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/meetings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Initial Discussion",
    "date": "2023-06-15T14:00:00Z",
    "notes": "Discuss funding requirements and business model"
  }'
```

## Add document
```bash
curl -X POST http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/documents \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Pitch Deck",
    "url": "https://example.com/pitchdeck.pdf",
    "type": "presentation"
  }'
```

## Add task
```bash
curl -X POST http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Review financial statements",
    "description": "Analyze last 3 years of financial data",
    "dueDate": "2023-06-20T00:00:00Z",
    "completed": false
  }'
```

## Update task status
```bash
curl -X PATCH http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/tasks/60d21b4667d0d8992e610c86 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

## Update deal stage
```bash
curl -X PATCH http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/stage \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stage": "dueDiligence"
  }'
```

## Update deal status
```bash
curl -X PATCH http://localhost:8080/api/v1/dealflow/60d21b4667d0d8992e610c85/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "active"
  }'
```

# Additional API Endpoints - curl Examples

## Founder Grant Routes

### Get Available Grants
```bash
curl -X GET http://localhost:8080/api/v1/founder/grants \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "tech",
    "region": "US"
  }'
```

### Submit Grant Application
```bash
curl -X POST http://localhost:8080/api/v1/founder/grants/apply \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: multipart/form-data" \
  -F "grantId=60d21b4667d0d8992e610c85" \
  -F "startupName=EcoTech Solutions" \
  -F "contactEmail=founder@ecotech.com" \
  -F "contactPhone=+1234567890" \
  -F "description=Sustainable technology solutions for renewable energy" \
  -F "website=https://ecotech.com" \
  -F "teamSize=12" \
  -F "previousFunding=250000" \
  -F "pitchDeck=@/path/to/pitchdeck.pdf"
```

## Founder Notification Routes

### Get All Notifications
```bash
curl -X GET http://localhost:8080/api/v1/founder/notifications \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Notification
```bash
curl -X PUT http://localhost:8080/api/v1/founder/notification/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "read": true
  }'
```

### Delete Notification
```bash
curl -X DELETE http://localhost:8080/api/v1/founder/notification/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Investor Meeting Routes

### Add Meeting (Investor-specific endpoint)
```bash
curl -X POST http://localhost:8080/api/v1/investor/investor/60d21b4667d0d8992e610c85/meeting \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Investment Discussion",
    "date": "2023-07-15T10:00:00Z",
    "location": "Virtual",
    "description": "Initial discussion about potential investment",
    "attendees": ["60d21b4667d0d8992e610c86", "60d21b4667d0d8992e610c87"]
  }'
```

### Get Meetings (Investor-specific endpoint)
```bash
curl -X GET http://localhost:8080/api/v1/investor/investor/60d21b4667d0d8992e610c85/meetings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Investor Dashboard Routes

### Get Investor Dashboard
```bash
curl -X GET http://localhost:8080/api/v1/investor/dashboard \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Portfolio Performance
```bash
curl -X GET http://localhost:8080/api/v1/investor/portfolio/performance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Investor Startup Routes

### Get Startup Details
```bash
curl -X GET http://localhost:8080/api/v1/investor/startups \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Founder Profiles
```bash
curl -X GET http://localhost:8080/api/v1/investor/founderProfile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Match Routes

### Get Match Data
```bash
curl -X GET http://localhost:8080/api/v1/match/data/60d21b4667d0d8992e610c85 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Authentication Routes

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Current User
```bash
curl -X GET http://localhost:8080/api/v1/get/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## User Routes

### Register
```bash
curl -X POST http://localhost:8080/api/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepassword",
    "name": "John Doe",
    "role": "founder"
  }'
```
