# Roadmap y Diferencial para eks-review-cli

Este documento describe las estrategias y funcionalidades propuestas para diferenciar `eks-review-cli` de herramientas estándar como `kubectl`, junto con un checklist de desarrollo a corto y mediano plazo.

## 🎯 Análisis Diferencial: ¿Cómo Destacar?

Si bien `kubectl` es la herramienta fundamental para interactuar con Kubernetes, `eks-review-cli` puede ofrecer un valor agregado significativo al enfocarse en:

1.  **Análisis y Perspectivas (Insights):**
    * Ir más allá del listado crudo de recursos.
    * Proporcionar resúmenes inteligentes que prioricen problemas o áreas de atención.
    * Detectar problemas comunes y ofrecer sugerencias.
    * Correlacionar información de diferentes fuentes para dar un contexto más rico.

2.  **Enfoque Específico en EKS:**
    * Integrar comprobaciones y listados que son particulares del ecosistema Amazon EKS (configuración del plano de control, node groups, IRSA, CNI, etc.).
    * Validar contra las buenas prácticas recomendadas por AWS para EKS.

3.  **Optimización y Gestión de Costos:**
    * Identificar recursos infrautilizados o no utilizados (ConfigMaps, Secrets, PVCs, PVs).
    * Proveer recomendaciones sobre el dimensionamiento de `requests` y `limits`.
    * (Avanzado) Sugerir optimizaciones de costos en nodos.

4.  **Seguridad Proactiva:**
    * Realizar auditorías básicas de RBAC (permisos excesivos).
    * Analizar NetworkPolicies para identificar configuraciones de riesgo.
    * Chequear configuraciones de seguridad en imágenes o ServiceAccounts.

5.  **Experiencia de Usuario Mejorada:**
    * Flujos de trabajo guiados para diagnóstico.
    * Salidas más amigables, con colores o sugerencias accionables.
    * Perfiles de revisión personalizables.

El objetivo no es reemplazar `kubectl`, sino complementarlo, ofreciendo una capa de inteligencia, automatización de tareas de revisión comunes, y un enfoque en EKS.

---

## ✅ Checklist de Desarrollo y Funcionalidades Diferenciales

Esta es una lista de posibles funcionalidades a desarrollar, priorizadas para agregar valor diferencial progresivamente.

### Fase 1: Mejoras en Monitorización y Primeros Insights (Próximos Días/Semanas)

* **[ ] `monitor status` Mejorado:**
    * **[ ]** En la salida de `monitor status`, resaltar en color (ej. rojo/amarillo) los recursos con problemas evidentes (ej. Pods en `Error` o `CrashLoopBackOff`, Deployments con 0 réplicas listas).
    * **[ ]** Añadir una pequeña sección de "Alertas Rápidas" al final de `monitor status` que liste los problemas más críticos encontrados (ej. "3 Pods fallando, 1 Nodo NotReady").
* **[ ] Nuevo Comando: `monitor health-check` (o `monitor insights-rapidos`)**
    * **[ ]** Implementación básica que verifique:
        * Nodos en estado `NotReady`.
        * Pods en estados problemáticos (`Failed`, `CrashLoopBackOff`, `ImagePullBackOff`, `Pending` por mucho tiempo).
        * Deployments/StatefulSets/Daemonsets que no tienen todas sus réplicas listas/disponibles.
        * Eventos de tipo `Warning` muy recientes (últimos 5-10 minutos).
    * **[ ]** Presentar un resumen claro de los hallazgos.
* **[ ] `monitor get pods` Mejorado:**
    * **[ ]** En la opción `-o wide` o una nueva flag (ej. `--show-events`), listar los últimos 2-3 eventos de tipo `Warning` asociados directamente a cada pod listado.
* **[ ] Documentación:**
    * **[ ]** Actualizar `COMMANDS.md` y `README.md` con los nuevos comandos o mejoras.

