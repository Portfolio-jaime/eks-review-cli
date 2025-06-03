# Guía de Pruebas para eks-review-cli

Este documento proporciona una lista completa de comandos para probar la herramienta `eks-review-cli`, utilizando la información detallada de `COMMANDS.md` para cada comando y sus flags.

**Prerrequisitos:**
* Asegúrate de haber compilado tu CLI con `go build -o eks-review`.
* Ejecuta los comandos desde el directorio donde se encuentra el ejecutable `eks-review`, o asegúrate de que esté en tu PATH.
* Configura tu archivo `~/.kube/config` para apuntar a un clúster de Kubernetes accesible y con los permisos adecuados.

---

## Comando Raíz (`eks-review`)

* **Ayuda General:**
    ```bash
    ./eks-review --help
    ```
* **Flag Global `--verbose` / `-v`:**
    * Descripción: Habilita la salida detallada (logs de DEBUG).
    * Ejemplo de uso (aplicable a cualquier comando):
        ```bash
        ./eks-review --verbose monitor status
        ./eks-review -v monitor events -n default
        ```

---

## Subcomando `monitor`
* Descripción: Comandos para monitoreo y visibilidad del clúster.
* **Ayuda del subcomando `monitor`:**
    ```bash
    ./eks-review monitor --help
    ```

### 1. `eks-review monitor status`
* Descripción: Muestra un resumen del estado de recursos como Pods, Deployments, Services e Ingresses.
* **Comandos de prueba:**
    * Ver estado en el namespace por defecto (o el configurado en tu kubeconfig):
        ```bash
        ./eks-review monitor status
        ```
    * Ver estado en un namespace específico (ej. `kube-system`):
        ```bash
        ./eks-review monitor status --namespace kube-system
        ./eks-review monitor status -n kube-system
        ```
        * Flag: `-n, --namespace <namespace>`
            * Descripción: Si está presente, el ámbito del namespace para esta solicitud CLI.
    * Ver estado en todos los namespaces:
        ```bash
        ./eks-review monitor status --all-namespaces
        ./eks-review monitor status -A
        ```
        * Flag: `-A, --all-namespaces`
            * Descripción: Si es true, lista el/los objeto(s) solicitado(s) en todos los namespaces.
    * Ayuda específica para `status`:
        ```bash
        ./eks-review monitor status --help
        ```

### 2. `eks-review monitor events`
* Descripción: Muestra eventos recientes del clúster, útiles para la resolución de problemas.
* **Comandos de prueba:**
    * Ver eventos en el namespace por defecto:
        ```bash
        ./eks-review monitor events
        ```
    * Ver eventos en un namespace específico (ej. `default` o `all`):
        ```bash
        ./eks-review monitor events --namespace default
        ./eks-review monitor events -n all
        ```
        * Flag: `-n, --namespace <namespace>`
            * Descripción: Si está presente, el ámbito del namespace para esta solicitud CLI. Usa 'all' para todos los namespaces.
    * Filtrar eventos por tipo (ej. `Warning` o `Normal`):
        ```bash
        ./eks-review monitor events --type Warning
        ./eks-review monitor events -T Normal
        ```
        * Flag: `-T, --type <tipo_evento>`
            * Descripción: Filtra eventos por tipo (ej., 'Warning', 'Normal'). No sensible a mayúsculas.
    * Combinar namespace y tipo:
        ```bash
        ./eks-review monitor events -n kube-system -T Warning
        ```
    * Ayuda específica para `events`:
        ```bash
        ./eks-review monitor events --help
        ```

### 3. `eks-review monitor nodes`
* Descripción: Muestra información detallada sobre los nodos del clúster, incluyendo su estado, roles, versiones y uso de recursos (si el servidor de métricas está disponible).
* **Comandos de prueba:**
    * Ver información de los nodos:
        ```bash
        ./eks-review monitor nodes
        ```
    * Ayuda específica para `nodes`:
        ```bash
        ./eks-review monitor nodes --help
        ```
    *(Este comando no tiene flags específicas adicionales por el momento según `COMMANDS.md`).*

