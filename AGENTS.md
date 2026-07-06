## Imported Claude Cowork project instructions

Claude AI Master Prompt – BitwiseLearn V2 SaaS EdTech Platform with Subscription Business Model
Act as a Chief Technology Officer (CTO), Principal Software Architect, Azure Cloud Architect, Product Architect, SaaS Architect, FinOps Consultant, Revenue Model Strategist, Security Architect, and DevOps Architect.
Design a complete production-ready SaaS EdTech ecosystem named BitwiseLearn V2 that serves Students, Institutions, and Vendors through a subscription-based business model. The platform must be optimized for Azure Cloud, use Hybrid Microservices Architecture, support online coding assessments and LMS, and automatically generate pricing plans based on cloud operational costs with a minimum revenue target of 3X infrastructure cost.
Business Vision
BitwiseLearn V2 is a B2B + B2C SaaS EdTech platform providing:

* LMS Platform
* Coding Platform
* Placement Training
* Assessment Platform
* AI Learning Assistant
* Resume Analysis
* Coding Contests
* Institution Analytics
* Training Vendor Management
* Certification Platform
The system must support:
Customer Types

1. Students
2. Colleges / Universities
3. Training Institutes
4. Corporate Training Organizations
5. Vendors / Freelance Trainers
Revenue Model
The platform must automatically calculate subscription plans based on:

```text
Revenue Target = 3 × Total Cloud Operational Cost

```

The architecture should include a Financial Planning Module that:

* Tracks Azure Monthly Cost
* Tracks Active Users
* Tracks Storage Consumption
* Tracks Compiler Usage
* Tracks Assessment Usage
* Tracks AI Usage
* Calculates Cost Per User
* Recommends Subscription Pricing
Subscription Architecture
Design complete subscription management microservices.
Student Subscription Plans
Free Plan
Features:

* Limited LMS Access
* 5 Coding Problems Per Day
* 1 Mock Test Per Month
* Community Access
Restrictions:

* No Certificates
* No AI Tutor
* No Placement Analytics
Student Pro Plan
Target:
Individual Learners
Features:

* Unlimited LMS
* Unlimited Coding Practice
* Weekly Assessments
* Resume Builder
* AI Tutor
* Certificates
* Placement Dashboard
Suggested Pricing Formula:

```text
Monthly Cost Per Student × 3

```

Expected Price Range:
₹199 – ₹499/month
Student Premium Plan
Features:

* Everything in Pro
* AI Mock Interviews
* Company Wise Preparation
* Advanced Analytics
* Personalized Learning Paths
Expected Price Range:
₹499 – ₹999/month
Institution Subscription Plans
Starter Institution Plan
Target:
1 College
Limits:

* 100 Students
* 10 Trainers
Features:

* LMS
* Assessments
* Reports
* Coding Platform
Expected Pricing:
₹25,000–₹40,000/year
Growth Institution Plan
Target:
Up to 500 Students
Features:

* AI Analytics
* Placement Dashboard
* Coding Contests
* API Access
Expected Pricing:
₹75,000–₹2 Lakhs/year
Enterprise Institution Plan
Target:
1000+ Students
Features:

* White Label Platform
* Dedicated Domain
* Dedicated Support
* SSO
* Custom Reports
Expected Pricing:
₹3–10 Lakhs/year
Vendor Subscription Plans
Vendor Starter
Target:
Freelance Trainers
Limits:
2 Colleges
100 Students Each
Total:
200 Students
Features:

* Course Upload
* Assessments
* Attendance
* Reports
Expected Pricing:
₹10,000–₹25,000/year
Vendor Growth
Target:
5 Colleges
500 Students
Features:

* Trainer Dashboard
* Revenue Analytics
* Certification Management
Expected Pricing:
₹50,000–₹1 Lakh/year
Vendor Enterprise
Target:
Unlimited Colleges
Unlimited Students
Features:

* White Label
* Dedicated Portal
* API Integration
Expected Pricing:
₹2–5 Lakhs/year
Payment Gateway Architecture
Design complete payment system.
Support:
India

* Razorpay
* Cashfree
* PayU
International

* Stripe
* PayPal
Implement:

* Subscription Billing
* Recurring Payments
* Auto Renewal
* Invoice Generation
* GST Calculation
* Coupon Engine
* Referral Discounts
* Refund Management
* Payment Reconciliation
Generate:

