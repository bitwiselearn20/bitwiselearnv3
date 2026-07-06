// Azure Container Apps deployment for the Go rewrite (BACKEND_REWRITE_PLAN.md
// section 6). Consumption plan: bills per vCPU-second + memory-GiB-second
// only while replicas run.
//
// NOT VALIDATED against a live subscription — no Azure CLI/Bicep tooling was
// available in the environment this was written in. Before applying, run:
//   az bicep build --file main.bicep
//   az deployment group what-if -g <rg> -f main.bicep -p main.parameters.json
targetScope = 'resourceGroup'

@description('Short name used as a prefix for every resource (e.g. "bitwise").')
param namePrefix string = 'bitwise'

@description('Azure region for all resources.')
param location string = resourceGroup().location

@description('Environment tag: dev, staging, prod.')
param environmentName string = 'prod'

@description('Container image tag to deploy (set by CI to the commit SHA).')
param imageTag string = 'latest'

@description('MongoDB connection string (Atlas or self-hosted). Stored as a Container Apps secret.')
@secure()
param databaseUrl string

@description('JWT signing secrets. Stored as Container Apps secrets.')
@secure()
param jwtAccessSecret string
@secure()
param jwtRefreshSecret string
@secure()
param resetTokenSecret string

@description('Frontend origin for CORS.')
param frontendUrl string

@description('SMTP credentials for the worker (Gmail app password).')
param emailUser string
@secure()
param emailPass string

@description('Azure Blob Storage connection string (course files/certificates).')
@secure()
param azureStorageConnectionString string
param azureStorageContainer string = 'bitwise-learn'

@description('RabbitMQ connection string. First deploy provisions the in-cluster rabbitmq module below with default guest/guest credentials — set this to amqp://guest:guest@<rabbitmq module fqdn>:5672/ using that module\'s output, or point it at an external managed broker (CloudAMQP, Azure Service Bus) instead.')
@secure()
param mqClient string

var acrName = toLower(replace('${namePrefix}acr${environmentName}', '-', ''))
var envName = '${namePrefix}-env-${environmentName}'
var logAnalyticsName = '${namePrefix}-logs-${environmentName}'

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: logAnalyticsName
  location: location
  properties: {
    sku: { name: 'PerGB2018' }
    retentionInDays: 30
  }
}

resource acr 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  name: acrName
  location: location
  sku: { name: 'Basic' }
  properties: {
    adminUserEnabled: true
  }
}

resource containerAppEnv 'Microsoft.App/managedEnvironments@2024-03-01' = {
  name: envName
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: logAnalytics.properties.customerId
        sharedKey: logAnalytics.listKeys().primarySharedKey
      }
    }
  }
}

// RabbitMQ — single replica, ephemeral storage (see module for the
// persistence trade-off). Swap this module out for a managed broker
// (CloudAMQP, Azure Service Bus premium tier) if you need HA or durability
// across restarts.
module rabbitmq 'modules/containerapp-rabbitmq.bicep' = {
  name: 'rabbitmq'
  params: {
    name: '${namePrefix}-rabbitmq-${environmentName}'
    location: location
    environmentId: containerAppEnv.id
  }
}

module api 'modules/containerapp-api.bicep' = {
  name: 'api'
  params: {
    name: '${namePrefix}-api-${environmentName}'
    location: location
    environmentId: containerAppEnv.id
    acrLoginServer: acr.properties.loginServer
    acrUsername: acr.listCredentials().username
    acrPassword: acr.listCredentials().passwords[0].value
    imageTag: imageTag
    databaseUrl: databaseUrl
    jwtAccessSecret: jwtAccessSecret
    jwtRefreshSecret: jwtRefreshSecret
    resetTokenSecret: resetTokenSecret
    frontendUrl: frontendUrl
    azureStorageConnectionString: azureStorageConnectionString
    azureStorageContainer: azureStorageContainer
    mqClient: mqClient
    codeExecutionServer: 'https://${piston.outputs.fqdn}'
  }
}

module worker 'modules/containerapp-worker.bicep' = {
  name: 'worker'
  params: {
    name: '${namePrefix}-worker-${environmentName}'
    location: location
    environmentId: containerAppEnv.id
    acrLoginServer: acr.properties.loginServer
    acrUsername: acr.listCredentials().username
    acrPassword: acr.listCredentials().passwords[0].value
    imageTag: imageTag
    mqClient: mqClient
    emailUser: emailUser
    emailPass: emailPass
  }
}

module piston 'modules/containerapp-piston.bicep' = {
  name: 'piston'
  params: {
    name: '${namePrefix}-piston-${environmentName}'
    location: location
    environmentId: containerAppEnv.id
  }
}

output apiFqdn string = api.outputs.fqdn
output pistonFqdn string = piston.outputs.fqdn
output acrLoginServer string = acr.properties.loginServer