### 4. `eks-review monitor logs`
* Descripción: Imprime los logs de un contenedor en un pod, deployment o servicio.
* **Comandos de prueba (reemplaza `<...>` con nombres reales de tu clúster):**
    * Logs de un pod específico:
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod>
        ```
        * Flag: `--pod <nombre_pod>`
            * Descripción: Nombre del pod del que obtener logs.
    * Logs de un deployment:
        ```bash
        ./eks-review monitor logs --deployment <nombre-del-deployment>
        ```
        * Flag: `--deployment <nombre_deployment>`
            * Descripción: Nombre del deployment del que obtener logs.
    * Logs de un servicio:
        ```bash
        ./eks-review monitor logs --service <nombre-del-servicio>
        ```
        * Flag: `--service <nombre_service>`
            * Descripción: Nombre del service del que obtener logs.
    * Logs en un namespace específico:
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> -n <nombre-del-namespace>
        ```
        * Flag: `-n, --namespace <namespace>`
            * Descripción: Si está presente, el ámbito del namespace para esta solicitud CLI.
    * Logs de un contenedor específico:
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> -c <nombre-del-contenedor>
        ```
        * Flag: `-c, --container <nombre_contenedor>`
            * Descripción: Nombre del contenedor. Si se omite, se elegirá el primer contenedor del pod.
    * Seguir los logs (stream):
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> -f
        ```
        * Flag: `-f, --follow`
            * Descripción: Especificar si los logs deben ser transmitidos (streamed).
    * Logs de la instancia previa del contenedor:
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> -p
        ```
        * Flag: `-p, --previous`
            * Descripción: Si es true, imprime los logs de la instancia previa del contenedor en un pod si existe.
    * Filtrar logs con grep:
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> --grep "terminating"
        ```
        * Flag: `--grep <cadena_texto>`
            * Descripción: Filtrar logs por una cadena de texto.
    * Mostrar las últimas N líneas (tail):
        ```bash
        ./eks-review monitor logs --pod <nombre-del-pod> --tail 20
        ```
        * Flag: `--tail <lineas>`
            * Descripción: Líneas desde el final de los logs a mostrar. Si es -1, muestra todas las líneas (valor por defecto: -1).
    * Ayuda específica para `logs`:
        ```bash
        ./eks-review monitor logs --help
        ```

---

## Subcomandos Planificados (Placeholders)
Estos comandos existen pero solo imprimirán un mensaje indicando que han sido llamados, ya que su funcionalidad completa aún no está implementada. Puedes probar que se ejecutan y muestran su mensaje de ayuda.

### `eks-review security`
* Descripción: Auditoría de Network Policies, RBAC, imágenes de contenedores y Secrets.
* **Comandos de prueba:**
    ```bash
    ./eks-review security
    ./eks-review security --help
    ```

### `eks-review optimize`
* Descripción: Identificación de recursos no utilizados y revisión de autoescalado.
* **Comandos de prueba:**
    ```bash
    ./eks-review optimize
    ./eks-review optimize --help
    ```

### `eks-review diagnose`
* Descripción: Diagnóstico de problemas en Pods, Services e Ingresses.
* **Comandos de prueba:**
    ```bash
    ./eks-review diagnose
    ./eks-review diagnose --help
    ```

---

**Notas Adicionales para Probar:**

* **Nombres de Recursos Reales:** Para los comandos que interactúan con recursos específicos (como `logs --pod <nombre-del-pod>`), asegúrate de usar nombres que existan en tu clúster y en el namespace correcto. Utiliza `kubectl get pods`, `kubectl get deployments`, etc., para encontrar nombres válidos.
* **Interacción:** Observa la salida de cada comando. ¿Es la esperada? ¿Se formatea correctamente la tabla? ¿Los filtros funcionan como se indica?
* **Errores:** Prueba también casos que deberían generar errores (ej. un pod que no existe, un namespace incorrecto) para ver cómo maneja los errores tu CLI.
* **Ayuda de Comandos:** No olvides que la flag `--help` es tu amiga para cada comando y subcomando. Te mostrará todas las flags disponibles y una breve descripción de lo que hace el comando.