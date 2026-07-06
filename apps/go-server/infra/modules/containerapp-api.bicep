// api (Echo) — Container App per BACKEND_REWRITE_PLAN.md section 6.1.
param name string
param location string
param environmentId string
param acrLoginServer string
param acrUsername string
@secure()
param acrPassword string
param imageTag string

@secure()
param databaseUrl string
@secure()
param jwtAccessSecret string
@secure()
param jwtRefreshSecret string
@secure()
param resetTokenSecret string
param frontendUrl string
@secure()
param azureStorageConnectionString string
param azureStorageContainer string
@secure()
param mqClient string
param codeExecutionServer string

resource api 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  properties: {
    environmentId: environmentId
    configuration: {
      ingress: {
        external: true
        targetPort: 8080
        transport: 'auto'
      }
      registries: [
        {
          server: acrLoginServer
          username: acrUsername
          passwordSecretRef: 'acr-password'
        }
      ]
      secrets: [
        { name: 'acr-password', value: acrPassword }
        { name: 'database-url', value: databaseUrl }
        { name: 'jwt-access-secret', value: jwtAccessSecret }
        { name: 'jwt-refresh-secret', value: jwtRefreshSecret }
        { name: 'reset-token-secret', value: resetTokenSecret }
        { name: 'azure-storage-connection-string', value: azureStorageConnectionString }
        { name: 'mq-client', value: mqClient }
      ]
    }
    template: {
      containers: [
        {
          name: 'api'
          image: '${acrLoginServer}/bitwise-api:${imageTag}'
          resources: {
            cpu: json('0.25')
            memory: '0.5Gi'
          }
          env: [
            { name: 'PORT', value: '8080' }
            { name: 'DATABASE_URL', secretRef: 'database-url' }
            { name: 'JWT_ACCESS_SECRET', secretRef: 'jwt-access-secret' }
            { name: 'JWT_REFRESH_SECRET', secretRef: 'jwt-refresh-secret' }
            { name: 'RESET_TOKEN_SECRET', secretRef: 'reset-token-secret' }
            { name: 'FRONTEND_URL', value: frontendUrl }
            { name: 'AZURE_STORAGE_CONNECTION_STRING', secretRef: 'azure-storage-connection-string' }
            { name: 'AZURE_STORAGE_CONTAINER', value: azureStorageContainer }
            { name: 'MQ_CLIENT', secretRef: 'mq-client' }
            { name: 'CODE_EXECUTION_SERVER', value: codeExecutionServer }
          ]
          probes: [
            {
              type: 'Liveness'
              httpGet: { path: '/health', port: 8080 }
              periodSeconds: 15
            }
            {
              type: 'Readiness'
              httpGet: { path: '/ready', port: 8080 }
              periodSeconds: 10
            }
            {
              type: 'Startup'
              httpGet: { path: '/health', port: 8080 }
              failureThreshold: 10
            }
          ]
        }
      ]
      scale: {
        // minReplicas: 2 avoids cold start on the critical path at the cost of
        // a slightly higher floor. Drop to 0 off-hours if cold start is
        // acceptable for your traffic pattern (see plan section 10, open
        // decision #3 — left at the safer default here).
        minReplicas: 2
        maxReplicas: 25
        rules: [
          {
            name: 'http-concurrency'
            http: { metadata: { concurrentRequests: '100' } }
          }
          {
            name: 'cpu'
            custom: {
              type: 'cpu'
              metadata: { type: 'Utilization', value: '70' }
            }
          }
        ]
      }
    }
  }
}

output fqdn string = api.properties.configuration.ingress.fqdn