### Fase 2: Profundizando en el Análisis y EKS (Mediano Plazo)

* **[ ] `monitor get <recurso>` - Salida Enriquecida:**
    * **[ ]** Para `monitor get deployments`, añadir una opción `-o wide` que también muestre el estado de los ReplicaSets y Pods asociados, y si hay errores en el rollout.
    * **[ ]** Para `monitor get pvc`, mostrar si la PVC está vinculada y, si es posible, el nombre del PV.
* **[ ] Funcionalidades Específicas de EKS (Investigación e Implementación Inicial):**
    * **[ ]** **Investigar:** Identificar las APIs de AWS SDK Go v2 necesarias para obtener información del plano de control de EKS y Node Groups.
    * **[ ]** Nuevo Comando: `eks-review eks status` (o `monitor eks-info`)
        * **[ ]** Listar versión del clúster EKS, estado del plano de control.
        * **[ ]** Listar Node Groups: nombres, tipos de instancia, AMIs, estado.
    * **[ ]** Nuevo Comando: `eks-review eks check-irsa <namespace>/<serviceaccount_name>`
        * **[ ]** Verificar si un ServiceAccount está correctamente configurado para IRSA (anotaciones, rol de IAM existe y tiene la política de confianza correcta).
* **[ ] Módulo `optimize` (Básico):**
    * **[ ]** Nuevo Comando: `eks-review optimize unused-configmaps [-n namespace]`
        * **[ ]** Listar ConfigMaps que no parecen estar montados como volumen ni referenciados como `envFrom` por ningún Pod en el namespace.
    * **[ ]** Nuevo Comando: `eks-review optimize unused-secrets [-n namespace]`
        * **[ ]** Listar Secrets (excluyendo tipos `kubernetes.io/service-account-token` y otros gestionados por el sistema) que no parecen estar montados como volumen ni referenciados como `envFrom` por ningún Pod.

### Fase 3: Seguridad y Optimización Avanzada (Largo Plazo)

* **[ ] Módulo `security` (Básico):**
    * **[ ]** Nuevo Comando: `eks-review security overly-permissive-sa [-n namespace]`
        * **[ ]** Listar ServiceAccounts que tienen RoleBindings/ClusterRoleBindings a roles como `cluster-admin` o roles con `*` en resources/verbs.
    * **[ ]** Nuevo Comando: `eks-review security open-networkpolicies [-n namespace]`
        * **[ ]** Listar NetworkPolicies que permiten explícitamente `ingress` desde cualquier fuente (`0.0.0.0/0` o selector de pod vacío).
* **[ ] Módulo `diagnose` (Básico):**
    * **[ ]** Nuevo Comando: `eks-review diagnose pod <pod_name> [-n namespace]`
        * **[ ]** Mostrar un resumen consolidado: descripción del pod, últimos logs (ej. 20 líneas), últimos eventos relevantes, estado de los contenedores.
        * **[ ]** (Opcional) Ofrecer un modo interactivo simple.
* **[ ] Mejoras Generales:**
    * **[ ]** Implementar salida en color para resaltar problemas.
    * **[ ]** Explorar perfiles de revisión (ej. `./eks-review run-profile security-quick-check`).

---

**Próximos Pasos Inmediatos (Sugerencia para los próximos días):**

1.  **Empieza por la Fase 1:** `monitor status` mejorado o el nuevo comando `monitor health-check`. Estas funcionalidades ya pueden aportar mucho valor con la información que tu CLI ya puede obtener.
2.  **Elige una mejora pequeña y manejable** para empezar y ganar tracción. Por ejemplo, resaltar pods en mal estado en `monitor status`.
3.  **Itera:** Implementa, prueba tú mismo, y luego considera cómo mejorarlo.

Este roadmap es una guía; ajústalo según tus intereses, el feedback que recibas (si lo compartes) y el tiempo del que dispongas. ¡El objetivo es construir algo útil y diferenciado!