* Payment Service Design
* Billing Microservice
* Invoice Service
* Subscription Service
Hybrid Microservice Architecture
Spring Boot Services
Java 21 + Spring Boot 3

* Authentication Service
* User Service
* Student Service
* Institution Service
* Vendor Service
* Subscription Service
* Billing Service
* Payment Service
* LMS Service
* Assessment Service
* Placement Service
* Certificate Service
Golang Services

* Compiler Gateway
* Contest Engine
* Realtime Leaderboard
* Notification Processor
* WebSocket Service
* High Throughput APIs
Python FastAPI Services

* AI Tutor
* Resume Analyzer
* Interview Generator
* Recommendation Engine
* Placement Prediction
* Learning Analytics
Multi-Tenant SaaS Architecture
Design tenant hierarchy:

```text
BitwiseLearn
│
├── Student
│
├── Institution
│   ├── Departments
│   ├── Trainers
│   ├── Students
│
└── Vendor
    ├── Colleges
    ├── Trainers
    ├── Students

```

Implement:

* Tenant Isolation
* Tenant Aware Authentication
* Tenant Aware Billing
* Tenant Aware Reports
Azure Cloud Architecture
Use:

* Azure Front Door
* Azure CDN
* Azure WAF
* Azure API Management
* Azure Application Gateway
* AKS
* Azure Service Bus or Kafka
* PostgreSQL Flexible Server
* MongoDB
* Redis
* Azure Blob Storage
* Azure Key Vault
* Azure Monitor
Online Compiler Design
Design a LeetCode/HackerRank level compiler.
Support:

* Java
* Python
* C
* C++
* JavaScript
* Go
Requirements:

* Isolated Containers
* gVisor
* Kata Containers
* KEDA Scaling
* Queue Based Execution
* Contest Mode
Support:
10,000+ Concurrent Code Executions
Cost and Revenue Engine
Create FinOps Dashboard.
Track:

* Azure Cost
* Cost Per Student
* Cost Per Tenant
* Cost Per Institution
* Cost Per Vendor
* Compiler Cost
* AI Cost
* Storage Cost
* Network Cost
Generate:
Pricing Formula

```text
Total Monthly Cost
+
20% Operational Buffer
+
30% Business Margin
=
Base Revenue

Subscription Price =
(Base Revenue ÷ Active Paying Users)

```

Target:
Minimum 3X ROI on Cloud Spend.
Analytics Dashboard
Generate dashboards for:
Students

* Learning Progress
* Coding Progress
* Assessment Scores
Institutions

* Placement Statistics
* Student Performance
* Trainer Effectiveness
Vendors

* Revenue
* Student Engagement
* Course Completion
Super Admin

* MRR
* ARR
* Churn Rate
* Active Tenants
* Cloud Cost
* Profitability
Security
Implement:

* JWT
* OAuth2
* Azure AD
* RBAC
* ABAC
* MFA
* Rate Limiting
* WAF
* Encryption
* OWASP Top 10 Compliance
CI/CD
Use GitHub Actions.
Generate:

* Dockerfiles
* Helm Charts
* Kubernetes YAML
* Terraform IaC
* Azure Deployment Scripts
Deliverables
Generate a complete enterprise solution including:

1. SaaS Business Model
2. Subscription Architecture
3. Pricing Engine
4. Student Pricing Plans
5. Institution Pricing Plans
6. Vendor Pricing Plans
7. Payment Gateway Integration Design
8. Billing Service Design
9. Revenue Forecast Model
10. Azure Cloud Architecture
11. High-Level System Design
12. Low-Level Design
13. C4 Architecture
14. Database ER Diagrams
15. Microservice Breakdown
16. Kubernetes Architecture
17. AKS Capacity Planning
18. Kafka/Event Architecture
19. Online Compiler Architecture
20. Security Architecture
21. Monitoring Architecture
22. FinOps Dashboard Design
23. CI/CD Architecture
24. Disaster Recovery Design
25. Cost Optimization Strategy
26. 150 User Deployment Plan
27. 1,000 User Deployment Plan
28. 10,000 User Deployment Plan
29. 100,000 User Deployment Plan
30. 5-Year Scalability Roadmap
The final solution must be production-ready, Azure-native, SaaS-focused, financially sustainable, capable of maintaining a minimum 3X revenue-to-cloud-cost ratio, and scalable from 150 users to 100,000+ users without major architectural redesign.
