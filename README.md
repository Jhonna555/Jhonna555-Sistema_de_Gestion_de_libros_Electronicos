
SISTEMA DE GESTIÓN DE LIBROS ELECTRÓNICOS
============================================================

DESCRIPCIÓN GENERAL
------------------------------------------------------------
El Sistema de Gestión de Libros Electrónicos es una aplicación
backend desarrollada en Go (Golang) que permite administrar un
catálogo de libros electrónicos mediante una arquitectura por
capas, exposición de servicios REST y persistencia en SQL Server.

El sistema implementa operaciones CRUD completas, búsqueda,
estadísticas, control de concurrencia y pruebas automatizadas,
cumpliendo con los principios de Programación Orientada a Objetos
y buenas prácticas de arquitectura backend.

Proyecto desarrollado como parte del Proyecto Final Integrador
de la asignatura Programación Orientada a Objetos.


OBJETIVOS DEL PROYECTO
------------------------------------------------------------
- Implementar una aplicación web funcional con API REST.
- Aplicar principios de Programación Orientada a Objetos.
- Implementar mínimo 8 servicios REST.
- Integrar persistencia en SQL Server.
- Implementar control de concurrencia.
- Desarrollar pruebas unitarias, integración y aceptación.
- Documentar técnicamente el sistema.


ARQUITECTURA DEL SISTEMA
------------------------------------------------------------
Frontend (HTML)
    ↓
API REST (net/http)
    ↓
Catalog (Lógica de negocio)
    ↓
BookRepository (Interfaz)
    ↓
SQLServerStore (Persistencia)
    ↓
SQL Server


ESTRUCTURA DEL PROYECTO
------------------------------------------------------------
Sistema_Gestion_Libros_Electronicos/
|-- cmd/
|   |-- app/
|   |-- server/
|-- internal/
|   |-- model/
|   |-- catalog/
|   |-- storage/
|   |-- api/
|-- web/
|-- go.mod
|-- go.sum


FUNCIONALIDADES IMPLEMENTADAS
------------------------------------------------------------
- Crear libros electrónicos.
- Listar todos los libros.
- Obtener libro por ID.
- Actualizar información.
- Eliminar libros.
- Cambiar estado de disponibilidad.
- Buscar libros por criterio.
- Consultar estadísticas del catálogo.


ENDPOINTS REST
------------------------------------------------------------
GET     /api/libros
GET     /api/libros/{id}
POST    /api/libros
PUT     /api/libros/{id}
DELETE  /api/libros/{id}
PATCH   /api/libros/{id}/disponible
GET     /api/libros/buscar
GET     /api/estadisticas


BASE DE DATOS
------------------------------------------------------------
Motor utilizado: SQL Server

Tabla principal:

CREATE TABLE Libros (
    id INT PRIMARY KEY IDENTITY(1,1),
    titulo NVARCHAR(255) NOT NULL,
    autor NVARCHAR(255) NOT NULL,
    anio INT NOT NULL,
    disponible BIT NOT NULL
);


RESTAURAR BASE DE DATOS
------------------------------------------------------------
1. Abrir SQL Server Management Studio.
2. Click derecho en "Bases de datos".
3. Seleccionar "Restaurar base de datos".
4. Elegir "Dispositivo".
5. Agregar el archivo ubicado en backup/BibliotecaDB.bak.
6. Restaurar la base.
7. Ejecutar el servidor Go.


CÓMO EJECUTAR EL PROYECTO
------------------------------------------------------------
1. Restaurar la base de datos.
2. Verificar cadena de conexión en sqlserver_store.go.
3. Ejecutar desde la raíz del proyecto:

   go run ./cmd/server

Servidor disponible en:

   http://localhost:8080


PRUEBAS
------------------------------------------------------------
Ejecutar pruebas unitarias con:

   go test ./... -v

Las pruebas validan:
- Inserción
- Actualización
- Eliminación
- Búsqueda
- Cambio de disponibilidad


TECNOLOGÍAS UTILIZADAS
------------------------------------------------------------
- Go (Golang)
- SQL Server
- net/http
- sync.RWMutex
- Arquitectura REST
- Git & GitHub


CONCURRENCIA
------------------------------------------------------------
El sistema implementa control de concurrencia mediante
sync.RWMutex para permitir múltiples lecturas simultáneas
y proteger operaciones de escritura, evitando condiciones
de carrera.


AUTOR
------------------------------------------------------------
Jhonnatan Francisco Salazar Cadena
Ingeniería en Software
Programación Orientada a Objetos


NOTA FINAL
------------------------------------------------------------
El código fuente completo, pruebas, documentación y respaldo
de la base de datos se encuentran disponibles en este
repositorio para su revisión y ejecución.
