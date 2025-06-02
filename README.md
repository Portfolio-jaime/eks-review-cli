# eks-review-cli

**Herramienta de Revisión de Clústeres de Kubernetes (EKS)**

---

## 📖 Visión General

**eks-review-cli** es una herramienta de línea de comandos (CLI) escrita en Go, diseñada para simplificar la revisión y el diagnóstico de recursos en clústeres de Kubernetes, con un enfoque especial en Amazon EKS. Su objetivo es automatizar tareas repetitivas, estandarizar flujos de trabajo y proporcionar una visión rápida y clara del estado y la configuración de tus recursos de Kubernetes.

Actualmente, la CLI se enfoca en el subcomando `monitor`, permitiendo a los usuarios obtener un resumen del estado de los recursos clave (`status`) y visualizar los eventos del clúster (`events`).

---

## ✨ Características

- **monitor status:** Resumen tabular del estado de Pods, Deployments, Services e Ingresses.
- **monitor events:** Visualización de eventos recientes del clúster, con opciones de filtrado por tipo y namespace.
- **monitor nodes** *(Planificado):* Información detallada de los nodos, incluyendo uso de recursos.
- **monitor logs** *(Planificado):* Acceso y filtrado de logs de Pods/Deployments.
- **security** *(Planificado):* Auditoría de Network Policies, RBAC, imágenes de contenedores y Secrets.
- **optimize** *(Planificado):* Identificación de recursos no utilizados y revisión de autoescalado.
- **diagnose** *(Planificado):* Diagnóstico de problemas en Pods, Services e Ingresses.

---

## 🚀 Instalación

Asegúrate de tener Go instalado (versión 1.18+ recomendada).

### 1. Clonar el repositorio

```bash
git clone https://github.com/Portfolio-jaime/eks-review-cli.git
cd eks-review-cli
```

### 2. Inicializar módulos Go y descargar dependencias

```bash
go mod tidy
```
> Este comando descargará todas las librerías necesarias (Kubernetes client-go, AWS SDK v2, tablewriter, Cobra).

### 3. Compilar la CLI

```bash
go build -o eks-review
```
> Esto creará un ejecutable llamado `eks-review` en el directorio actual.

---

## 💡 Uso

Asegúrate de que tu kubeconfig esté configurado correctamente para apuntar a tu clúster de Kubernetes (Minikube, EKS, GKE, etc.).  
Por defecto, eks-review-cli leerá tu kubeconfig en `~/.kube/config`.

### **Estructura de Comandos**

```
eks-review
├── monitor               # Comandos para monitoreo y visibilidad
│   ├── status            # Resumen del estado de recursos (Pods, Deployments, Services, Ingresses)
│   └── events            # Visualización de eventos del clúster
├── security              # Seguridad y compliance (Planificado)
├── optimize              # Optimización y costos (Planificado)
└── diagnose              # Diagnóstico de problemas (Planificado)
```

### **Ejemplos de Comandos**

- Ver el estado de los recursos en el namespace actual/por defecto:
  ```bash
  ./eks-review monitor status
  ```

- Ver el estado de los recursos en un namespace específico:
  ```bash
  ./eks-review monitor status -n kube-system
  ```

- Ver el estado de los recursos en todos los namespaces:
  ```bash
  ./eks-review monitor status --all-namespaces
  ```

- Ver los eventos recientes en el namespace actual/por defecto:
  ```bash
  ./eks-review monitor events
  ```

- Ver los eventos de tipo 'Warning' en el namespace `my-app`:
  ```bash
  ./eks-review monitor events -n my-app --type Warning
  ```

- Ver todos los eventos en todos los namespaces:
  ```bash
  ./eks-review monitor events -n all
  ```

---

## 🏗️ Estructura del Proyecto

```
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
```

---

## 📊 Diagrama de Comandos

```mermaid
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
```
> Los nodos de color con línea discontinua (*security*, *optimize*, *diagnose*) representan funcionalidades planificadas.

---

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas!  
Si tienes ideas para nuevas características, mejoras o correcciones de errores, no dudes en abrir un issue o enviar un pull request.

---