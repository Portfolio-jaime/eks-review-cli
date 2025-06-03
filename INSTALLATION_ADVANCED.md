# Instalación Avanzada y Acceso Global a `kcli`

Esta guía explica cómo compilar tu herramienta CLI con el nombre `kcli` y cómo instalarla para que puedas ejecutarla directamente desde cualquier ubicación en tu terminal, sin necesidad de navegar a su directorio de código fuente o usar `./kcli`.

---

## 1. Cambiar el Nombre del Ejecutable Durante la Compilación

Por defecto, al compilar tu proyecto, el nombre del ejecutable podría basarse en el nombre del directorio o del módulo. Para asegurarte de que el binario se llame `kcli`, usa la opción `-o`:

```bash
go build -o kcli
```

Esto generará un archivo ejecutable llamado `kcli` (o `kcli.exe` en Windows) en el directorio actual del proyecto.

---

## 2. Instalar `kcli` para Acceso Global

Para ejecutar `kcli` desde cualquier lugar, el archivo ejecutable debe estar ubicado en un directorio que esté listado en la variable de entorno `PATH` de tu sistema.

### Opción A: Instalación Manual (Recomendada)

1. **Compila `kcli`:**

    ```bash
    go build -o kcli
    ```

2. **Mueve el ejecutable a un directorio en tu `PATH`:**

    - **Linux y macOS:**
        - Un directorio común es `/usr/local/bin` (requiere permisos de superusuario):

            ```bash
            sudo mv kcli /usr/local/bin/
            ```

        - Alternativamente, puedes usar un directorio personal como `$HOME/bin` o `$HOME/.local/bin`. Asegúrate de que esté en tu `PATH`:

            ```bash
            mkdir -p $HOME/bin
            mv kcli $HOME/bin/
            export PATH="$HOME/bin:$PATH"
            ```

            Añade la línea anterior a tu `~/.bashrc` o `~/.zshrc` para que sea permanente.

    - **Windows:**
        - Crea un directorio (ej. `C:\Program Files\kcli` o `C:\Users\TuUsuario\bin`).
        - Mueve el archivo `kcli.exe` a ese directorio.
        - Añade ese directorio a tu variable de entorno `PATH` a través de las "Variables de Entorno" en las propiedades del sistema.

3. **Verifica la instalación:**

    Abre una nueva terminal y ejecuta:

    ```bash
    kcli --help
    ```

    Debería mostrar la ayuda de tu CLI.

---

### Opción B: Usando `go install`

El comando `go install` compila e instala paquetes. Los binarios se colocan en el directorio `$GOPATH/bin` (o `$HOME/go/bin` si `GOPATH` no está definido).

1. **Asegúrate de que el directorio de binarios de Go esté en tu `PATH`:**

    ```bash
    export PATH="$(go env GOPATH)/bin:$PATH"
    # O, si GOPATH no está definido explícitamente:
    # export PATH="$HOME/go/bin:$PATH"
    ```

    Añade la línea anterior a tu archivo de configuración de shell.

2. **Ejecuta `go install`:**

    Desde la raíz de tu proyecto:

    ```bash
    go install
    ```

    El nombre del binario será el del módulo o directorio principal. Si quieres que se llame `kcli`, es más directo usar la Opción A.

---

### Opción C: Distribución a Otros Usuarios (Avanzado)

Si deseas distribuir `kcli` a usuarios que no necesariamente tienen Go instalado:

1. **Compilación Cruzada:**

    ```bash
    # Para Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o kcli_linux_amd64
    # Para Windows AMD64
    GOOS=windows GOARCH=amd64 go build -o kcli_windows_amd64.exe
    # Para macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -o kcli_darwin_amd64
    # Para macOS ARM64 (Apple Silicon)
    GOOS=darwin GOARCH=arm64 go build -o kcli_darwin_arm64
    ```

2. **Empaquetado:**
    - Crea archivos `.zip` o `.tar.gz` que contengan el binario, `README.md`, `COMMANDS.md`, y otros archivos necesarios.
    - Considera incluir un script de instalación o instrucciones claras.

3. **Releases de GitHub:**
    - Sube estos archivos empaquetados y/o los binarios directamente a la sección "Releases" de tu repositorio de GitHub.

4. **Gestores de Paquetes (más avanzado):**
    - Homebrew (macOS/Linux), Scoop (Windows), paquetes `.deb` o `.rpm` para Linux, etc.

---

## Recomendación Inicial

Para desarrollo propio y compartir con otros desarrolladores Go, la Opción A (compilación manual y mover a un directorio en el `PATH`) o una variante de la Opción B son las más directas y rápidas.

Para distribución más amplia, empieza con Releases de GitHub (Opción C.3) proporcionando binarios precompilados.

---

## 📚 Integración sugerida en el README

Puedes añadir una sección como esta en tu `README.md` principal:

```markdown
## ⚙️ Instalación Avanzada (Acceso Global como `kcli`)

Por defecto, después de compilar con `go build -o kcli`, puedes ejecutar la herramienta desde el directorio del proyecto con `./kcli`.

Si deseas poder ejecutar `kcli` desde cualquier ubicación en tu terminal, necesitarás instalar el binario en un directorio que esté en tu `PATH` del sistema.

Para instrucciones detalladas sobre cómo compilar con el nombre `kcli` e instalarlo globalmente en diferentes sistemas operativos, consulta nuestra [Guía de Instalación Avanzada](INSTALLATION_ADVANCED.md).
```
