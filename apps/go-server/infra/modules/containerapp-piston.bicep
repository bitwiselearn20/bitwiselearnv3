// piston (code execution) — Container App per plan section 6.3. Isolated,
// CPU-bound, hard-capped so it can never starve the api app of resources.
//
// NOTE: engineer-man/piston's isolate sandbox needs privileged
// container/cgroup access; verify your Container Apps environment supports
// it (Consumption plan sandboxing has historically been more restrictive
// than Dedicated/workload-profile environments — this is the same class of
// limitation this rewrite hit locally under Colima's VZ driver). If
// Consumption can't run it, use a Dedicated workload profile for this app.
param name string
param location string
param environmentId string

resource piston 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  properties: {
    environmentId: environmentId
    configuration: {
      ingress: {
        external: false
        targetPort: 2000
        transport: 'http'
      }
    }
    template: {
      containers: [
        {
          name: 'piston'
          image: 'ghcr.io/engineer-man/piston:latest'
          resources: {
            cpu: json('1.0')
            memory: '2.0Gi'
          }
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 15
        rules: [
          {
            name: 'cpu'
            custom: {
              type: 'cpu'
              metadata: { type: 'Utilization', value: '65' }
            }
          }
        ]
      }
    }
  }
}

output fqdn string = piston.properties.configuration.ingress.fqdn
