📘 Sistema de Gestión de Libros Electrónicos

📌 Descripción del Proyecto

Este proyecto corresponde al desarrollo de un Sistema de Gestión de Libros Electrónicos, implementado en el lenguaje de programación Go, como parte del Aprendizaje Autónomo 1 de la asignatura Programación Orientada a Objetos.

El sistema permite administrar un catálogo de libros electrónicos mediante operaciones básicas de gestión, aplicando conceptos de Programación Orientada a Objetos, programación funcional y organización modular mediante paquetes.

🎯 Objetivo del Sistema

Desarrollar un sistema que permita gestionar libros electrónicos, facilitando el registro, consulta, búsqueda, edición y eliminación de información, así como la persistencia de datos mediante archivos en formato JSON.

🧩 Funcionalidades

Registro de libros electrónicos

Listado de libros almacenados

Búsqueda por título, autor o categoría

Edición de información de libros

Eliminación de libros del sistema

Cambio del estado de disponibilidad

Guardado y carga de información en archivos JSON

🧱 Estructura del Proyecto
cmd/app            → Punto de entrada del sistema
internal/model     → Estructuras de datos (Libro)
internal/catalog   → Lógica de negocio
internal/storage   → Persistencia de datos (JSON)
internal/ui        → Interfaz de usuario (menú en consola)
data/              → Archivos de almacenamiento

⚙️ Tecnologías Utilizadas

Lenguaje de programación: Go

Paradigmas aplicados: Programación Orientada a Objetos y Programación Funcional

Persistencia de datos: Archivos JSON

Control de versiones: Git y GitHub

▶️ Ejecución del Proyecto

Para ejecutar el sistema, ubíquese en la carpeta raíz del proyecto y ejecute el siguiente comando:

go run ./cmd/app

📁 Repositorio GitHub

🔗 https://github.com/Jhonna555/Jhonna555-Sistema_de_Gestion_de_libros_Electronicos.git

👨‍🎓 Autor

Jhonnatan Francisco Salazar Cadena

Carrera: Ingeniería en Software

Materia: Programación Orientada a Objetos

📝 Observaciones

Este proyecto fue desarrollado con fines académicos, siguiendo las instrucciones establecidas para la planificación y diseño de un sistema de gestión empresarial.
