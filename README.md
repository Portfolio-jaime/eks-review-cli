eks-review-cli: Herramienta de Revisión de Clústeres de Kubernetes (EKS)
📖 Visión General
eks-review-cli es una herramienta de línea de comandos (CLI) escrita en Go, diseñada para simplificar la revisión y el diagnóstico de recursos en clústeres de Kubernetes, con un enfoque particular en entornos Amazon EKS. Esta CLI busca automatizar tareas repetitivas, estandarizar flujos de trabajo y proporcionar una visión rápida y clara del estado y la configuración de tus recursos de Kubernetes.

Actualmente, la CLI se enfoca en el subcomando monitor, permitiendo a los usuarios obtener un resumen del estado de los recursos clave (status) y visualizar los eventos del clúster (events).

✨ Características (Actuales y Planificadas)
monitor status: Proporciona un resumen tabular del estado de los Pods, Deployments, Services e Ingresses.
monitor events: Muestra los eventos recientes del clúster, con opciones de filtrado por tipo y namespace.
monitor nodes (Planificado): Información detallada de los nodos del clúster, incluyendo uso de recursos.
monitor logs (Planificado): Acceso y filtrado de logs de Pods/Deployments.
security (Planificado): Comandos para auditar Network Policies, RBAC, imágenes de contenedores y Secrets.
optimize (Planificado): Identificación de recursos no utilizados y revisión de configuraciones de autoescalado.
diagnose (Planificado): Herramientas para diagnosticar problemas específicos de Pods, Services e Ingresses.
🚀 Instalación
Para construir y ejecutar eks-review-cli, asegúrate de tener Go instalado (versión 1.18+ recomendada).

Clonar el Repositorio:

Bash

git clone https://github.com/Portfolio-jaime/eks-review-cli.git
cd eks-review-cli
Inicializar Módulos Go y Descargar Dependencias:

Bash

go mod tidy
Este comando descargará todas las librerías necesarias (Kubernetes client-go, AWS SDK v2, tablewriter, Cobra).

Compilar la CLI:

Bash

go build -o eks-review
Esto creará un ejecutable llamado eks-review en el directorio actual.

💡 Uso
Asegúrate de que tu kubeconfig esté configurado correctamente para apuntar a tu clúster de Kubernetes (Minikube, EKS, GKE, etc.). eks-review-cli leerá tu kubeconfig por defecto (~/.kube/config).

Estructura de Comandos:
eks-review
├── monitor               # Comandos para monitoreo y visibilidad
│   ├── status            # Resumen del estado de recursos (Pods, Deployments, Services, Ingresses)
│   └── events            # Visualización de eventos del cluster
├── security              # Comandos para seguridad y compliance (Planificado)
├── optimize              # Comandos para optimización y costos (Planificado)
└── diagnose              # Comandos para diagnóstico de problemas (Planificado)
Ejemplos de Comandos:
Ver el estado de los recursos en el namespace actual/por defecto:

Bash

./eks-review monitor status
Ver el estado de los recursos en un namespace específico:

Bash

./eks-review monitor status -n kube-system
Ver el estado de los recursos en todos los namespaces:

Bash

./eks-review monitor status --all-namespaces
Ver los eventos recientes en el namespace actual/por defecto:

Bash

./eks-review monitor events
Ver los eventos de tipo 'Warning' en el namespace my-app:

Bash

./eks-review monitor events -n my-app --type Warning
Ver todos los eventos en todos los namespaces:

Bash

./eks-review monitor events -n all
🏗️ Estructura del Proyecto
eks-review-cli/
├── cmd/
│   ├── diagnose.go     # Comandos de diagnóstico
│   ├── events.go       # Implementación de 'monitor events'
│   ├── monitor.go      # Comando base 'monitor'
│   ├── optimize.go     # Comandos de optimización
│   ├── root.go         # Comando raíz de la CLI
│   └── security.go     # Comandos de seguridad
├── go.mod              # Definición del módulo Go y dependencias
├── go.sum              # Sumas de verificación de dependencias
├── main.go             # Punto de entrada de la aplicación
├── README.md           # Este archivo
└── (otros archivos de configuración o scripts)
📊 Diagrama de Comandos (Estructura Actual y Futura)
Code snippet

graph TD
    A[eks-review] --> B(monitor)
    B --> C(status)
    B --> D(events)
    A --> E(security)
    A --> F(optimize)
    A --> G(diagnose)

    subgraph monitor_commands
        C
        D
    end

    style E fill:#f9f,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5;
    style F fill:#f9f,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5;
    style G fill:#f9f,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5;
    classDef planned fill:#f9f,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5;
    linkStyle 0,1,2,3,4,5,6 stroke:#333,stroke-width:2px;
Los nodos de color con línea discontinua (security, optimize, diagnose) representan funcionalidades planificadas.
🤝 Contribuciones
¡Las contribuciones son bienvenidas! Si tienes ideas para nuevas características, mejoras o correcciones de errores, no dudes en abrir un issue o enviar un pull request.