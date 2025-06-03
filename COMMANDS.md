# Referencia de Comandos de eks-review-cli

Este documento proporciona una lista detallada de todos los comandos disponibles en `eks-review-cli`, sus opciones y descripciones.

---

## Comando Raíz

`eks-review` es el comando principal para todas las operaciones.

### Flags Globales

- `-v, --verbose`  
  Habilita la salida detallada (logs de DEBUG).

---

## Subcomando `monitor`

Comandos para monitoreo y visibilidad del clúster.

**Uso general:**  
```bash
eks-review monitor [subcomando] [flags]
```

### `eks-review monitor status`

Muestra un resumen del estado de recursos como Pods, Deployments, Services e Ingresses.

**Flags:**
- `-A, --all-namespaces`  
  Lista los objetos en todos los namespaces.
- `-n, --namespace <namespace>`  
  Namespace específico para la consulta.

---

### `eks-review monitor events`

Muestra eventos recientes del clúster, útiles para la resolución de problemas.

**Flags:**
- `-n, --namespace <namespace>`  
  Namespace específico. Usa 'all' para todos los namespaces.
- `-T, --type <tipo_evento>`  
  Filtra eventos por tipo (ej., 'Warning', 'Normal').

---

### `eks-review monitor nodes`

Muestra información detallada sobre los nodos del clúster, incluyendo su estado, roles, versiones y uso de recursos (si el servidor de métricas está disponible).

**Flags:**  
(Este comando no tiene flags específicas adicionales por el momento).

---

### `eks-review monitor logs`

Imprime los logs de un contenedor en un pod, deployment o servicio.

**Flags:**
- `--pod <nombre_pod>`  
  Nombre del pod del que obtener logs.
- `--deployment <nombre_deployment>`  
  Nombre del deployment del que obtener logs.
- `--service <nombre_service>`  
  Nombre del service del que obtener logs.
- `-n, --namespace <namespace>`  
  Namespace específico.
- `-c, --container <nombre_contenedor>`  
  Nombre del contenedor. Si se omite, se elegirá el primero.
- `-f, --follow`  
  Stream de logs en tiempo real.
- `-p, --previous`  
  Muestra los logs de la instancia previa del contenedor si existe.
- `--grep <cadena_texto>`  
  Filtra logs por una cadena de texto.
- `--tail <lineas>`  
  Número de líneas desde el final de los logs a mostrar (por defecto: -1, todas).

---

## Subcomandos Planificados

### `eks-review security` *(Planificado)*

Auditoría de Network Policies, RBAC, imágenes de contenedores y Secrets.

### `eks-review optimize` *(Planificado)*

Identificación de recursos no utilizados y revisión de autoescalado.

### `eks-review diagnose` *(Planificado)*

Diagnóstico de problemas en Pods, Services e Ingresses.

---

## Ayuda

Para obtener ayuda sobre un comando específico, ejecuta:

```bash
eks-review [comando] [subcomando] --help
```